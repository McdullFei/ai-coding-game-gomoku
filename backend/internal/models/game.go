package models

import (
	"fmt"
	"math/rand"
	"time"
)

// GameMode 游戏模式
type GameMode string

const (
	ModeAI  GameMode = "ai"
	ModePVP GameMode = "pvp"
)

// Difficulty 难度级别
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// PlayerColor 玩家颜色
type PlayerColor int

const (
	Empty PlayerColor = 0
	Black PlayerColor = 1 // 黑棋
	White PlayerColor = 2 // 白棋
)

// GameStatus 游戏状态
type GameStatus string

const (
	StatusPlaying GameStatus = "playing"
	StatusWon     GameStatus = "won"
	StatusDraw    GameStatus = "draw"
)

// Position 落子位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GameState 游戏状态
type GameState struct {
	GameID       string      `json:"gameId"`
	Board        [][]int     `json:"board"` // 0:空, 1:黑, 2:白
	CurrentPlayer PlayerColor `json:"currentPlayer"`
	GameMode     GameMode    `json:"gameMode"`
	Difficulty   Difficulty  `json:"difficulty,omitempty"`
	Status       GameStatus  `json:"status"`
	Winner       PlayerColor `json:"winner,omitempty"`
	LastMove     *Position   `json:"lastMove,omitempty"`
	MoveHistory  []Position  `json:"moveHistory"`
}

// NewGameState 创建新的游戏状态
func NewGameState(gameMode GameMode, difficulty Difficulty) *GameState {
	board := make([][]int, 15)
	for i := range board {
		board[i] = make([]int, 15)
	}
	return &GameState{
		GameID:       generateGameID(),
		Board:        board,
		CurrentPlayer: Black, // 黑棋先行
		GameMode:     gameMode,
		Difficulty:   difficulty,
		Status:       StatusPlaying,
		MoveHistory:  make([]Position, 0),
	}
}

// IsValidPosition 检查位置是否有效
func (g *GameState) IsValidPosition(x, y int) bool {
	return x >= 0 && x < 15 && y >= 0 && y < 15
}

// IsEmpty 检查位置是否为空
func (g *GameState) IsEmpty(x, y int) bool {
	return g.IsValidPosition(x, y) && g.Board[y][x] == 0
}

// PlacePiece 落子
func (g *GameState) PlacePiece(x, y int, color PlayerColor) bool {
	if !g.IsEmpty(x, y) {
		return false
	}
	g.Board[y][x] = int(color)
	g.LastMove = &Position{X: x, Y: y}
	g.MoveHistory = append(g.MoveHistory, Position{X: x, Y: y})
	return true
}

// SwitchPlayer 切换玩家
func (g *GameState) SwitchPlayer() {
	if g.CurrentPlayer == Black {
		g.CurrentPlayer = White
	} else {
		g.CurrentPlayer = Black
	}
}

// CheckWin 检查是否获胜
func (g *GameState) CheckWin(x, y int) bool {
	color := g.Board[y][x]
	if color == 0 {
		return false
	}

	directions := []struct {
		dx int
		dy int
	}{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		count := 1

		// 正向检查
		for i := 1; i < 5; i++ {
			nx, ny := x+dir.dx*i, y+dir.dy*i
			if !g.IsValidPosition(nx, ny) || g.Board[ny][nx] != color {
				break
			}
			count++
		}

		// 反向检查
		for i := 1; i < 5; i++ {
			nx, ny := x-dir.dx*i, y-dir.dy*i
			if !g.IsValidPosition(nx, ny) || g.Board[ny][nx] != color {
				break
			}
			count++
		}

		if count >= 5 {
			return true
		}
	}

	return false
}

// IsBoardFull 检查棋盘是否已满
func (g *GameState) IsBoardFull() bool {
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if g.Board[y][x] == 0 {
				return false
			}
		}
	}
	return true
}

// GetWinLine 获取获胜连线（用于高亮）
func (g *GameState) GetWinLine(x, y int) []Position {
	if !g.CheckWin(x, y) {
		return nil
	}

	color := g.Board[y][x]
	directions := []struct {
		dx int
		dy int
	}{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		line := []Position{{X: x, Y: y}}

		// 正向
		for i := 1; i < 5; i++ {
			nx, ny := x+dir.dx*i, y+dir.dy*i
			if !g.IsValidPosition(nx, ny) || g.Board[ny][nx] != color {
				break
			}
			line = append(line, Position{X: nx, Y: ny})
		}

		// 反向
		for i := 1; i < 5; i++ {
			nx, ny := x-dir.dx*i, y-dir.dy*i
			if !g.IsValidPosition(nx, ny) || g.Board[ny][nx] != color {
				break
			}
			line = append([]Position{{X: nx, Y: ny}}, line...)
		}

		if len(line) >= 5 {
			return line
		}
	}

	return nil
}

func generateGameID() string {
	// 使用时间戳和随机数生成唯一ID
	return fmt.Sprintf("game_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
}