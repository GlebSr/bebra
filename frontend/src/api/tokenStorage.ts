const ACCESS_TOKEN_KEY = 'access_token';

let memoryToken: string | null = null;

export const tokenStorage = {
  getToken(): string | null {
    if (memoryToken) return memoryToken;
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },

  setToken(token: string): void {
    memoryToken = token;
    localStorage.setItem(ACCESS_TOKEN_KEY, token);
  },

  clearToken(): void {
    memoryToken = null;
    localStorage.removeItem(ACCESS_TOKEN_KEY);
  },

  hasToken(): boolean {
    return !!this.getToken();
  }
};
