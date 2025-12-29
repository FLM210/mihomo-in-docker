package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flm210/mihomo-in-k8s/parsers"

	"gopkg.in/yaml.v3"
)

// GenerateConfigFromEnv 从环境变量获取订阅链接和过滤器，生成mihomo配置并保存到本地
func GenerateConfigFromEnv() error {
	// 从环境变量获取订阅链接
	subscriptionURL := os.Getenv("SUBSCRIPTION_URL")
	if subscriptionURL == "" {
		return fmt.Errorf("环境变量 SUBSCRIPTION_URL 未设置")
	}

	// 从环境变量获取过滤器（可选）
	filterStr := os.Getenv("FILTER_KEYWORDS")
	var filters []string
	if filterStr != "" {
		filters = strings.Split(filterStr, ",")
		// 去除每个过滤器值的空格
		for i, f := range filters {
			filters[i] = strings.TrimSpace(f)
		}
	}

	// 获取输出文件名（可选，默认为 config.yaml）
	outputFile := os.Getenv("OUTPUT_FILE")
	if outputFile == "" {
		outputFile = "config.yaml"
	}

	// 获取输出格式（可选，默认为 yaml）
	outputFormat := os.Getenv("OUTPUT_FORMAT")
	if outputFormat == "" {
		outputFormat = "yaml"
	}

	// 获取订阅内容
	subContent, err := FetchSubscription(subscriptionURL)
	if err != nil {
		return fmt.Errorf("获取订阅内容失败: %v", err)
	}

	// 转换订阅为 mihomo 配置
	var config *map[string]interface{}
	if len(filters) > 0 {
		config, err = parsers.ConvertToMihomo(subContent, filters...)
	} else {
		config, err = parsers.ConvertToMihomo(subContent)
	}
	if err != nil {
		return fmt.Errorf("转换为 mihomo 配置失败: %v", err)
	}

	// 保存配置到本地文件
	if strings.ToLower(outputFormat) == "json" {
		jsonBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("生成 JSON 配置失败: %v", err)
		}
		err = os.WriteFile(outputFile, jsonBytes, 0644)
		if err != nil {
			return fmt.Errorf("保存 JSON 配置文件失败: %v", err)
		}
	} else {
		// 默认为 YAML 格式
		yamlBytes, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("生成 YAML 配置失败: %v", err)
		}
		err = os.WriteFile(outputFile, yamlBytes, 0644)
		if err != nil {
			return fmt.Errorf("保存 YAML 配置文件失败: %v", err)
		}
	}

	fmt.Printf("成功生成配置文件: %s (格式: %s)\n", outputFile, strings.ToUpper(outputFormat))
	return nil
}

func FetchSubscription(url string) (string, error) {
	// Create a custom HTTP client that ignores HTTPS certificate errors
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Add timeout for requests
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
