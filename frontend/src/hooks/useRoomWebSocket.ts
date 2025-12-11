import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { createRoomWebSocket } from '@/api';

export function useRoomWebSocket(roomId: string | null) {
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!roomId) return;

    const ws = createRoomWebSocket(roomId);
    ws.connect();

    ws.on('game.added', () => {
      queryClient.invalidateQueries({ queryKey: ['games', roomId] });
    });

    ws.on('game.deleted', () => {
      queryClient.invalidateQueries({ queryKey: ['games', roomId] });
    });

    ws.on('vote.added', () => {
      queryClient.invalidateQueries({ queryKey: ['votes', roomId] });
    });

    ws.on('vote.deleted', () => {
      queryClient.invalidateQueries({ queryKey: ['votes', roomId] });
    });

    ws.on('room.updated', () => {
      queryClient.invalidateQueries({ queryKey: ['rooms'] });
    });

    ws.on('participant.added', () => {
      queryClient.invalidateQueries({ queryKey: ['participants', roomId] });
    });

    ws.on('participant.left', () => {
      queryClient.invalidateQueries({ queryKey: ['participants', roomId] });
    });

    ws.on('results.updated', () => {
      queryClient.invalidateQueries({ queryKey: ['random', 'history', roomId] });
    });

    return () => ws.disconnect();
  }, [roomId, queryClient]);
}
