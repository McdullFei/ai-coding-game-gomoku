#!/bin/bash

# 后端启动脚本

cd "$(dirname "$0")/../backend"

echo "启动后端服务..."
go mod download
go run cmd/server/main.go