import { apiClient } from '../client';
import type { RandomResult, RandomHistoryResponse } from '../types';

export const randomApi = {
  generate(roomId: string): Promise<RandomResult> {
    return apiClient.get<RandomResult>(`/rooms/${roomId}/random`);
  },

  getLast(roomId: string): Promise<RandomResult> {
    return apiClient.get<RandomResult>(`/rooms/${roomId}/random/last`);
  },

  getHistory(roomId: string): Promise<RandomHistoryResponse> {
    return apiClient.get<RandomHistoryResponse>(`/rooms/${roomId}/random/history`);
  },
};
