import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { randomApi } from '@/api';

export function useLastRandom(roomId: string) {
  return useQuery({
    queryKey: ['random', 'last', roomId],
    queryFn: () => randomApi.getLast(roomId),
    enabled: !!roomId,
  });
}

export function useRandomHistory(roomId: string) {
  return useQuery({
    queryKey: ['random', 'history', roomId],
    queryFn: async () => {
      const response = await randomApi.getHistory(roomId);
      // API может возвращать массив напрямую или объект
      return Array.isArray(response) ? response : [];
    },
    enabled: !!roomId,
  });
}

export function useGetRandom(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: () => randomApi.generate(roomId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['random', roomId] });
      queryClient.invalidateQueries({ queryKey: ['random', 'history', roomId] });
    },
  });
}
