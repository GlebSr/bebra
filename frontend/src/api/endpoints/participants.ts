import { apiClient } from '../client';
import type { ParticipantsResponse, InviteParticipantRequest } from '../types';

export const participantsApi = {
  invite(roomId: string, data: InviteParticipantRequest): Promise<void> {
    return apiClient.post<void>(`/rooms/${roomId}/participants`, data);
  },

  getAll(roomId: string): Promise<ParticipantsResponse> {
    return apiClient.get<ParticipantsResponse>(`/rooms/${roomId}/participants`);
  },

  leave(roomId: string): Promise<void> {
    return apiClient.delete<void>(`/rooms/${roomId}/participants`);
  },
};
