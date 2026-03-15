import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Player, RankSegment } from '../types'
import { playerApi } from '../utils/api'
import { RankSegments } from '../types'

const STORAGE_KEY = 'gomoku_player'

export const usePlayerStore = defineStore('player', () => {
  const player = ref<Player | null>(null)
  const isLoading = ref(false)

  // 计算当前段位进度
  const currentSegment = computed((): RankSegment | null => {
    if (!player.value) return null
    for (let i = RankSegments.length - 1; i >= 0; i--) {
      if (player.value.rankLevel >= RankSegments[i].minLevel) {
        return RankSegments[i]
      }
    }
    return RankSegments[0]
  })

  // 计算升级所需经验
  const expToNextLevel = computed(() => {
    if (!player.value || !currentSegment.value) return 0
    const currentIdx = RankSegments.findIndex(s => s.name === currentSegment.value?.name)
    if (currentIdx >= RankSegments.length - 1) {
      return 10000 // 大师后每级10000
    }
    return RankSegments[currentIdx + 1].expThreshold - player.value.experience
  })

  // 计算当前经验进度百分比
  const expProgress = computed(() => {
    if (!player.value || !currentSegment.value) return 0
    const currentIdx = RankSegments.findIndex(s => s.name === currentSegment.value?.name)
    if (currentIdx >= RankSegments.length - 1) {
      return Math.min(100, (player.value.experience % 10000) / 100)
    }
    const prevExp = currentSegment.value.expThreshold
    const nextExp = RankSegments[currentIdx + 1].expThreshold
    const progress = ((player.value.experience - prevExp) / (nextExp - prevExp)) * 100
    return Math.min(100, Math.max(0, progress))
  })

  // 从本地存储加载玩家
  function loadFromStorage() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        player.value = JSON.parse(stored)
      } catch (e) {
        console.error('加载玩家数据失败', e)
      }
    }
  }

  // 保存到本地存储
  function saveToStorage() {
    if (player.value) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(player.value))
    }
  }

  // 注册新玩家
  async function register(username: string) {
    isLoading.value = true
    try {
      const res = await playerApi.register(username)
      if (res.data.code === 0 && res.data.data) {
        player.value = res.data.data
        saveToStorage()
        return player.value
      }
      throw new Error(res.data.error || '注册失败')
    } finally {
      isLoading.value = false
    }
  }

  // 从服务器同步玩家数据
  async function syncFromServer(userId: string) {
    isLoading.value = true
    try {
      const res = await playerApi.getPlayer(userId)
      if (res.data.code === 0 && res.data.data) {
        player.value = res.data.data
        saveToStorage()
        return player.value
      }
    } catch (e) {
      console.error('同步玩家数据失败', e)
    } finally {
      isLoading.value = false
    }
  }

  // 更新玩家数据（本地）
  function updateLocal(data: Partial<Player>) {
    if (player.value) {
      player.value = { ...player.value, ...data }
      saveToStorage()
    }
  }

  // 初始化玩家
  async function initPlayer() {
    loadFromStorage()
    if (!player.value) {
      // 首次访问，创建本地玩家
      const guestId = 'guest_' + Date.now()
      player.value = {
        userId: guestId,
        username: '游客',
        rankLevel: 1,
        rankSegment: '新手',
        experience: 0,
        totalWins: 0,
        totalLosses: 0
      }
      saveToStorage()
    }
    return player.value
  }

  return {
    player,
    isLoading,
    currentSegment,
    expToNextLevel,
    expProgress,
    initPlayer,
    register,
    syncFromServer,
    updateLocal,
    loadFromStorage
  }
})