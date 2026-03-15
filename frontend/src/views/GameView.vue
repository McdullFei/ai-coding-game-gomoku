<template>
  <div class="game-view">
    <div class="header">
      <div class="title">五子棋</div>
      <RankBadge v-if="playerStore.player" :player="playerStore.player" />
    </div>

    <div class="game-container">
      <div class="game-content">
        <ChessBoard
          v-if="gameStore.game"
          :board="gameStore.game.board"
          :lastMove="gameStore.game.lastMove"
          :winLine="gameStore.winLine"
          :currentPlayer="gameStore.game.currentPlayer"
          :isPlaying="gameStore.isPlaying"
          :gameMode="gameStore.game.gameMode"
          @move="handleMove"
        />
        <div v-else class="loading">加载中...</div>
      </div>

      <div class="sidebar">
        <GameInfo
          v-if="gameStore.game"
          :currentPlayer="gameStore.game.currentPlayer"
          :status="gameStore.game.status"
          :winner="gameStore.game.winner"
          :gameMode="gameStore.game.gameMode"
          :difficulty="gameStore.game.difficulty"
        />

        <GameControl @restart="handleRestart" />
      </div>
    </div>

    <!-- 游戏结束弹窗 -->
    <div class="modal" v-if="showResult" @click.self="closeResult">
      <div class="result-modal">
        <h2>{{ resultTitle }}</h2>
        <p class="result-text" :class="resultClass">{{ resultMessage }}</p>
        <div class="exp-gain" v-if="expGain > 0">
          <span class="exp-icon">⭐</span>
          <span>+{{ expGain }} 经验值</span>
        </div>
        <div class="result-buttons">
          <button @click="handleRestart">再来一局</button>
          <button class="back" @click="goHome">返回菜单</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/game'
import { usePlayerStore } from '../stores/player'
import ChessBoard from '../components/ChessBoard.vue'
import GameInfo from '../components/GameInfo.vue'
import GameControl from '../components/GameControl.vue'
import RankBadge from '../components/RankBadge.vue'

const router = useRouter()
const gameStore = useGameStore()
const playerStore = usePlayerStore()

const showResult = ref(false)
const expGain = ref(0)

// 监听游戏状态变化
watch(() => gameStore.game?.status, async (newStatus) => {
  if (newStatus === 'won' || newStatus === 'draw') {
    // 计算经验值
    if (gameStore.game?.gameMode === 'ai' && gameStore.game.winner === 1) {
      const diff = gameStore.game.difficulty || 'easy'
      expGain.value = diff === 'easy' ? 10 : diff === 'medium' ? 25 : 50
    } else if (gameStore.game?.status === 'draw') {
      expGain.value = 3
    } else if (gameStore.game?.gameMode === 'ai') {
      expGain.value = 5 // 失败参与奖
    }

    // 更新本地经验值
    if (gameStore.game?.gameMode === 'ai' && playerStore.player) {
      const currentExp = playerStore.player.experience
      const newExp = currentExp + expGain.value
      const newLevel = Math.floor(newExp / 100) + 1

      const segments = ['新手', '初级', '中级', '高级', '专家', '大师']
      let segment = '新手'
      if (newLevel >= 51) segment = '大师'
      else if (newLevel >= 41) segment = '专家'
      else if (newLevel >= 31) segment = '高级'
      else if (newLevel >= 21) segment = '中级'
      else if (newLevel >= 11) segment = '初级'

      playerStore.updateLocal({
        experience: newExp,
        rankLevel: newLevel,
        rankSegment: segment,
        totalWins: gameStore.game.winner === 1 ? playerStore.player.totalWins + 1 : playerStore.player.totalWins,
        totalLosses: gameStore.game.winner !== 1 && gameStore.game.status === 'won' ? playerStore.player.totalLosses + 1 : playerStore.player.totalLosses
      })
    }

    showResult.value = true
  }
})

const resultTitle = computed(() => {
  if (!gameStore.game) return ''
  if (gameStore.game.status === 'draw') return '平局'
  if (gameStore.game.winner === 1) return '你赢了！'
  return '你输了'
})

const resultMessage = computed(() => {
  if (!gameStore.game) return ''
  if (gameStore.game.status === 'draw') return '棋盘已满，难分胜负'
  if (gameStore.game.winner === 1) return '恭喜你获得胜利！'
  return gameStore.game.gameMode === 'ai' ? '再接再厉！' : '白方获胜'
})

const resultClass = computed(() => {
  if (!gameStore.game) return ''
  if (gameStore.game.status === 'draw') return 'draw'
  return gameStore.game.winner === 1 ? 'win' : 'lose'
})

function handleMove(x: number, y: number) {
  gameStore.makeMove(x, y)
}

function handleRestart() {
  showResult.value = false
  gameStore.restartGame()
  expGain.value = 0
}

function closeResult() {
  showResult.value = false
}

function goHome() {
  showResult.value = false
  gameStore.clearGame()
  router.push('/')
}

onMounted(() => {
  if (!gameStore.game) {
    router.push('/')
  }
})
</script>

<style scoped>
.game-view {
  min-height: 100vh;
  background: linear-gradient(145deg, #f5e6d3, #e8d9c5);
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 800px;
  margin: 0 auto 20px;
}

.title {
  font-size: 28px;
  color: #5a4a3a;
  font-weight: bold;
}

.game-container {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.game-content {
  flex: 1;
}

.sidebar {
  width: 200px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #8b7355;
}

/* 结果弹窗 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.result-modal {
  background: white;
  padding: 32px 48px;
  border-radius: 16px;
  text-align: center;
}

.result-modal h2 {
  font-size: 28px;
  margin-bottom: 12px;
  color: #5a4a3a;
}

.result-text {
  font-size: 18px;
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 8px;
}

.result-text.win {
  color: #27ae60;
  background: rgba(39, 174, 96, 0.1);
}

.result-text.lose {
  color: #e74c3c;
  background: rgba(231, 76, 60, 0.1);
}

.result-text.draw {
  color: #3498db;
  background: rgba(52, 152, 219, 0.1);
}

.exp-gain {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 20px;
  color: #f39c12;
  margin-bottom: 20px;
}

.exp-icon {
  font-size: 24px;
}

.result-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.result-buttons button {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  background: linear-gradient(145deg, #667eea, #764ba2);
  color: white;
}

.result-buttons button.back {
  background: #e0e0e0;
  color: #666;
}
</style>