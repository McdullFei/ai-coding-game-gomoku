package services

import (
	"github.com/ai-coding-game-gomoku/backend/internal/ai"
	"github.com/ai-coding-game-gomoku/backend/internal/models"
)

// GameService 游戏服务
type GameService struct {
	playerService *PlayerService
}

// NewGameService 创建游戏服务
func NewGameService(ps *PlayerService) *GameService {
	return &GameService{
		playerService: ps,
	}
}

// GameStore 游戏存储（内存）
var GameStore = make(map[string]*models.GameState)

// CreateGame 创建新游戏
func (s *GameService) CreateGame(gameMode models.GameMode, difficulty models.Difficulty) *models.GameState {
	game := models.NewGameState(gameMode, difficulty)
	GameStore[game.GameID] = game
	return game
}

// GetGame 获取游戏状态
func (s *GameService) GetGame(gameID string) (*models.GameState, error) {
	game, ok := GameStore[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}
	return game, nil
}

// PlayerMove 玩家落子
func (s *GameService) PlayerMove(gameID string, x, y int) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.StatusPlaying {
		return nil, ErrGameEnded
	}

	// 验证落子
	if !game.IsEmpty(x, y) {
		return nil, ErrInvalidMove
	}

	// 落子
	game.PlacePiece(x, y, game.CurrentPlayer)

	// 检查胜利
	if game.CheckWin(x, y) {
		game.Status = models.StatusWon
		game.Winner = game.CurrentPlayer
		return game, nil
	}

	// 检查平局
	if game.IsBoardFull() {
		game.Status = models.StatusDraw
		return game, nil
	}

	// 切换玩家
	game.SwitchPlayer()

	// 如果是人机模式，且当前是AI（白棋），则AI落子
	if game.GameMode == models.ModeAI && game.CurrentPlayer == models.White {
		aiMove := ai.GetAIMove(game.Board, string(game.Difficulty))
		if aiMove != nil {
			game.PlacePiece(aiMove.X, aiMove.Y, models.White)

			if game.CheckWin(aiMove.X, aiMove.Y) {
				game.Status = models.StatusWon
				game.Winner = models.White
				return game, nil
			}

			if game.IsBoardFull() {
				game.Status = models.StatusDraw
				return game, nil
			}

			game.SwitchPlayer()
		}
	}

	return game, nil
}

// AIGo 电脑落子（返回AI落子位置）
func (s *GameService) AIGo(gameID string) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.GameMode != models.ModeAI {
		return nil, ErrInvalidGameMode
	}

	if game.CurrentPlayer != models.White {
		return nil, ErrNotAITurn
	}

	aiMove := ai.GetAIMove(game.Board, string(game.Difficulty))
	if aiMove == nil {
		return nil, ErrNoValidMove
	}

	game.PlacePiece(aiMove.X, aiMove.Y, models.White)

	if game.CheckWin(aiMove.X, aiMove.Y) {
		game.Status = models.StatusWon
		game.Winner = models.White
		return game, nil
	}

	if game.IsBoardFull() {
		game.Status = models.StatusDraw
		return game, nil
	}

	game.SwitchPlayer()
	return game, nil
}

// Surrender 认输
func (s *GameService) Surrender(gameID string) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.StatusPlaying {
		return nil, ErrGameEnded
	}

	// 当前玩家认输，另一方获胜
	if game.CurrentPlayer == models.Black {
		game.Winner = models.White
	} else {
		game.Winner = models.Black
	}
	game.Status = models.StatusWon

	return game, nil
}

// RestartGame 重新开始游戏
func (s *GameService) RestartGame(gameID string, difficulty models.Difficulty) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	newGame := models.NewGameState(game.GameMode, difficulty)
	GameStore[game.GameID] = newGame

	return newGame, nil
}

// UndoMove 悔棋（删除最后一步）
func (s *GameService) UndoMove(gameID string) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if len(game.MoveHistory) == 0 {
		return nil, ErrNoMoveToUndo
	}

	// 移除最后一步
	lastMove := game.MoveHistory[len(game.MoveHistory)-1]
	game.Board[lastMove.Y][lastMove.X] = 0
	game.MoveHistory = game.MoveHistory[:len(game.MoveHistory)-1]
	game.LastMove = nil
	if len(game.MoveHistory) > 0 {
		game.LastMove = &game.MoveHistory[len(game.MoveHistory)-1]
	}

	// 切换回上一个玩家
	game.SwitchPlayer()
	game.Status = models.StatusPlaying
	game.Winner = 0

	return game, nil
}

// 错误定义
var (
	ErrGameNotFound   = &GameError{"游戏不存在"}
	ErrGameEnded      = &GameError{"游戏已结束"}
	ErrInvalidMove    = &GameError{"无效的落子位置"}
	ErrInvalidGameMode = &GameError{"无效的游戏模式"}
	ErrNotAITurn      = &GameError{"当前不是AI回合"}
	ErrNoValidMove    = &GameError{"没有有效的落子位置"}
	ErrNoMoveToUndo   = &GameError{"没有可悔的棋"}
)

type GameError struct {
	Message string
}

func (e *GameError) Error() string {
	return e.Message
}