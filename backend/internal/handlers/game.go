package handlers

import (
	"net/http"

	"github.com/ai-coding-game-gomoku/backend/internal/models"
	"github.com/ai-coding-game-gomoku/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService    *services.GameService
	playerService  *services.PlayerService
}

// NewGameHandler 创建游戏处理器
func NewGameHandler(gs *services.GameService, ps *services.PlayerService) *GameHandler {
	return &GameHandler{
		gameService:   gs,
		playerService: ps,
	}
}

// StartGameRequest 开始游戏请求
type StartGameRequest struct {
	GameMode  string `json:"gameMode" binding:"required"`  // "ai" 或 "pvp"
	Difficulty string `json:"difficulty"`                  // "easy", "medium", "hard"
	UserID    string `json:"userId"`                       // 玩家ID
}

// StartGame 开始游戏
func (h *GameHandler) StartGame(c *gin.Context) {
	var req StartGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供游戏模式"})
		return
	}

	gameMode := models.GameMode(req.GameMode)
	difficulty := models.Difficulty(req.Difficulty)

	if gameMode != models.ModeAI && gameMode != models.ModePVP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的游戏模式"})
		return
	}

	if gameMode == models.ModeAI {
		if difficulty == "" {
			difficulty = models.DifficultyEasy
		}
		if difficulty != models.DifficultyEasy && difficulty != models.DifficultyMedium && difficulty != models.DifficultyHard {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的难度级别"})
			return
		}
	}

	game := h.gameService.CreateGame(gameMode, difficulty)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}

// MoveRequest 落子请求
type MoveRequest struct {
	GameID string `json:"gameId" binding:"required"`
	X      int    `json:"x" binding:"required"`
	Y      int    `json:"y" binding:"required"`
}

// Move 玩家落子
func (h *GameHandler) Move(c *gin.Context) {
	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供落子位置"})
		return
	}

	game, err := h.gameService.PlayerMove(req.GameID, req.X, req.Y)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}

// EndGameRequest 结束游戏请求
type EndGameRequest struct {
	GameID string `json:"gameId" binding:"required"`
	UserID string `json:"userId"` // 玩家ID，用于更新经验值
}

// EndGame 结束游戏
func (h *GameHandler) EndGame(c *gin.Context) {
	var req EndGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供游戏ID"})
		return
	}

	game, err := h.gameService.GetGame(req.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "游戏不存在"})
		return
	}

	response := gin.H{
		"code":           0,
		"gameId":         game.GameID,
		"status":         game.Status,
		"winner":         game.Winner,
		"experienceGain": 0,
	}

	// 人机对战模式更新经验值
	if game.GameMode == models.ModeAI && req.UserID != "" && game.Status == models.StatusWon {
		isWin := game.Winner == models.Black // 玩家执黑棋
		isDraw := game.Status == models.StatusDraw

		difficulty := string(game.Difficulty)
		expGain := h.playerService.GetExperienceGain(difficulty, isWin, isDraw)

		if isWin {
			_, err := h.playerService.UpdatePlayerExperience(req.UserID, expGain, true, false)
			if err == nil {
				response["experienceGain"] = expGain
				response["newRank"] = true
			}
		} else if isDraw {
			_, err := h.playerService.UpdatePlayerExperience(req.UserID, expGain, false, true)
			if err == nil {
				response["experienceGain"] = expGain
			}
		} else {
			// 失败也获得参与奖
			_, err := h.playerService.UpdatePlayerExperience(req.UserID, expGain, false, false)
			if err == nil {
				response["experienceGain"] = expGain
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetGame 获取游戏状态
func (h *GameHandler) GetGame(c *gin.Context) {
	gameID := c.Param("gameId")

	game, err := h.gameService.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "游戏不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}

// Surrender 认输
func (h *GameHandler) Surrender(c *gin.Context) {
	var req struct {
		GameID string `json:"gameId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供游戏ID"})
		return
	}

	game, err := h.gameService.Surrender(req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}

// Restart 重新开始
func (h *GameHandler) Restart(c *gin.Context) {
	var req struct {
		GameID    string `json:"gameId" binding:"required"`
		Difficulty string `json:"difficulty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供游戏ID"})
		return
	}

	difficulty := models.Difficulty(req.Difficulty)
	if difficulty == "" {
		difficulty = models.DifficultyEasy
	}

	game, err := h.gameService.RestartGame(req.GameID, difficulty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}

// Undo 悔棋
func (h *GameHandler) Undo(c *gin.Context) {
	var req struct {
		GameID string `json:"gameId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供游戏ID"})
		return
	}

	game, err := h.gameService.UndoMove(req.GameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": game,
	})
}