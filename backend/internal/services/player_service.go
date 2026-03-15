package services

import (
	"database/sql"
	"fmt"

	"github.com/ai-coding-game-gomoku/backend/internal/models"
	"github.com/ai-coding-game-gomoku/backend/internal/repository"
	"github.com/google/uuid"
)

// PlayerService 玩家服务
type PlayerService struct{}

// NewPlayerService 创建玩家服务
func NewPlayerService() *PlayerService {
	return &PlayerService{}
}

// RegisterPlayer 注册新玩家
func (s *PlayerService) RegisterPlayer(username string) (*models.Player, error) {
	userID := uuid.New().String()

	player := &models.Player{
		UserID:      userID,
		Username:    username,
		RankLevel:   1,
		RankSegment: "新手",
		Experience:  0,
		TotalWins:   0,
		TotalLosses: 0,
	}

	// 如果数据库可用，存入数据库
	if repository.DB != nil {
		err := repository.DB.QueryRow(`
			INSERT INTO player_ranks (user_id, username, rank_level, rank_segment, experience, total_wins, total_losses)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, created_at, updated_at`,
			player.UserID, player.Username, player.RankLevel, player.RankSegment,
			player.Experience, player.TotalWins, player.TotalLosses,
		).Scan(&player.ID, &player.CreatedAt, &player.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("注册玩家失败: %w", err)
		}
		return player, nil
	}

	// 数据库不可用时使用内存存储
	player.ID = len(players) + 1
	player.CreatedAt = "now"
	player.UpdatedAt = "now"
	players[userID] = player

	return player, nil
}

// GetPlayer 获取玩家信息
func (s *PlayerService) GetPlayer(userID string) (*models.Player, error) {
	if repository.DB != nil {
		player := &models.Player{}
		err := repository.DB.QueryRow(`
			SELECT id, user_id, username, rank_level, rank_segment, experience, total_wins, total_losses, created_at, updated_at
			FROM player_ranks WHERE user_id = $1`, userID).Scan(
			&player.ID, &player.UserID, &player.Username, &player.RankLevel,
			&player.RankSegment, &player.Experience, &player.TotalWins,
			&player.TotalLosses, &player.CreatedAt, &player.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("玩家不存在")
		}
		if err != nil {
			return nil, fmt.Errorf("获取玩家失败: %w", err)
		}
		return player, nil
	}

	// 内存存储
	if player, ok := players[userID]; ok {
		return player, nil
	}
	return nil, fmt.Errorf("玩家不存在")
}

// UpdatePlayerExperience 更新玩家经验值
func (s *PlayerService) UpdatePlayerExperience(userID string, expGain int, isWin bool, isDraw bool) (*models.Player, error) {
	player, err := s.GetPlayer(userID)
	if err != nil {
		return nil, err
	}

	player.Experience += expGain

	// 计算新等级
	newLevel := calculateLevel(player.Experience)
	player.RankLevel = newLevel
	player.RankSegment = models.GetSegmentByLevel(newLevel)

	if isWin {
		player.TotalWins++
	} else if !isDraw {
		player.TotalLosses++
	}

	// 更新数据库
	if repository.DB != nil {
		_, err = repository.DB.Exec(`
			UPDATE player_ranks
			SET rank_level = $1, rank_segment = $2, experience = $3, total_wins = $4, total_losses = $5, updated_at = NOW()
			WHERE user_id = $6`,
			player.RankLevel, player.RankSegment, player.Experience,
			player.TotalWins, player.TotalLosses, userID,
		)
		if err != nil {
			return nil, fmt.Errorf("更新玩家经验失败: %w", err)
		}
	} else {
		players[userID] = player
	}

	return player, nil
}

// calculateLevel 根据经验值计算等级
func calculateLevel(exp int) int {
	// 简化版：每100经验升一级
	level := exp/100 + 1
	if level < 1 {
		return 1
	}
	if level > 100 {
		return 100
	}
	return level
}

// GetExperienceGain 获取经验值增量
func (s *PlayerService) GetExperienceGain(difficulty string, isWin bool, isDraw bool) int {
	if isDraw {
		return models.ExperienceConfig["draw"]
	}
	if !isWin {
		return models.ExperienceConfig["lose"]
	}

	switch difficulty {
	case "easy":
		return models.ExperienceConfig["win_easy"]
	case "medium":
		return models.ExperienceConfig["win_medium"]
	case "hard":
		return models.ExperienceConfig["win_hard"]
	default:
		return 0
	}
}

// 内存存储（数据库不可用时）
var players = make(map[string]*models.Player)