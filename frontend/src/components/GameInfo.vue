<template>
  <div class="game-info">
    <div class="player-indicator">
      <div class="player black" :class="{ active: currentPlayer === 1 }">
        <span class="piece">●</span>
        <span class="name">黑方</span>
      </div>
      <div class="vs">VS</div>
      <div class="player white" :class="{ active: currentPlayer === 2 }">
        <span class="piece">●</span>
        <span class="name">{{ opponentName }}</span>
      </div>
    </div>

    <div class="status" v-if="status !== 'playing'">
      <span class="result" :class="resultClass">{{ resultText }}</span>
    </div>
    <div class="turn" v-else>
      {{ turnText }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPlayer: number
  status: string
  winner?: number
  gameMode: string
  difficulty?: string
}>()

const opponentName = computed(() => {
  if (props.gameMode === 'pvp') return '白方'
  const names: Record<string, string> = {
    easy: '小白',
    medium: '小灵',
    hard: '大师'
  }
  return names[props.difficulty || 'easy'] || '小白'
})

const turnText = computed(() => {
  return props.currentPlayer === 1 ? '轮到你落子' : '等待对方落子...'
})

const resultText = computed(() => {
  if (props.status === 'draw') return '平局！'
  if (props.winner === 1) return '黑方获胜！'
  return '白方获胜！'
})

const resultClass = computed(() => {
  if (props.status === 'draw') return 'draw'
  if (props.winner === 1) return 'win'
  return 'lose'
})
</script>

<style scoped>
.game-info {
  text-align: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.player-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 12px;
}

.player {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.player.active {
  background: linear-gradient(145deg, #f0e6d3, #e8d9c5);
  box-shadow: 0 2px 8px rgba(139, 69, 19, 0.3);
}

.player .piece {
  font-size: 24px;
  line-height: 1;
}

.player.black .piece {
  color: #1a1a1a;
}

.player.white .piece {
  color: #f0f0f0;
  text-shadow: 0 0 1px #666;
}

.player .name {
  font-size: 14px;
  margin-top: 4px;
  color: #5a4a3a;
}

.vs {
  font-size: 14px;
  font-weight: bold;
  color: #8b4513;
}

.status {
  margin-top: 12px;
}

.result {
  font-size: 20px;
  font-weight: bold;
  padding: 8px 16px;
  border-radius: 8px;
}

.result.win {
  color: #2ecc71;
  background: rgba(46, 204, 113, 0.1);
}

.result.lose {
  color: #e74c3c;
  background: rgba(231, 76, 60, 0.1);
}

.result.draw {
  color: #3498db;
  background: rgba(52, 152, 219, 0.1);
}

.turn {
  font-size: 16px;
  color: #5a4a3a;
}
</style>