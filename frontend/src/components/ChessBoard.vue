<template>
  <div class="chess-board" :class="{ 'disabled': !canClick }">
    <div class="board-grid">
      <div
        v-for="(row, y) in board"
        :key="y"
        class="board-row"
      >
        <div
          v-for="(cell, x) in row"
          :key="x"
          class="board-cell"
          :class="{
            'is-last-move': isLastMove(x, y),
            'is-win-line': isWinLine(x, y)
          }"
          @click="handleClick(x, y)"
        >
          <!-- 棋子 -->
          <div
            v-if="cell !== 0"
            class="chess-piece"
            :class="cell === 1 ? 'black' : 'white'"
          >
            <!-- 最后落子标记 -->
            <div v-if="isLastMove(x, y)" class="last-move-dot"></div>
          </div>
          <!-- 交叉点标记（星位） -->
          <div v-else-if="isStarPoint(x, y)" class="star-point"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Position } from '../types'

const props = defineProps<{
  board: number[][]
  lastMove?: Position
  winLine: Position[]
  currentPlayer: number
  isPlaying: boolean
  gameMode: string
}>()

const emit = defineEmits<{
  (e: 'move', x: number, y: number): void
}>()

// 星位点（天元和四角星）
const starPoints = [
  [3, 3], [3, 7], [3, 11],
  [7, 3], [7, 7], [7, 11],
  [11, 3], [11, 7], [11, 11]
]

function isStarPoint(x: number, y: number): boolean {
  return starPoints.some(p => p[0] === x && p[1] === y)
}

function isLastMove(x: number, y: number): boolean {
  return props.lastMove?.x === x && props.lastMove?.y === y
}

function isWinLine(x: number, y: number): boolean {
  return props.winLine.some(p => p.x === x && p.y === y)
}

// 判断是否可以点击
const canClick = computed(() => {
  if (!props.isPlaying) return false
  // 人机模式下，只有黑方（玩家）可以点击
  if (props.gameMode === 'ai' && props.currentPlayer !== 1) return false
  return true
})

function handleClick(x: number, y: number) {
  if (!canClick.value) return
  if (props.board[y][x] !== 0) return
  emit('move', x, y)
}
</script>

<style scoped>
.chess-board {
  background: linear-gradient(145deg, #deb887, #d2a86e);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.chess-board.disabled {
  pointer-events: none;
  opacity: 0.95;
}

.board-grid {
  display: flex;
  flex-direction: column;
}

.board-row {
  display: flex;
}

.board-cell {
  width: 36px;
  height: 36px;
  position: relative;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.board-cell::before {
  content: '';
  position: absolute;
  background: #8b4513;
}

.board-cell::after {
  content: '';
  position: absolute;
  background: #8b4513;
}

/* 横线 */
.board-cell::before {
  width: 100%;
  height: 1px;
  top: 50%;
  left: 0;
}

/* 竖线 */
.board-cell::after {
  width: 1px;
  height: 100%;
  top: 0;
  left: 50%;
}

/* 边缘处理 */
.board-row:first-child .board-cell::before {
  width: 50%;
  left: 50%;
}

.board-row:last-child .board-cell::before {
  width: 50%;
  left: 0;
}

.board-cell:first-child::after {
  height: 50%;
  top: 50%;
}

.board-cell:last-child::after {
  height: 50%;
  top: 0;
}

/* 四个角不显示线 */
.board-row:first-child .board-cell:first-child::before,
.board-row:first-child .board-cell:first-child::after,
.board-row:first-child .board-cell:last-child::before,
.board-row:first-child .board-cell:last-child::after,
.board-row:last-child .board-cell:first-child::before,
.board-row:last-child .board-cell:first-child::after,
.board-row:last-child .board-cell:last-child::before,
.board-row:last-child .board-cell:last-child::after {
  display: none;
}

.board-cell:hover:not(.chess-piece) {
  background-color: rgba(139, 69, 19, 0.1);
}

.chess-piece {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  position: relative;
  z-index: 1;
  box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.chess-piece.black {
  background: radial-gradient(circle at 30% 30%, #4a4a4a, #1a1a1a);
}

.chess-piece.white {
  background: radial-gradient(circle at 30% 30%, #ffffff, #d0d0d0);
}

.last-move-dot {
  position: absolute;
  width: 8px;
  height: 8px;
  background: #e74c3c;
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.star-point {
  width: 8px;
  height: 8px;
  background: #8b4513;
  border-radius: 50%;
  position: relative;
  z-index: 0;
}

.is-win-line {
  background-color: rgba(255, 215, 0, 0.4);
}
</style>