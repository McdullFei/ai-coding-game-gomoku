-- 五子棋游戏数据库初始化脚本

-- 创建数据库（如果不存在）
-- CREATE DATABASE gomoku;

-- 玩家等级表
CREATE TABLE IF NOT EXISTS player_ranks (
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
);

-- 对局记录表
CREATE TABLE IF NOT EXISTS game_records (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL,
    game_mode VARCHAR(16) NOT NULL,
    difficulty VARCHAR(16),
    result VARCHAR(16) NOT NULL,
    experience_gained INTEGER DEFAULT 0,
    played_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_player_user_id ON player_ranks(user_id);
CREATE INDEX IF NOT EXISTS idx_game_records_user_id ON game_records(user_id);

-- 插入测试数据（可选）
-- INSERT INTO player_ranks (user_id, username, rank_level, rank_segment, experience)
-- VALUES ('test_user', '测试玩家', 1, '新手', 0)
-- ON CONFLICT (user_id) DO NOTHING;