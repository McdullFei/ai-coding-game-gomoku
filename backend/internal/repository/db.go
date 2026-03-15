package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ai-coding-game-gomoku/backend/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %w", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("数据库连接失败，将使用内存存储: %v", err)
		return nil
	}

	log.Println("数据库连接成功")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func RunMigrations() error {
	if DB == nil {
		log.Println("数据库未连接，跳过迁移")
		return nil
	}

	migrations := []string{
		`CREATE TABLE IF NOT EXISTS player_ranks (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(64) UNIQUE NOT NULL,
			username VARCHAR(64) NOT NULL,
			rank_level INTEGER DEFAULT 1,
			rank_segment VARCHAR(16) DEFAULT '新手',
			experience INTEGER DEFAULT 0,
			total_wins INTEGER DEFAULT 0,
			total_losses INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS game_records (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL,
			game_mode VARCHAR(16) NOT NULL,
			difficulty VARCHAR(16),
			result VARCHAR(16) NOT NULL,
			experience_gained INTEGER DEFAULT 0,
			played_at TIMESTAMP DEFAULT NOW()
		)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return fmt.Errorf("数据库迁移失败: %w", err)
		}
	}

	log.Println("数据库迁移完成")
	return nil
}