import { apiClient } from '../client';
import type { Game, GamesResponse, CreateGameRequest } from '../types';

export const gamesApi = {
  add(roomId: string, data: CreateGameRequest): Promise<Game> {
    return apiClient.post<Game>(`/rooms/${roomId}/games`, data);
  },

  getAll(roomId: string): Promise<GamesResponse> {
    return apiClient.get<GamesResponse>(`/rooms/${roomId}/games`);
  },

  delete(roomId: string, gameId: string): Promise<void> {
    return apiClient.delete<void>(`/rooms/${roomId}/games/${gameId}`);
  },
};
