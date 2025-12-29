#!/bin/bash

# 生成配置文件
./config-generator 

# 检查配置文件是否存在
if [ -f "config.yaml" ]; then
    echo "启动 mihomo..."
    exec /usr/local/bin/mihomo -d /root/.config/mihomo -f config.yaml
else
    echo "配置文件不存在，无法启动 mihomo"
    exit 1
fi