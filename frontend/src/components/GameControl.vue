<template>
  <div class="game-control">
    <button class="btn" @click="handleUndo" :disabled="!canUndo">
      悔棋
    </button>
    <button class="btn" @click="handleRestart">
      重新开始
    </button>
    <button class="btn surrender" @click="handleSurrender" :disabled="!canSurrender">
      认输
    </button>
    <button class="btn back" @click="handleBack">
      返回菜单
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/game'

const router = useRouter()
const gameStore = useGameStore()

const canUndo = computed(() => {
  if (!gameStore.game || gameStore.game.status !== 'playing') return false
  if (gameStore.game.moveHistory.length === 0) return false
  // 人机模式需要至少2步才能悔棋
  if (gameStore.game.gameMode === 'ai' && gameStore.game.moveHistory.length < 2) return false
  return true
})

const canSurrender = computed(() => {
  return gameStore.game?.status === 'playing'
})

function handleUndo() {
  gameStore.undo()
}

function handleRestart() {
  gameStore.restartGame()
}

function handleSurrender() {
  if (confirm('确定要认输吗？')) {
    gameStore.surrender()
  }
}

function handleBack() {
  gameStore.clearGame()
  router.push('/')
}
</script>

<style scoped>
.game-control {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}

.btn {
  padding: 10px 20px;
  font-size: 14px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  background: linear-gradient(145deg, #f5e6d3, #e8d9c5);
  color: #5a4a3a;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
}

.btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.btn:active:not(:disabled) {
  transform: translateY(0);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn.surrender {
  background: linear-gradient(145deg, #ffe0e0, #ffcccc);
  color: #c0392b;
}

.btn.back {
  background: linear-gradient(145deg, #e0e8ff, #ccd4ff);
  color: #2c3e50;
}
</style>