package handlers

import (
	"net/http"

	"github.com/ai-coding-game-gomoku/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *services.PlayerService
}

func NewPlayerHandler(ps *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: ps,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
}

// Register 注册玩家
func (h *PlayerHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户名"})
		return
	}

	player, err := h.playerService.RegisterPlayer(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": player,
	})
}

// GetPlayer 获取玩家信息
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	userID := c.Param("userId")

	player, err := h.playerService.GetPlayer(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": player,
	})
}

// GetPlayerByUserID 通过userID获取或创建玩家（简化版）
func (h *PlayerHandler) GetPlayerByUserID(c *gin.Context) {
	userID := c.Param("userId")

	player, err := h.playerService.GetPlayer(userID)
	if err != nil {
		// 如果不存在，创建新玩家
		prefix := userID
		if len(userID) > 8 {
			prefix = userID[:8]
		}
		player, err = h.playerService.RegisterPlayer("玩家" + prefix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": player,
	})
}