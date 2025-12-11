import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { votesApi } from '@/api';
import type { Vote } from '@/api/types';

export function useVotes(roomId: string) {
  return useQuery({
    queryKey: ['votes', roomId],
    queryFn: async () => {
      const response = await votesApi.getAll(roomId);
      // API может возвращать массив напрямую или { votes: [...] }
      return Array.isArray(response) ? response : (response as { votes: Vote[] }).votes || [];
    },
    enabled: !!roomId,
  });
}

export function useAddVote(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (gameId: string) => votesApi.add(roomId, { game_id: gameId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['votes', roomId] });
    },
  });
}

export function useDeleteVote(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (voteId: string) => votesApi.delete(roomId, voteId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['votes', roomId] });
    },
  });
}
