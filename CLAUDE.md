# 五子棋游戏项目规范

## 1. 项目概述

这是一个五子棋（Gomoku）游戏项目，包含人机对战、双人对战和等级成长系统。

### 技术栈

- **前端**: Vue 3 + TypeScript + Vite
- **后端**: Go + Gin
- **数据库**: PostgreSQL
- **AI算法**: Minimax + Alpha-Beta剪枝

### 项目结构

```
ai-coding-game-gomoku/
├── frontend/          # Vue3前端项目
│   ├── src/
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面视图
│   │   ├── composables/   # 组合式API
│   │   ├── stores/        # 状态管理
│   │   ├── utils/         # 工具函数
│   │   └── types/         # TypeScript类型
│   └── index.html
├── backend/           # Go后端项目
│   ├── cmd/
│   │   └── server/    # 入口文件
│   ├── internal/
│   │   ├── handlers/  # HTTP处理器
│   │   ├── services/  # 业务逻辑
│   │   ├── models/    # 数据模型
│   │   ├── repository/# 数据访问层
│   │   └── ai/        # AI算法
│   └── go.mod
├── doc/               # 文档
└── CLAUDE.md          # 项目规范（本文件）
```

---

## 2. 数据库设计

### 2.1 PostgreSQL表结构

```sql
-- 玩家等级表
CREATE TABLE player_ranks (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(64) UNIQUE NOT NULL,  -- 用户唯一标识
    username VARCHAR(64) NOT NULL,
    rank_level INTEGER DEFAULT 1,         -- 等级 1-100+
    rank_segment VARCHAR(16) DEFAULT '新手', -- 段位
    experience INTEGER DEFAULT 0,          -- 经验值
    total_wins INTEGER DEFAULT 0,          -- 总胜场
    total_losses INTEGER DEFAULT 0,        -- 总负场
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 对局记录表（可选）
CREATE TABLE game_records (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    game_mode VARCHAR(16) NOT NULL,        -- 'ai' 或 'pvp'
    difficulty VARCHAR(16),                -- 'easy', 'medium', 'hard'
    result VARCHAR(16) NOT NULL,           -- 'win', 'lose', 'draw'
    experience_gained INTEGER DEFAULT 0,
    played_at TIMESTAMP DEFAULT NOW()
);
```

### 2.2 段位经验值对照表

| 段位 | 等级范围 | 升级所需经验 |
|------|----------|--------------|
| 新手 | 1-10 | 0 |
| 初级 | 11-20 | 100 |
| 中级 | 21-30 | 500 |
| 高级 | 31-40 | 1500 |
| 专家 | 41-50 | 4000 |
| 大师 | 51+ | 10000 |

### 2.3 经验值获取规则

| 对局结果 | 经验值 |
|----------|--------|
| 击败简单AI | +10 |
| 击败中级AI | +25 |
| 击败困难AI | +50 |
| 负于AI | +5 |
| 平局 | +3 |

---

## 3. AI算法设计

### 3.1 难度级别

| 难度 | 算法 | 搜索深度 | 特点 |
|------|------|----------|------|
| 简单 | 随机 + 基础防守 | 0-1 | 随机落子，概率性防守 |
| 中等 | Minimax | 2-3层 | 具备攻守能力 |
| 困难 | Minimax + Alpha-Beta | 4-6层 | 深度计算，识别棋型 |

### 3.2 棋型评估

AI评估以下棋型（从高到低）：

1. **连五**（AAAAA/BBBBB）：必胜，优先级最高
2. **活四**（AAAA_ / _AAAA）：必胜
3. **冲四**（AAA_A / A_AAA）：即将获胜
4. **活三**（AA_A / A_AA / _AA_）：重要进攻棋型
5. **眠三**（_AAA_）：中等威胁
6. **活二**：基础进攻形

### 3.3 评估函数

```go
// 棋型分值
const (
    ScoreFive   int = 100000  // 连五
    ScoreFour   int = 10000   // 活四/冲四
    ScoreThree  int = 1000    // 活三
    ScoreTwo    int = 100     // 活二
    ScoreOne    int = 10      // 眠三/眠二
)
```

---

## 4. API设计

### 4.1 玩家相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/player/register | 注册/创建玩家 |
| GET | /api/player/:userId | 获取玩家信息 |
| PUT | /api/player/:userId | 更新玩家信息 |

### 4.2 对局相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/game/start | 开始游戏 |
| POST | /api/game/move | 落子 |
| POST | /api/game/end | 结束游戏 |

### 4.3 数据格式

```typescript
// 玩家信息
interface Player {
  userId: string;
  username: string;
  rankLevel: number;
  rankSegment: string;
  experience: number;
  totalWins: number;
  totalLosses: number;
}

// 游戏状态
interface GameState {
  gameId: string;
  board: number[][];  // 0:空, 1:黑, 2:白
  currentPlayer: 1 | 2;
  gameMode: 'ai' | 'pvp';
  difficulty?: 'easy' | 'medium' | 'hard';
  status: 'playing' | 'won' | 'draw';
  winner?: 1 | 2;
}
```

---

## 5. 前端开发规范

### 5.1 技术规范

- **Vue 3**：使用Composition API + `<script setup>`
- **TypeScript**：严格模式，开启所有严格检查
- **状态管理**：Pinia
- **样式**：SCSS 或 CSS Variables

### 5.2 组件结构

```
components/
├── ChessBoard.vue       # 棋盘组件
├── ChessPiece.vue       # 棋子组件
├── GameInfo.vue         # 游戏信息（当前玩家、状态）
├── RankBadge.vue        # 等级徽章
├── ModeSelector.vue     # 模式选择器
└── GameResult.vue       # 游戏结果弹窗
```

### 5.3 页面视图

```
views/
├── HomeView.vue         # 主菜单
├── GameView.vue         # 游戏页面
└── RankView.vue         # 等级查看
```

---

## 6. 后端开发规范

### 6.1 目录结构

```
internal/
├── handlers/       # HTTP处理器
│   ├── player.go
│   └── game.go
├── services/       # 业务逻辑
│   ├── player_service.go
│   └── game_service.go
├── models/         # 数据模型
│   ├── player.go
│   └── game.go
├── repository/     # 数据访问
│   ├── db.go
│   └── player_repo.go
├── ai/             # AI算法
│   ├── evaluator.go
│   ├── minimax.go
│   └── strategy.go
└── config/         # 配置
    └── config.go
```

### 6.2 命名规范

- **Go**: 使用驼峰命名，公开方法首字母大写，私有方法小写
- **数据库表**: 使用蛇蛇命名法（snake_case）
- **API路径**: 使用RESTful风格，路径用连字符

---

## 7. 开发流程

### 7.1 项目结构

```
ai-coding-game-gomoku/
├── frontend/                    # Vue3前端项目
│   ├── src/
│   │   ├── components/          # 组件
│   │   │   ├── ChessBoard.vue   # 棋盘组件
│   │   │   ├── ChessPiece.vue   # 棋子组件
│   │   │   ├── GameInfo.vue     # 游戏信息
│   │   │   ├── GameControl.vue  # 游戏控制
│   │   │   └── RankBadge.vue    # 等级徽章
│   │   ├── views/               # 页面视图
│   │   │   ├── HomeView.vue     # 主菜单
│   │   │   └── GameView.vue     # 游戏页面
│   │   ├── stores/              # 状态管理
│   │   │   ├── game.ts          # 游戏状态
│   │   │   └── player.ts        # 玩家状态
│   │   ├── router/              # 路由配置
│   │   ├── types/               # TypeScript类型
│   │   └── utils/               # 工具函数
│   │       └── api.ts           # API调用
│   ├── package.json
│   ├── vite.config.ts
│   └── index.html
├── backend/                     # Go后端项目
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # 入口文件
│   ├── internal/
│   │   ├── handlers/            # HTTP处理器
│   │   │   ├── game.go
│   │   │   └── player.go
│   │   ├── services/            # 业务逻辑
│   │   │   ├── game_service.go
│   │   │   └── player_service.go
│   │   ├── models/              # 数据模型
│   │   │   ├── game.go
│   │   │   └── player.go
│   │   ├── repository/          # 数据访问
│   │   │   └── db.go
│   │   ├── ai/                  # AI算法
│   │   │   ├── strategy.go
│   │   │   └── minimax.go
│   │   └── config/              # 配置
│   │       └── config.go
│   ├── db/
│   │   └── init.sql             # 数据库初始化
│   └── go.mod
├── doc/                         # 文档
│   └── PRD.md                   # 产品设计文档
├── scripts/                     # 启动脚本
│   ├── run.sh                   # 一键启动脚本
│   └── run-backend.sh           # 后端启动脚本
├── CLAUDE.md                    # 项目规范
└── README.md                    # 项目说明
```

### 7.2 初始化项目

```bash
# 前端
cd frontend
npm install

# 后端（如需后端）
cd backend
go mod download
```

### 7.3 运行项目

#### 方式1：前端独立运行（推荐，无需Go）

```bash
cd frontend
npm run dev
```

访问 http://localhost:5173

#### 方式2：前端 + 后端

```bash
# 终端1：启动后端
cd backend
go run cmd/server/main.go

# 终端2：启动前端
cd frontend
npm run dev
```

#### 方式3：使用启动脚本

```bash
./scripts/run.sh
```

### 7.4 数据库初始化（如需PostgreSQL）

```bash
# 创建数据库
createdb gomoku

# 运行迁移
psql -U postgres -d gomoku -f backend/db/init.sql
```

### 7.5 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 5432 | 数据库端口 |
| DB_USER | postgres | 数据库用户 |
| DB_PASSWORD | postgres | 数据库密码 |
| DB_NAME | gomoku | 数据库名称 |
| SERVER_PORT | 8080 | 服务端口 |

---

## 8. 验收检查点

### 功能验收

- [ ] 15×15棋盘正确渲染
- [ ] 黑棋先行，轮流落子正确
- [ ] 连五/长连判定正确
- [ ] 三种AI难度可切换
- [ ] 双人对战正常
- [ ] 等级经验值计算正确
- [ ] 数据库正确存储和读取

### 浏览器测试（必须执行）

每次代码修改后，必须在Chrome浏览器中验证：

```bash
# 1. 确保前端服务运行中
curl -s http://localhost:5173 | head -5

# 2. 在Chrome浏览器中打开 http://localhost:5173
# 3. 检查控制台是否有报错
# 4. 验证页面正常渲染，无504/Outdated Dep等错误
```

### 深度功能测试（必须执行）

每次发布前必须验证以下功能：

#### 1. AI智能测试（根据victory.md口诀）
- [ ] **简单难度**：AI能防守四连（"冲四必挡"）
- [ ] **简单难度**：AI有70%概率防守活三（"活三必应"）
- [ ] **简单难度**：AI能进攻（自己能形成四连）
- [ ] **中等难度**：AI必须防守三连及以上（"活三必应"）
- [ ] **中等难度**：AI会主动进攻，评估落子位置
- [ ] **困难难度**：AI防守权重更高，优先处理威胁
- [ ] **困难难度**：AI有明显的攻防策略，根据棋型评估选点

#### 2. 页面导航测试
- [ ] 首页可以正常显示
- [ ] 选择难度后能进入游戏页面
- [ ] 双人对战能正常进入游戏
- [ ] 点击"返回菜单"能正确返回首页（不能一直加载）
- [ ] "重新开始"能重置棋盘
- [ ] "悔棋"能撤销上一步

#### 3. 游戏流程测试
- [ ] 黑棋先行
- [ ] 落子后切换玩家
- [ ] 连成5子判定获胜
- [ ] 棋盘满判定平局
- [ ] 游戏结束后不能继续落子

**常见问题排查：**

- **504 Outdated Optimize Dep**: 执行 `rm -rf node_modules/.vite` 清除缓存，然后重启开发服务器
- **空白页面**: 检查控制台错误，检查端口是否被占用
- **API请求失败**: 确认后端服务是否启动（需要Go）
- **返回菜单一直加载**: 检查router.push是否正确调用，检查gameStore.clearGame()是否清理状态
- **AI太弱**: 检查evaluatePosition函数是否正确实现，是否能识别四连/三连

### 性能要求

- 简单/中等AI响应 < 100ms
- 困难AI响应 < 3秒
- 前端帧率 > 30fps

---

## 9. 语言设定

### 9.1 默认语言

- **游戏语言**：简体中文（zh-CN）
- 所有UI文本、按钮、提示信息均使用中文
- 不提供多语言切换功能

### 9.2 文本规范

| 场景 | 中文文本 |
|------|----------|
| 主菜单 | 开始游戏、我的等级、双人对战、人机对战 |
| 难度选择 | 简单、中级、困难 |
| 游戏状态 | 轮到你落子、黑方回合、白方回合 |
| 游戏结果 | 你赢了！你输了！平局！ |
| 按钮 | 重新开始、悔棋、认输、返回菜单 |
| 段位 | 新手、初级、中级、高级、专家、大师 |

---

## 10. 电脑角色设定

### 10.1 AI角色昵称

电脑玩家使用以下中文昵称：

| 难度 | 昵称 | 描述 |
|------|------|------|
| 简单 | 小白 | 新手级别的AI，反应迟缓 |
| 中级 | 小灵 | 具备一定思考能力的AI |
| 困难 | 大师 | 顶尖高手，极难战胜 |

### 10.2 角色特点

- **小白**：落子随机会有长时间思考的假象，偶尔会走出明显失误
- **小灵**：思考时间中等，会防守但进攻意识不强
- **大师**：落子迅速，攻防兼备，擅长制造陷阱

---

## 11. 注意事项

1. **前后端分离**：前端通过HTTP API与后端通信
2. **CORS配置**：后端需配置跨域资源共享
3. **本地存储**：可先用localStorage作为过渡，后续迁移到PostgreSQL
4. **WebSocket**（可选）：实时双人对战可用WebSocket增强体验