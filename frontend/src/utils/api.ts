import axios from 'axios'
import type { GameState, Player, ApiResponse } from '../types'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 玩家相关API
export const playerApi = {
  register: (username: string) =>
    api.post<ApiResponse<Player>>('/player/register', { username }),

  getPlayer: (userId: string) =>
    api.get<ApiResponse<Player>>(`/player/get/${userId}`)
}

// 游戏相关API
export const gameApi = {
  start: (gameMode: string, difficulty?: string, userId?: string) =>
    api.post<ApiResponse<GameState>>('/game/start', { gameMode, difficulty, userId }),

  move: (gameId: string, x: number, y: number) =>
    api.post<ApiResponse<GameState>>('/game/move', { gameId, x, y }),

  end: (gameId: string, userId?: string) =>
    api.post<ApiResponse<any>>('/game/end', { gameId, userId }),

  getGame: (gameId: string) =>
    api.get<ApiResponse<GameState>>(`/game/${gameId}`),

  surrender: (gameId: string) =>
    api.post<ApiResponse<GameState>>('/game/surrender', { gameId }),

  restart: (gameId: string, difficulty?: string) =>
    api.post<ApiResponse<GameState>>('/game/restart', { gameId, difficulty }),

  undo: (gameId: string) =>
    api.post<ApiResponse<GameState>>('/game/undo', { gameId })
}

export default api