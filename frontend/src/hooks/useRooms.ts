import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { roomsApi, type Room, type CreateRoomRequest, type UpdateRoomRequest } from '@/api';

export function useRooms() {
  return useQuery({
    queryKey: ['rooms'],
    queryFn: async () => {
      const response = await roomsApi.getAll();
      return response.rooms;
    },
  });
}

export function useRoom(roomId: string) {
  return useQuery({
    queryKey: ['room', roomId],
    queryFn: () => roomsApi.getById(roomId),
    enabled: !!roomId,
  });
}

export function useCreateRoom() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateRoomRequest) => roomsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });
}

export function useUpdateRoom(roomId: string) {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: UpdateRoomRequest) => roomsApi.update(roomId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
      queryClient.invalidateQueries({ queryKey: ['room', roomId] });
    },
  });
}

export function useDeleteRoom() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (roomId: string) => roomsApi.delete(roomId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    },
  });
}
