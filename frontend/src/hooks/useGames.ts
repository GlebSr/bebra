import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { gamesApi, type CreateGameRequest } from '@/api';

export function useGames(roomId: string) {
  return useQuery({
    queryKey: ['games', roomId],
    queryFn: () => gamesApi.getAll(roomId),
    enabled: !!roomId,
  });
}

export function useAddGame(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateGameRequest) => gamesApi.add(roomId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['games', roomId] });
    },
  });
}

export function useDeleteGame(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (gameId: string) => gamesApi.delete(roomId, gameId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['games', roomId] });
    },
  });
}
