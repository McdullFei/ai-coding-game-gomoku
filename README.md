# 五子棋游戏 (Gomoku)

一款支持人机对战、双人对战的五子棋游戏，包含等级成长系统。

## 功能特性

- **人机对战**：三种AI难度（小白/小灵/大师）
- **双人对战**：双人对战模式
- **等级系统**：段位升级（新手→初级→中级→高级→专家→大师）
- **中文界面**：完整的简体中文UI

## 技术栈

- **前端**: Vue 3 + TypeScript + Vite + Pinia
- **后端**: Go + Gin
- **数据库**: PostgreSQL

## 快速开始

### 前端运行（无需后端）

前端已内置本地游戏模式，可独立运行：

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:5173

### 完整运行（前端 + 后端）

需要安装 [Go](https://go.dev/dl/) 和 [Node.js](https://nodejs.org/)

```bash
# 1. 启动后端
cd backend
go mod download
go run cmd/server/main.go

# 2. 启动前端（新终端）
cd frontend
npm install
npm run dev
```

访问 http://localhost:5173

### 使用启动脚本

```bash
# 方式1：同时启动前端+后端（如已安装Go）
./scripts/run.sh

# 方式2：仅启动前端
cd frontend && npm run dev
```

## 游戏规则

- 棋盘：15×15
- 黑棋先行，轮流落子
- 横、竖、斜任意方向连成5子及以上获胜

## 经验值规则

| 结果 | 经验值 |
|------|--------|
| 击败简单AI | +10 |
| 击败中级AI | +25 |
| 击败困难AI | +50 |
| 负于AI | +5 |
| 平局 | +3 |

## 项目结构

```
ai-coding-game-gomoku/
├── frontend/          # Vue3前端
├── backend/           # Go后端
├── doc/               # 文档
├── scripts/           # 启动脚本
└── CLAUDE.md          # 项目规范
```

## 环境变量

### 后端

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 5432 | 数据库端口 |
| DB_USER | postgres | 数据库用户 |
| DB_PASSWORD | postgres | 数据库密码 |
| DB_NAME | gomoku | 数据库名称 |
| SERVER_PORT | 8080 | 服务端口 |

## 许可证

MIT