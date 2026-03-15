package ai

import (
	"math/rand"
	"time"
)

// Position 位置
type Position struct {
	X int
	Y int
}

// GetAIMove 获取AI落子
func GetAIMove(board [][]int, difficulty string) *Position {
	switch difficulty {
	case "easy":
		return getEasyMove(board)
	case "medium":
		return getMediumMove(board)
	case "hard":
		return getHardMove(board)
	default:
		return getEasyMove(board)
	}
}

// getEasyMove 简单难度：随机落子 + 基础防守
func getEasyMove(board [][]int) *Position {
	// 1. 检查是否可以赢（直接获胜）
	if win := findWinningMove(board, 2); win != nil {
		return win
	}

	// 2. 检查是否需要防守（对手即将获胜）
	if block := findWinningMove(board, 1); block != nil {
		return block
	}

	// 3. 随机落子
	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	return &validMoves[rand.Intn(len(validMoves))]
}

// getMediumMove 中等难度：Minimax算法
func getMediumMove(board [][]int) *Position {
	// 先检查必赢/必防
	if win := findWinningMove(board, 2); win != nil {
		return win
	}
	if block := findWinningMove(board, 1); block != nil {
		return block
	}

	// 使用简化的Minimax
	move := minimax(board, 2, 3, -1000000, 1000000, false)
	if move != nil {
		return move
	}

	// 后备：中心优先随机
	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return nil
	}

	// 优先选择靠近中心的点
	center := 7
	bestMove := &validMoves[0]
	bestDist := 1000
	for _, m := range validMoves {
		dist := abs(m.X-center) + abs(m.Y-center)
		if dist < bestDist {
			bestDist = dist
			bestMove = &m
		}
	}

	return bestMove
}

// getHardMove 困难难度：深度Minimax + Alpha-Beta剪枝
func getHardMove(board [][]int) *Position {
	// 先检查必赢/必防
	if win := findWinningMove(board, 2); win != nil {
		return win
	}
	if block := findWinningMove(board, 1); block != nil {
		return block
	}

	// 使用深度Minimax
	move := minimax(board, 2, 4, -1000000, 1000000, false)
	if move != nil {
		return move
	}

	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return nil
	}

	// 优先选择靠近中心的点
	center := 7
	bestMove := &validMoves[0]
	bestDist := 1000
	for _, m := range validMoves {
		dist := abs(m.X-center) + abs(m.Y-center)
		if dist < bestDist {
			bestDist = dist
			bestMove = &m
		}
	}

	return bestMove
}

// findWinningMove 寻找获胜落子点
func findWinningMove(board [][]int, player int) *Position {
	validMoves := getValidMoves(board)
	for _, move := range validMoves {
		// 模拟落子
		board[move.Y][move.X] = player
		if checkWin(board, move.X, move.Y, player) {
			board[move.Y][move.X] = 0 // 恢复
			return &move
		}
		board[move.Y][move.X] = 0 // 恢复
	}
	return nil
}

// getValidMoves 获取所有有效落子点
func getValidMoves(board [][]int) []Position {
	moves := make([]Position, 0)
	visited := make(map[string]bool)

	// 搜索落子点周围2格内有棋子的位置
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] != 0 {
				// 检查周围
				for dy := -2; dy <= 2; dy++ {
					for dx := -2; dx <= 2; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < 15 && ny >= 0 && ny < 15 && board[ny][nx] == 0 {
							key := string(rune(nx)) + "," + string(rune(ny))
							if !visited[key] {
								visited[key] = true
								moves = append(moves, Position{X: nx, Y: ny})
							}
						}
					}
				}
			}
		}
	}

	// 如果棋盘为空，下中心点
	if len(moves) == 0 {
		moves = append(moves, Position{X: 7, Y: 7})
	}

	return moves
}

// checkWin 检查是否获胜
func checkWin(board [][]int, x, y, player int) bool {
	directions := []struct{ dx, dy int }{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		count := 1

		for i := 1; i < 5; i++ {
			nx, ny := x+dir.dx*i, y+dir.dy*i
			if !inBoard(nx, ny) || board[ny][nx] != player {
				break
			}
			count++
		}

		for i := 1; i < 5; i++ {
			nx, ny := x-dir.dx*i, y-dir.dy*i
			if !inBoard(nx, ny) || board[ny][nx] != player {
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

func inBoard(x, y int) bool {
	return x >= 0 && x < 15 && y >= 0 && y < 15
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}