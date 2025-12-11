import { tokenStorage } from './tokenStorage';
import type { ApiError } from './types';

const BASE_URL = '/api/v1';
const AUTH_BASE_URL = '/api';

export class ApiClientError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: ApiError
  ) {
    super(message);
    this.name = 'ApiClientError';
  }
}

async function refreshAccessToken(): Promise<boolean> {
  try {
    const response = await fetch(`${AUTH_BASE_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    });

    if (!response.ok) return false;

    const data = await response.json();
    if (data.access_token) {
      tokenStorage.setToken(data.access_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
  retry = true,
  useAuthBase = false
): Promise<T> {
  const baseUrl = useAuthBase ? AUTH_BASE_URL : BASE_URL;
  const url = `${baseUrl}${endpoint}`;
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const token = tokenStorage.getToken();
  if (token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(url, {
    ...options,
    headers,
    credentials: 'include',
  });

  // Попытка обновить токен при 401
  if (response.status === 401 && retry) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      return request<T>(endpoint, options, false);
    }
    tokenStorage.clearToken();
  }

  if (!response.ok) {
    let errorData: ApiError | undefined;
    try {
      errorData = await response.json();
    } catch {
      // Игнорим ошибки парсинга JSON
    }
    throw new ApiClientError(
      errorData?.error || `HTTP ${response.status}`,
      response.status,
      errorData
    );
  }

  // Для 204 No Content возвращаем undefined
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

export const apiClient = {
  get<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return request<T>(endpoint, { ...options, method: 'GET' });
  },

  post<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  },

  put<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  },

  delete<T>(endpoint: string, options?: RequestInit): Promise<T> {
    return request<T>(endpoint, { ...options, method: 'DELETE' });
  },

  // Аутентифицированный POST запрос к AUTH_BASE_URL
  authPost<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }, true, true);
  },
};
