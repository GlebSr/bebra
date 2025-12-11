import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import { authApi, userApi, tokenStorage, type User } from '@/api';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  signUp: (name: string, password: string) => Promise<{ error?: string }>;
  signIn: (name: string, password: string) => Promise<{ error?: string }>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Проверка аутентификации при инициализации
  useEffect(() => {
    const initAuth = async () => {
      if (!tokenStorage.hasToken()) {
        setIsLoading(false);
        return;
      }

      try {
        const userData = await userApi.getMe();
        setUser(userData);
      } catch {
        // Токен недействителен, очистить его
        tokenStorage.clearToken();
      } finally {
        setIsLoading(false);
      }
    };

    initAuth();
  }, []);

  const signUp = useCallback(async (name: string, password: string) => {
    try {
      const response = await authApi.signUp({ name, password });
      const userData = await userApi.getMe();
      setUser(userData);
      return {};
    } catch (error: any) {
      return { error: error.data?.error || error.message || 'Ошибка регистрации' };
    }
  }, []);

  const signIn = useCallback(async (name: string, password: string) => {
    try {
      await authApi.signIn({ name, password });
      const userData = await userApi.getMe();
      setUser(userData);
      return {};
    } catch (error: any) {
      return { error: error.data?.error || error.message || 'Ошибка входа' };
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      await authApi.logout();
    } catch {
      // Игнорировать ошибки выхода
    } finally {
      setUser(null);
    }
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        signUp,
        signIn,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
