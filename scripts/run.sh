#!/bin/bash

# 五子棋游戏启动脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== 五子棋游戏启动脚本 ===${NC}"

# 检查Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}错误: Node.js未安装${NC}"
    exit 1
fi

# 启动前端
echo -e "${YELLOW}启动前端服务...${NC}"
cd "$(dirname "$0")/../frontend"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}安装前端依赖...${NC}"
    npm install
fi

npm run dev &
FRONTEND_PID=$!

# 等待前端启动
sleep 3

# 检查Go是否安装，如果安装则启动后端
if command -v go &> /dev/null; then
    echo -e "${YELLOW}启动后端服务...${NC}"
    cd "$(dirname "$0")/../backend"
    go run cmd/server/main.go &
    BACKEND_PID=$!
    echo -e "后端: http://localhost:8080"
else
    echo -e "${YELLOW}注意: Go未安装，后端服务将不会启动${NC}"
    echo -e "${YELLOW}前端将使用模拟数据运行${NC}"
fi

echo -e "${GREEN}=== 服务启动完成 ===${NC}"
echo -e "前端: http://localhost:5173"
echo ""
echo -e "按 ${RED}Ctrl+C${NC} 停止服务"

# 捕获Ctrl+C
if [ -n "$BACKEND_PID" ]; then
    trap "kill $FRONTEND_PID $BACKEND_PID 2>/dev/null; exit" INT TERM
else
    trap "kill $FRONTEND_PID 2>/dev/null; exit" INT TERM
fi

wait