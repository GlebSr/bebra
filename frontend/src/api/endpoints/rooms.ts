import { apiClient } from '../client';
import type { Room, RoomsResponse, CreateRoomRequest, UpdateRoomRequest } from '../types';

export const roomsApi = {
  create(data: CreateRoomRequest): Promise<Room> {
    return apiClient.post<Room>('/rooms', data);
  },

  getAll(): Promise<RoomsResponse> {
    return apiClient.get<RoomsResponse>('/rooms');
  },

  getById(roomId: string): Promise<Room> {
    return apiClient.get<Room>(`/rooms/${roomId}`);
  },

  update(roomId: string, data: UpdateRoomRequest): Promise<Room> {
    return apiClient.put<Room>(`/rooms/${roomId}`, data);
  },

  delete(roomId: string): Promise<void> {
    return apiClient.delete<void>(`/rooms/${roomId}`);
  },
};
