#!/bin/bash

echo "正在启动 Blog API 服务..."
echo "Swagger 文档地址: http://localhost:8080/swagger/index.html"
echo ""

# 编译并运行项目
go run cmd/api/main.go