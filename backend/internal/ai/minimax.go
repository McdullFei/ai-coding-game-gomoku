package ai

import (
	"math"
)

// 棋型分值
const (
	ScoreFive   int = 100000 // 连五
	ScoreFour   int = 10000  // 活四/冲四
	ScoreThree  int = 1000   // 活三
	ScoreTwo    int = 100    // 活二
	ScoreOne    int = 10     // 眠三/眠二
)

// minimax Minimax算法
func minimax(board [][]int, player int, depth int, alpha, beta int, isMaximizing bool) *Position {
	if depth == 0 {
		score := evaluateBoard(board, player)
		return &Position{X: -1, Y: -1}
	}

	validMoves := getValidMoves(board)
	if len(validMoves) == 0 {
		return nil
	}

	var bestMove *Position

	if isMaximizing {
		maxScore := -math.MaxInt

		for _, move := range validMoves {
			board[move.Y][move.X] = player

			var childScore int
			if checkWin(board, move.X, move.Y, player) {
				childScore = ScoreFive
			} else {
				childScore = -evaluateBoard(board, 3-player)
			}

			board[move.Y][move.X] = 0

			if childScore > maxScore {
				maxScore = childScore
				bestMove = &move
			}

			if maxScore > alpha {
				alpha = maxScore
			}
			if beta <= alpha {
				break
			}
		}
	} else {
		minScore := math.MaxInt

		for _, move := range validMoves {
			board[move.Y][move.X] = 3 - player

			var childScore int
			if checkWin(board, move.X, move.Y, 3-player) {
				childScore = -ScoreFive
			} else {
				childScore = evaluateBoard(board, player)
			}

			board[move.Y][move.X] = 0

			if childScore < minScore {
				minScore = childScore
				bestMove = &move
			}

			if minScore < beta {
				beta = minScore
			}
			if beta <= alpha {
				break
			}
		}
	}

	return bestMove
}

// evaluateBoard 评估棋盘局面
func evaluateBoard(board [][]int, player int) int {
	score := 0
	opponent := 3 - player

	// 评估所有连子情况
	for y := 0; y < 15; y++ {
		for x := 0; x < 15; x++ {
			if board[y][x] == 0 {
				// 评估这个空位对双方的潜力
				score += evaluatePosition(board, x, y, player)
				score -= evaluatePosition(board, x, y, opponent)
			}
		}
	}

	return score
}

// evaluatePosition 评估单个位置
func evaluatePosition(board [][]int, x, y, player int) int {
	score := 0
	directions := []struct{ dx, dy int }{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		line := getLine(board, x, y, dir.dx, dir.dy, player)
		patternScore := evaluatePattern(line)
		score += patternScore
	}

	return score
}

// getLine 获取指定方向上的棋型
func getLine(board [][]int, x, y, dx, dy, player int) []int {
	line := make([]int, 0)

	// 向前4格
	for i := -4; i <= 4; i++ {
		nx, ny := x+i*dx, y+i*dy
		if nx >= 0 && nx < 15 && ny >= 0 && ny < 15 {
			line = append(line, board[ny][nx])
		}
	}

	return line
}

// evaluatePattern 评估棋型
func evaluatePattern(line []int) int {
	if len(line) < 5 {
		return 0
	}

	playerCount := 0
	opponentCount := 0
	emptyCount := 0

	for _, v := range line {
		if v == 1 {
			playerCount++
		} else if v == 2 {
			opponentCount++
		} else {
			emptyCount++
		}
	}

	// 连五
	if playerCount >= 5 {
		return ScoreFive
	}

	// 活四
	if playerCount == 4 && emptyCount >= 1 {
		return ScoreFour
	}

	// 冲四
	if playerCount == 3 && emptyCount >= 2 {
		midEmpty := 0
		for _, v := range line {
			if v == 0 {
				midEmpty++
			}
		}
		if midEmpty >= 2 {
			return ScoreFour
		}
	}

	// 活三
	if playerCount == 3 && emptyCount >= 2 {
		return ScoreThree
	}

	// 活二
	if playerCount == 2 && emptyCount >= 3 {
		return ScoreTwo
	}

	// 眠三
	if playerCount == 2 && emptyCount >= 1 {
		return ScoreOne
	}

	return 0
}