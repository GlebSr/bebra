import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { participantsApi } from '@/api';

export function useParticipants(roomId: string) {
  return useQuery({
    queryKey: ['participants', roomId],
    queryFn: () => participantsApi.getAll(roomId),
    enabled: !!roomId,
  });
}

export function useInviteParticipant(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (name: string) => participantsApi.invite(roomId, { name }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['participants', roomId] });
    },
  });
}

export function useLeaveRoom(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: () => participantsApi.leave(roomId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
      queryClient.invalidateQueries({ queryKey: ['participants', roomId] });
    },
  });
}
