import { apiClient } from '../client';
import { tokenStorage } from '../tokenStorage';
import type { AuthResponse, SignUpRequest, SignInRequest } from '../types';

export const authApi = {
  async signUp(data: SignUpRequest): Promise<AuthResponse> {
    const response = await apiClient.authPost<AuthResponse>('/auth/signup', data);
    tokenStorage.setToken(response.access_token);
    return response;
  },

  async signIn(data: SignInRequest): Promise<AuthResponse> {
    const response = await apiClient.authPost<AuthResponse>('/auth/signin', data);
    tokenStorage.setToken(response.access_token);
    return response;
  },

  async refresh(): Promise<AuthResponse> {
    const response = await apiClient.authPost<AuthResponse>('/auth/refresh');
    tokenStorage.setToken(response.access_token);
    return response;
  },

  async logout(): Promise<void> {
    await apiClient.authPost<void>('/auth/logout');
    tokenStorage.clearToken();
  },
};
