package main

import (
	"log"

	"github.com/ai-coding-game-gomoku/backend/internal/config"
	"github.com/ai-coding-game-gomoku/backend/internal/handlers"
	"github.com/ai-coding-game-gomoku/backend/internal/repository"
	"github.com/ai-coding-game-gomoku/backend/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := repository.InitDB(cfg); err != nil {
		log.Printf("数据库初始化失败: %v，将使用内存存储", err)
	}

	// 运行数据库迁移
	if err := repository.RunMigrations(); err != nil {
		log.Printf("数据库迁移失败: %v", err)
	}

	// 初始化服务
	playerService := services.NewPlayerService()
	gameService := services.NewGameService(playerService)

	// 初始化处理器
	playerHandler := handlers.NewPlayerHandler(playerService)
	gameHandler := handlers.NewGameHandler(gameService, playerService)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由
	api := r.Group("/api")
	{
		// 玩家相关
		api.POST("/player/register", playerHandler.Register)
		api.GET("/player/:userId", playerHandler.GetPlayer)
		api.GET("/player/get/:userId", playerHandler.GetPlayerByUserID)

		// 游戏相关
		api.POST("/game/start", gameHandler.StartGame)
		api.POST("/game/move", gameHandler.Move)
		api.POST("/game/end", gameHandler.EndGame)
		api.GET("/game/:gameId", gameHandler.GetGame)
		api.POST("/game/surrender", gameHandler.Surrender)
		api.POST("/game/restart", gameHandler.Restart)
		api.POST("/game/undo", gameHandler.Undo)
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}