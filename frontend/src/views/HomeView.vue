<template>
  <div class="home">
    <div class="container">
      <h1 class="title">五子棋</h1>

      <div class="rank-section" v-if="playerStore.player">
        <RankBadge :player="playerStore.player" />
        <div class="exp-bar">
          <div class="exp-fill" :style="{ width: playerStore.expProgress + '%' }"></div>
        </div>
        <div class="exp-text">{{ playerStore.player.experience }} 经验值</div>
      </div>

      <div class="menu">
        <button class="menu-btn" @click="startPVP">
          双人对战
        </button>

        <div class="difficulty-section">
          <button class="menu-btn ai" @click="showDifficulty = true">
            人机对战
          </button>
        </div>
      </div>

      <!-- 难度选择弹窗 -->
      <div class="modal" v-if="showDifficulty" @click.self="showDifficulty = false">
        <div class="modal-content">
          <h2>选择难度</h2>
          <div class="difficulty-list">
            <button
              v-for="diff in difficulties"
              :key="diff.value"
              class="difficulty-btn"
              @click="startAI(diff.value)"
            >
              <span class="diff-name">{{ diff.name }}</span>
              <span class="diff-desc">{{ diff.desc }}</span>
            </button>
          </div>
          <button class="close-btn" @click="showDifficulty = false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/game'
import { usePlayerStore } from '../stores/player'
import RankBadge from '../components/RankBadge.vue'
import type { Difficulty } from '../types'

const router = useRouter()
const gameStore = useGameStore()
const playerStore = usePlayerStore()

const showDifficulty = ref(false)

const difficulties = [
  { value: 'easy', name: '小白', desc: '新手级别，反应迟缓' },
  { value: 'medium', name: '小灵', desc: '具备一定思考能力' },
  { value: 'hard', name: '大师', desc: '顶尖高手，极难战胜' }
]

onMounted(() => {
  playerStore.initPlayer()
})

async function startPVP() {
  await gameStore.startGame('pvp')
  router.push('/game')
}

async function startAI(difficulty: Difficulty) {
  showDifficulty.value = false
  await gameStore.startGame('ai', difficulty)
  router.push('/game')
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(145deg, #f5e6d3, #e8d9c5);
}

.container {
  text-align: center;
  padding: 40px;
}

.title {
  font-size: 48px;
  color: #5a4a3a;
  margin-bottom: 40px;
  text-shadow: 2px 2px 4px rgba(139, 69, 19, 0.2);
}

.rank-section {
  margin-bottom: 40px;
}

.exp-bar {
  width: 200px;
  height: 8px;
  background: rgba(139, 69, 19, 0.2);
  border-radius: 4px;
  margin: 12px auto 8px;
  overflow: hidden;
}

.exp-fill {
  height: 100%;
  background: linear-gradient(90deg, #e67e22, #f39c12);
  transition: width 0.3s ease;
}

.exp-text {
  font-size: 12px;
  color: #8b7355;
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.menu-btn {
  width: 200px;
  padding: 16px 32px;
  font-size: 18px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  background: linear-gradient(145deg, #fff, #f0f0f0);
  color: #5a4a3a;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.2s ease;
}

.menu-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.menu-btn.ai {
  background: linear-gradient(145deg, #667eea, #764ba2);
  color: white;
}

/* 弹窗样式 */
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

.modal-content {
  background: white;
  padding: 32px;
  border-radius: 16px;
  text-align: center;
}

.modal-content h2 {
  margin-bottom: 24px;
  color: #5a4a3a;
}

.difficulty-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.difficulty-btn {
  padding: 16px 24px;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  background: linear-gradient(145deg, #f5e6d3, #e8d9c5);
  text-align: left;
  transition: all 0.2s ease;
}

.difficulty-btn:hover {
  transform: translateX(4px);
}

.diff-name {
  display: block;
  font-size: 18px;
  font-weight: bold;
  color: #5a4a3a;
}

.diff-desc {
  display: block;
  font-size: 12px;
  color: #8b7355;
  margin-top: 4px;
}

.close-btn {
  padding: 10px 24px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  background: #e0e0e0;
  color: #666;
}
</style>