// 游戏模式
export type GameMode = 'ai' | 'pvp'

// 难度级别
export type Difficulty = 'easy' | 'medium' | 'hard'

// 玩家颜色
export type PlayerColor = 0 | 1 | 2 // 0:空, 1:黑, 2:白

// 游戏状态
export type GameStatus = 'playing' | 'won' | 'draw'

// 位置
export interface Position {
  x: number
  y: number
}

// 游戏状态
export interface GameState {
  gameId: string
  board: number[][] // 0:空, 1:黑, 2:白
  currentPlayer: PlayerColor
  gameMode: GameMode
  difficulty?: Difficulty
  status: GameStatus
  winner?: PlayerColor
  lastMove?: Position
  moveHistory: Position[]
}

// 玩家信息
export interface Player {
  userId: string
  username: string
  rankLevel: number
  rankSegment: string
  experience: number
  totalWins: number
  totalLosses: number
}

// AI角色配置
export const AICharacters: Record<Difficulty, { name: string; description: string }> = {
  easy: { name: '小白', description: '新手级别的AI，反应迟缓' },
  medium: { name: '小灵', description: '具备一定思考能力的AI' },
  hard: { name: '大师', description: '顶尖高手，极难战胜' }
}

// 段位信息
export interface RankSegment {
  name: string
  minLevel: number
  expThreshold: number
}

export const RankSegments: RankSegment[] = [
  { name: '新手', minLevel: 1, expThreshold: 0 },
  { name: '初级', minLevel: 11, expThreshold: 100 },
  { name: '中级', minLevel: 21, expThreshold: 500 },
  { name: '高级', minLevel: 31, expThreshold: 1500 },
  { name: '专家', minLevel: 41, expThreshold: 4000 },
  { name: '大师', minLevel: 51, expThreshold: 10000 }
]

// API响应
export interface ApiResponse<T> {
  code: number
  data?: T
  error?: string
}