<template>
  <div class="rank-badge" v-if="player">
    <div class="segment">{{ player.rankSegment }}</div>
    <div class="level">Lv.{{ player.rankLevel }}</div>
  </div>
  <div class="rank-badge" v-else>
    <div class="segment">新手</div>
    <div class="level">Lv.1</div>
  </div>
</template>

<script setup lang="ts">
import type { Player } from '../types'

defineProps<{
  player?: Player | null
}>()
</script>

<style scoped>
.rank-badge {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 16px;
  background: linear-gradient(145deg, #f0e6d3, #e8d9c5);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(139, 69, 19, 0.2);
}

.segment {
  font-size: 16px;
  font-weight: bold;
  color: #8b4513;
}

.level {
  font-size: 12px;
  color: #a08060;
}
</style>