import { apiClient } from '../client';
import type { Vote, VotesResponse, AddVoteRequest } from '../types';

export const votesApi = {
  add(roomId: string, data: AddVoteRequest): Promise<Vote> {
    return apiClient.post<Vote>(`/rooms/${roomId}/votes`, data);
  },

  getAll(roomId: string): Promise<VotesResponse> {
    return apiClient.get<VotesResponse>(`/rooms/${roomId}/votes`);
  },

  delete(roomId: string, voteId: string): Promise<void> {
    return apiClient.delete<void>(`/rooms/${roomId}/votes/${voteId}`);
  },
};
