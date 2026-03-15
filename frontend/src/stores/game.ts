import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GameState, GameMode, Difficulty, PlayerColor, Position } from '../types'
import { gameApi } from '../utils/api'
import { usePlayerStore } from './player'

export const useGameStore = defineStore('game', () => {
  const game = ref<GameState | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const winLine = ref<Position[]>([])
  const useLocalMode = ref(true) // 使用本地模拟模式

  // 判断游戏是否进行中
  const isPlaying = computed(() => game.value?.status === 'playing')

  // 获取当前玩家名称
  const currentPlayerName = computed(() => {
    if (!game.value) return ''
    return game.value.currentPlayer === 1 ? '黑方' : '白方'
  })

  // 玩家执黑棋（先行）
  const isPlayerTurn = computed(() => {
    if (!game.value || game.value.gameMode === 'pvp') {
      return game.value?.currentPlayer === 1
    }
    // 人机对战时，玩家始终执黑棋
    return game.value.currentPlayer === 1
  })

  // 棋型分值
  const SCORE_FIVE = 100000   // 连五
  const SCORE_FOUR = 10000    // 活四/冲四
  const SCORE_THREE = 1000    // 活三
  const SCORE_TWO = 100       // 活二
  const SCORE_ONE = 10        // 眠三/眠二

  // 棋型类型
  type PatternType = 'five' | 'four' | '活三' | '冲四' | '活二' | '眠三' | '眠二' | 'none'

  // 获取棋型名称
  function getPatternType(count: number, openEnds: number): PatternType {
    if (count >= 5) return 'five'
    if (count === 4 && openEnds > 0) return 'four'
    if (count === 3 && openEnds === 2) return '活三'
    if (count === 3 && openEnds === 1) return '冲四'
    if (count === 2 && openEnds === 2) return '活二'
    if (count === 2 && openEnds === 1) return '眠二'
    if (count === 1 && openEnds === 2) return '活二'
    return 'none'
  }

  // 评估棋型
  function evaluatePattern(count: number, openEnds: number): number {
    if (count >= 5) return SCORE_FIVE
    if (count === 4 && openEnds > 0) return SCORE_FOUR
    if (count === 3 && openEnds === 2) return SCORE_THREE
    if (count === 3 && openEnds === 1) return SCORE_ONE
    if (count === 2 && openEnds === 2) return SCORE_TWO
    if (count === 2 && openEnds === 1) return SCORE_ONE
    return 0
  }

  // 统计线上的棋子数量和两端空位
  function countLine(board: number[][], x: number, y: number, dx: number, dy: number, player: number): { count: number; openEnds: number } {
    let count = 1
    let openEnds = 0

    // 正向检查
    let i = 1
    while (true) {
      const nx = x + dx * i
      const ny = y + dy * i
      if (nx < 0 || nx >= 15 || ny < 0 || ny >= 15) break
      if (board[ny][nx] === player) {
        count++
      } else if (board[ny][nx] === 0) {
        openEnds++
        break
      } else {
        break
      }
      i++
    }

    // 反向检查
    i = 1
    while (true) {
      const nx = x - dx * i
      const ny = y - dy * i
      if (nx < 0 || nx >= 15 || ny < 0 || ny >= 15) break
      if (board[ny][nx] === player) {
        count++
      } else if (board[ny][nx] === 0) {
        openEnds++
        break
      } else {
        break
      }
      i++
    }

    return { count, openEnds }
  }

  // 检查已有的棋型（不落子）
  function checkExistingPattern(board: number[][], x: number, y: number, player: number): PatternType | null {
    const directions = [
      { dx: 1, dy: 0 },
      { dx: 0, dy: 1 },
      { dx: 1, dy: 1 },
      { dx: 1, dy: -1 }
    ]

    for (const dir of directions) {
      const { count, openEnds } = countLine(board, x, y, dir.dx, dir.dy, player)
      const pattern = getPatternType(count, openEnds)
      if (pattern === 'four' || pattern === '活三') {
        return pattern
      }
    }
    return null
  }

  // 寻找需要防守的关键位置
  function findDefensiveMoves(board: number[][], player: number): Position[] {
    const defensiveMoves: Position[] = []
    const directions = [
      { dx: 1, dy: 0 },
      { dx: 0, dy: 1 },
      { dx: 1, dy: 1 },
      { dx: 1, dy: -1 }
    ]

    // 遍历每个方向
    for (let y = 0; y < 15; y++) {
      for (let x = 0; x < 15; x++) {
        if (board[y][x] === player) {
          // 检查这个棋子所在的线
          for (const dir of directions) {
            // 向一个方向统计连续的棋子
            let count = 1
            let endX = x, endY = y

            // 正向
            let i = 1
            while (true) {
              const nx = x + dir.dx * i
              const ny = y + dir.dy * i
              if (nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || board[ny][nx] !== player) break
              count++
              endX = nx
              endY = ny
              i++
            }

            // 检查正向的端点是否为空
            let fx = endX + dir.dx
            let fy = endY + dir.dy
            if (fx >= 0 && fx < 15 && fy >= 0 && fy < 15 && board[fy][fx] === 0) {
              if (count === 4) {
                // 活四！必须防守
                defensiveMoves.push({ x: fx, y: fy })
              } else if (count === 3) {
                // 检查反向是否有空格
                let bx = x - dir.dx
                let by = y - dir.dy
                if (bx >= 0 && bx < 15 && by >= 0 && by < 15 && board[by][bx] === 0) {
                  // 活三
                  defensiveMoves.push({ x: fx, y: fy })
                  defensiveMoves.push({ x: bx, y: by })
                } else {
                  // 眠三，只防守一端
                  defensiveMoves.push({ x: fx, y: fy })
                }
              }
            }

            // 检查反向的端点
            let bx = x - dir.dx
            let by = y - dir.dy
            if (bx >= 0 && bx < 15 && by >= 0 && by < 15 && board[by][bx] === 0) {
              if (count === 4) {
                defensiveMoves.push({ x: bx, y: by })
              } else if (count === 3) {
                defensiveMoves.push({ x: bx, y: by })
              }
            }
          }
        }
      }
    }

    return defensiveMoves
  }

  // 获取棋盘上所有威胁位置（落子后形成的棋型）
  function findThreats(board: number[][], player: number): { x: number, y: number, type: PatternType }[] {
    const threats: { x: number, y: number, type: PatternType }[] = []
    const directions = [
      { dx: 1, dy: 0 },
      { dx: 0, dy: 1 },
      { dx: 1, dy: 1 },
      { dx: 1, dy: -1 }
    ]

    // 遍历所有空位
    for (let y = 0; y < 15; y++) {
      for (let x = 0; x < 15; x++) {
        if (board[y][x] === 0) {
          for (const dir of directions) {
            const { count, openEnds } = countLine(board, x, y, dir.dx, dir.dy, player)
            const pattern = getPatternType(count, openEnds)

            if (pattern === 'four' || pattern === '活三' || pattern === '冲四') {
              threats.push({ x, y, type: pattern })
            }
          }
        }
      }
    }

    // 排序
    const priority: Record<PatternType, number> = {
      'five': 4, 'four': 3, '冲四': 3, '活三': 2,
      '活二': 1, '眠三': 1, '眠二': 0, 'none': 0
    }
    threats.sort((a, b) => priority[b.type] - priority[a.type])

    return threats
  }

  // 检查玩家是否在某位置落子后获胜
  function wouldWin(board: number[][], x: number, y: number, player: number): boolean {
    if (x < 0 || x >= 15 || y < 0 || y >= 15 || board[y][x] !== 0) return false

    board[y][x] = player
    const win = checkWinLocal(board, x, y)
    board[y][x] = 0
    return win
  }

  // 评估单个位置
  function evaluatePosition(board: number[][], x: number, y: number, player: number): number {
    if (board[y][x] !== 0) return -999999

    let score = 0
    const directions = [
      { dx: 1, dy: 0 },
      { dx: 0, dy: 1 },
      { dx: 1, dy: 1 },
      { dx: 1, dy: -1 }
    ]

    for (const dir of directions) {
      const { count, openEnds } = countLine(board, x, y, dir.dx, dir.dy, player)
      score += evaluatePattern(count, openEnds)

      const oppResult = countLine(board, x, y, dir.dx, dir.dy, 3 - player)
      score += evaluatePattern(oppResult.count, oppResult.openEnds)
    }

    const centerDist = Math.abs(x - 7) + Math.abs(y - 7)
    score += (14 - centerDist) * 2

    return score
  }

  // 本地模式下的AI
  function getLocalAIMove(board: number[][], difficulty: Difficulty): Position | null {
    // 获取有效落子点
    const validMoves: Position[] = []
    const visited = new Set<string>()

    for (let y = 0; y < 15; y++) {
      for (let x = 0; x < 15; x++) {
        if (board[y][x] !== 0) {
          for (let dy = -2; dy <= 2; dy++) {
            for (let dx = -2; dx <= 2; dx++) {
              const nx = x + dx
              const ny = y + dy
              if (nx >= 0 && nx < 15 && ny >= 0 && ny < 15 && board[ny][nx] === 0) {
                const key = `${nx},${ny}`
                if (!visited.has(key)) {
                  visited.add(key)
                  validMoves.push({ x: nx, y: ny })
                }
              }
            }
          }
        }
      }
    }

    // 棋盘为空，下中心
    if (validMoves.length === 0) {
      return { x: 7, y: 7 }
    }

    // AI执白棋(2)，玩家执黑棋(1)
    const aiPlayer = 2
    const humanPlayer = 1

    // ====== 1. AI能赢直接赢 ======
    for (const move of validMoves) {
      if (wouldWin(board, move.x, move.y, aiPlayer)) {
        return move
      }
    }

    // ====== 2. 防守：先检查已有的棋型（不落子）======
    const defensive = findDefensiveMoves(board, humanPlayer)

    // 优先防守活四（已有4子）
    for (const move of defensive) {
      // 这个位置防守活四
      return move
    }

    // ====== 3. 防守：检查落子后形成的威胁 ======
    const threats = findThreats(board, humanPlayer)

    // 冲四必挡
    for (const threat of threats) {
      if (threat.type === 'four' || threat.type === '冲四') {
        return { x: threat.x, y: threat.y }
      }
    }

    // 活三必应
    for (const threat of threats) {
      if (threat.type === '活三') {
        if (difficulty === 'easy' && Math.random() < 0.7) {
          return { x: threat.x, y: threat.y }
        }
        if (difficulty === 'medium' || difficulty === 'hard') {
          return { x: threat.x, y: threat.y }
        }
      }
    }

    // ====== 4. 根据难度选择落子 ======
    if (difficulty === 'easy') {
      const center = 7
      const candidates = validMoves.filter(m => {
        const dist = Math.abs(m.x - center) + Math.abs(m.y - center)
        return dist <= 8
      })
      const pool = candidates.length > 0 ? candidates : validMoves
      return pool[Math.floor(Math.random() * pool.length)]
    }

    if (difficulty === 'medium') {
      let bestScore = -Infinity
      let bestMove = validMoves[0]

      for (const move of validMoves) {
        board[move.y][move.x] = aiPlayer
        const attackScore = evaluatePosition(board, move.x, move.y, aiPlayer)
        board[move.y][move.x] = 0

        board[move.y][move.x] = humanPlayer
        const defenseScore = evaluatePosition(board, move.x, move.y, humanPlayer)
        board[move.y][move.x] = 0

        let totalScore = attackScore * 0.6 + defenseScore * 0.5
        totalScore += Math.random() * 30

        if (totalScore > bestScore) {
          bestScore = totalScore
          bestMove = move
        }
      }
      return bestMove
    }

    // 困难难度
    let bestScore = -Infinity
    let bestMove = validMoves[0]

    for (const move of validMoves) {
      board[move.y][move.x] = aiPlayer
      const attackScore = evaluatePosition(board, move.x, move.y, aiPlayer)
      board[move.y][move.x] = 0

      board[move.y][move.x] = humanPlayer
      const defenseScore = evaluatePosition(board, move.x, move.y, humanPlayer)
      board[move.y][move.x] = 0

      let totalScore = attackScore * 0.4 + defenseScore * 0.7

      if (totalScore > bestScore) {
        bestScore = totalScore
        bestMove = move
      }
    }
    return bestMove
  }

  // 本地模式下的胜负判定
  function checkWinLocal(board: number[][], x: number, y: number): boolean {
    const color = board[y][x]
    if (color === 0) return false

    const directions = [{ dx: 1, dy: 0 }, { dx: 0, dy: 1 }, { dx: 1, dy: 1 }, { dx: 1, dy: -1 }]

    for (const dir of directions) {
      let count = 1

      for (let i = 1; i < 5; i++) {
        const nx = x + dir.dx * i
        const ny = y + dir.dy * i
        if (nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || board[ny][nx] !== color) break
        count++
      }

      for (let i = 1; i < 5; i++) {
        const nx = x - dir.dx * i
        const ny = y - dir.dy * i
        if (nx < 0 || nx >= 15 || ny < 0 || ny >= 15 || board[ny][nx] !== color) break
        count++
      }

      if (count >= 5) return true
    }
    return false
  }

  // 本地落子
  function localMove(x: number, y: number) {
    if (!game.value || !isPlaying.value) return

    // 玩家落子
    game.value.board[y][x] = game.value.currentPlayer
    game.value.lastMove = { x, y }
    game.value.moveHistory.push({ x, y })

    // 检查胜利
    if (checkWinLocal(game.value.board, x, y)) {
      game.value.status = 'won'
      game.value.winner = game.value.currentPlayer
      winLine.value = calculateWinLine(game.value.board, x, y)
      return
    }

    // 检查平局
    let isFull = true
    for (let i = 0; i < 15 && isFull; i++) {
      for (let j = 0; j < 15; j++) {
        if (game.value.board[i][j] === 0) {
          isFull = false
          break
        }
      }
    }
    if (isFull) {
      game.value.status = 'draw'
      return
    }

    // 切换玩家
    game.value.currentPlayer = game.value.currentPlayer === 1 ? 2 : 1

    // 人机对战，AI落子
    if (game.value.gameMode === 'ai' && game.value.currentPlayer === 2) {
      const aiMove = getLocalAIMove(game.value.board, game.value.difficulty || 'easy')
      if (aiMove) {
        game.value.board[aiMove.y][aiMove.x] = 2
        game.value.lastMove = aiMove
        game.value.moveHistory.push(aiMove)

        if (checkWinLocal(game.value.board, aiMove.x, aiMove.y)) {
          game.value.status = 'won'
          game.value.winner = 2
          winLine.value = calculateWinLine(game.value.board, aiMove.x, aiMove.y)
          return
        }

        // 检查平局
        isFull = true
        for (let i = 0; i < 15 && isFull; i++) {
          for (let j = 0; j < 15; j++) {
            if (game.value.board[i][j] === 0) {
              isFull = false
              break
            }
          }
        }
        if (isFull) {
          game.value.status = 'draw'
          return
        }

        // 切换回玩家
        game.value.currentPlayer = 1
      }
    }
  }

  // 落子
  async function makeMove(x: number, y: number) {
    if (!game.value || !isPlaying.value) return
    if (game.value.gameMode === 'ai' && !isPlayerTurn.value) {
      return
    }

    if (useLocalMode.value) {
      localMove(x, y)
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const res = await gameApi.move(game.value.gameId, x, y)
      if (res.data.code === 0 && res.data.data) {
        game.value = res.data.data

        if (game.value.status === 'won' && game.value.lastMove) {
          winLine.value = calculateWinLine(game.value.board, game.value.lastMove.x, game.value.lastMove.y)
        }
      } else {
        error.value = res.data.error || '落子失败'
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || '网络错误'
    } finally {
      isLoading.value = false
    }
  }

  // 开始游戏
  async function startGame(mode: GameMode, difficulty: Difficulty = 'easy') {
    isLoading.value = true
    error.value = null
    winLine.value = []

    if (useLocalMode.value) {
      // 本地模式
      const board: number[][] = []
      for (let i = 0; i < 15; i++) {
        board.push(new Array(15).fill(0))
      }
      game.value = {
        gameId: 'local_' + Date.now(),
        board,
        currentPlayer: 1,
        gameMode: mode,
        difficulty: mode === 'ai' ? difficulty : undefined,
        status: 'playing',
        moveHistory: []
      }
      isLoading.value = false
      return
    }

    const playerStore = usePlayerStore()
    const userId = playerStore.player?.userId

    try {
      const res = await gameApi.start(mode, mode === 'ai' ? difficulty : undefined, userId)
      if (res.data.code === 0 && res.data.data) {
        game.value = res.data.data
      } else {
        error.value = res.data.error || '开始游戏失败'
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || '网络错误'
    } finally {
      isLoading.value = false
    }
  }

  // 重新开始
  async function restartGame() {
    if (!game.value) return

    if (useLocalMode.value) {
      const board: number[][] = []
      for (let i = 0; i < 15; i++) {
        board.push(new Array(15).fill(0))
      }
      game.value = {
        ...game.value,
        board,
        currentPlayer: 1,
        status: 'playing',
        winner: undefined,
        lastMove: undefined,
        moveHistory: []
      }
      winLine.value = []
      return
    }

    isLoading.value = true
    error.value = null
    winLine.value = []

    try {
      const res = await gameApi.restart(game.value.gameId, game.value.difficulty)
      if (res.data.code === 0 && res.data.data) {
        game.value = res.data.data
      } else {
        error.value = res.data.error || '重新开始失败'
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || '网络错误'
    } finally {
      isLoading.value = false
    }
  }

  // 认输
  async function surrender() {
    if (!game.value || !isPlaying.value) return

    if (useLocalMode.value) {
      game.value.status = 'won'
      game.value.winner = game.value.currentPlayer === 1 ? 2 : 1

      // 更新玩家数据
      const playerStore = usePlayerStore()
      if (playerStore.player) {
        playerStore.updateLocal({
          totalLosses: playerStore.player.totalLosses + 1
        })
      }
      return
    }

    isLoading.value = true
    try {
      const res = await gameApi.surrender(game.value.gameId)
      if (res.data.code === 0 && res.data.data) {
        game.value = res.data.data

        if (game.value.gameMode === 'ai') {
          const playerStore = usePlayerStore()
          if (playerStore.player) {
            playerStore.updateLocal({
              totalLosses: playerStore.player.totalLosses + 1
            })
          }
        }
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || '网络错误'
    } finally {
      isLoading.value = false
    }
  }

  // 悔棋
  async function undo() {
    if (!game.value || !isPlaying.value) return
    if (game.value.gameMode === 'ai' && game.value.moveHistory.length < 2) {
      return
    }

    if (useLocalMode.value) {
      if (game.value.moveHistory.length === 0) return

      const lastMove = game.value.moveHistory[game.value.moveHistory.length - 1]
      game.value.board[lastMove.y][lastMove.x] = 0
      game.value.moveHistory.pop()
      game.value.lastMove = game.value.moveHistory.length > 0
        ? game.value.moveHistory[game.value.moveHistory.length - 1]
        : undefined
      game.value.status = 'playing'
      game.value.winner = undefined
      game.value.currentPlayer = 1
      winLine.value = []
      return
    }

    isLoading.value = true
    try {
      const res = await gameApi.undo(game.value.gameId)
      if (res.data.code === 0 && res.data.data) {
        game.value = res.data.data
        winLine.value = []
      }
    } catch (e: any) {
      error.value = e.response?.data?.error || '网络错误'
    } finally {
      isLoading.value = false
    }
  }

  // 计算获胜连线
  function calculateWinLine(board: number[][], x: number, y: number): Position[] {
    const color = board[y][x]
    if (color === 0) return []

    const directions = [
      { dx: 1, dy: 0 },
      { dx: 0, dy: 1 },
      { dx: 1, dy: 1 },
      { dx: 1, dy: -1 }
    ]

    for (const dir of directions) {
      const line: Position[] = [{ x, y }]

      for (let i = 1; i < 5; i++) {
        const nx = x + dir.dx * i
        const ny = y + dir.dy * i
        if (nx >= 0 && nx < 15 && ny >= 0 && ny < 15 && board[ny][nx] === color) {
          line.push({ x: nx, y: ny })
        } else {
          break
        }
      }

      for (let i = 1; i < 5; i++) {
        const nx = x - dir.dx * i
        const ny = y - dir.dy * i
        if (nx >= 0 && nx < 15 && ny >= 0 && ny < 15 && board[ny][nx] === color) {
          line.unshift({ x: nx, y: ny })
        } else {
          break
        }
      }

      if (line.length >= 5) {
        return line
      }
    }

    return []
  }

  // 清除游戏状态
  function clearGame() {
    game.value = null
    error.value = null
    winLine.value = []
  }

  return {
    game,
    isLoading,
    error,
    winLine,
    isPlaying,
    currentPlayerName,
    isPlayerTurn,
    startGame,
    makeMove,
    restartGame,
    surrender,
    undo,
    clearGame
  }
})