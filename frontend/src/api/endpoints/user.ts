import { apiClient } from '../client';
import type { User, UpdateUserRequest } from '../types';

export const userApi = {
  getMe(): Promise<User> {
    return apiClient.get<User>('/user');
  },

  getById(userId: string): Promise<User> {
    return apiClient.get<User>(`/user/${userId}`);
  },

  getByName(name: string): Promise<User> {
    return apiClient.get<User>(`/user/name/${encodeURIComponent(name)}`);
  },

  update(data: UpdateUserRequest): Promise<User> {
    return apiClient.put<User>('/user', data);
  },
};
