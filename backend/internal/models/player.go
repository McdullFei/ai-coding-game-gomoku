package models

// Player 玩家信息
type Player struct {
	ID            int    `json:"id"`
	UserID        string `json:"userId"`
	Username      string `json:"username"`
	RankLevel     int    `json:"rankLevel"`
	RankSegment   string `json:"rankSegment"`
	Experience    int    `json:"experience"`
	TotalWins     int    `json:"totalWins"`
	TotalLosses   int    `json:"totalLosses"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// RankSegmentData 段位数据
var RankSegments = []struct {
	Segment      string
	MinLevel     int
	ExpThreshold int
}{
	{"新手", 1, 0},
	{"初级", 11, 100},
	{"中级", 21, 500},
	{"高级", 31, 1500},
	{"专家", 41, 4000},
	{"大师", 51, 10000},
}

// GetSegmentByLevel 根据等级获取段位
func GetSegmentByLevel(level int) string {
	for i := len(RankSegments) - 1; i >= 0; i-- {
		if level >= RankSegments[i].MinLevel {
			return RankSegments[i].Segment
		}
	}
	return "新手"
}

// GetExpToNextLevel 获取升级到下一级所需经验
func GetExpToNextLevel(level int) int {
	for i := 0; i < len(RankSegments)-1; i++ {
		if level >= RankSegments[i].MinLevel && level < RankSegments[i+1].MinLevel {
			return RankSegments[i+1].ExpThreshold
		}
	}
	// 大师级别后每级需要10000经验
	return 10000
}

// ExperienceConfig 经验值配置
var ExperienceConfig = map[string]int{
	"win_easy":   10,
	"win_medium": 25,
	"win_hard":   50,
	"lose":       5,
	"draw":       3,
}