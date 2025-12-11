import { useState, useMemo } from "react";
import RoomsList from "@/components/RoomsList";
import GamesList from "@/components/GamesList";
import RoomSettings from "@/components/RoomSettings";
import { useToast } from "@/hooks/use-toast";
import { useRooms, useUpdateRoom, useDeleteRoom } from "@/hooks/useRooms";
import { useGames, useAddGame, useDeleteGame } from "@/hooks/useGames";
import { useVotes, useAddVote, useDeleteVote } from "@/hooks/useVotes";
import { useInviteParticipant, useLeaveRoom } from "@/hooks/useParticipants";
import { useGetRandom, useRandomHistory } from "@/hooks/useRandom";
import { useRoomWebSocket } from "@/hooks/useRoomWebSocket";
import { useAuth } from "@/contexts/AuthContext";
import { Loader2 } from "lucide-react";

const Index = () => {
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null);
  const [winner, setWinner] = useState<string | null>(null);
  const { toast } = useToast();
  const { user } = useAuth();

  // Получение комнат
  const { data: rooms = [], isLoading: roomsLoading } = useRooms();

  // Автоматический выбор первой комнаты, если ни одна не выбрана
  const activeRoomId = selectedRoomId || rooms[0]?.id || null;
  const hasRooms = rooms.length > 0;

  // WebSocket для обновлений в реальном времени
  useRoomWebSocket(activeRoomId);

  // Получение игр и голосов для выбранной комнаты (только если комната существует)
  const { data: games = [], isLoading: gamesLoading } = useGames(activeRoomId || '');
  const { data: votes = [] } = useVotes(activeRoomId || '');

  // Mutations
  const addGameMutation = useAddGame(activeRoomId || '');
  const deleteGameMutation = useDeleteGame(activeRoomId || '');
  const addVoteMutation = useAddVote(activeRoomId || '');
  const deleteVoteMutation = useDeleteVote(activeRoomId || '');
  const updateRoomMutation = useUpdateRoom(activeRoomId || '');
  const deleteRoomMutation = useDeleteRoom();
  const inviteParticipantMutation = useInviteParticipant(activeRoomId || '');
  const leaveRoomMutation = useLeaveRoom(activeRoomId || '');
  const getRandomMutation = useGetRandom(activeRoomId || '');
  const { data: randomHistory = [] } = useRandomHistory(activeRoomId || '');

  // Подсчёт голосов для каждой игры
  const gamesWithVotes = useMemo(() => {
    const voteCount: Record<string, number> = {};
    votes.forEach(vote => {
      voteCount[vote.game_id] = (voteCount[vote.game_id] || 0) + 1;
    });
    return games.map(game => ({
      id: game.id,
      name: game.title,
      votes: voteCount[game.id] || 0,
    }));
  }, [games, votes]);

  // Получение ID игр, за которые проголосовал текущий пользователь
  const userVotedGameIds = useMemo(() => {
    if (!user) return [];
    return votes
      .filter(vote => vote.user_id === user.id)
      .map(vote => vote.game_id);
  }, [votes, user]);

  // Карта голосов пользователя по game_id
  const userVotesByGameId = useMemo(() => {
    if (!user) return {};
    const map: Record<string, string> = {};
    votes
      .filter(vote => vote.user_id === user.id)
      .forEach(vote => {
        map[vote.game_id] = vote.id;
      });
    return map;
  }, [votes, user]);

  // Преобразование истории из game_id в объекты с названиями
  const historyWithNames = useMemo(() => {
    return randomHistory.map(gameId => {
      const game = games.find(g => g.id === gameId);
      return {
        gameId,
        gameName: game?.title || 'Неизвестная игра',
      };
    });
  }, [randomHistory, games]);

  const currentRoom = rooms.find((room) => room.id === activeRoomId);
  const isOwner = !!(user && currentRoom && currentRoom.owner_id === user.id);

  const handleVote = async (gameId: string) => {
    try {
      await addVoteMutation.mutateAsync(gameId);
      toast({
        title: "Голос учтён",
        description: "Ваш голос успешно добавлен",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось проголосовать",
      });
    }
  };

  const handleRemoveVote = async (gameId: string) => {
    const voteId = userVotesByGameId[gameId];
    if (!voteId) return;

    try {
      await deleteVoteMutation.mutateAsync(voteId);
      toast({
        title: "Голос отменён",
        description: "Ваш голос успешно удалён",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось отменить голос",
      });
    }
  };

  const handleAddGame = async (gameName: string) => {
    try {
      await addGameMutation.mutateAsync({ title: gameName });
      toast({
        title: "Игра добавлена",
        description: `"${gameName}" добавлена в список`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось добавить игру",
      });
    }
  };

  const handleDeleteGame = async (gameId: string) => {
    try {
      await deleteGameMutation.mutateAsync(gameId);
      toast({
        title: "Игра удалена",
        description: "Игра успешно удалена из списка",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось удалить игру",
      });
    }
  };

  const handleRandomSelect = async () => {
    try {
      const gameId = await getRandomMutation.mutateAsync();
      const game = games.find(g => g.id === gameId);
      if (game) {
        setWinner(game.title);
        toast({
          title: "Победитель выбран!",
          description: `Случайно выбрана игра: ${game.title}`,
        });
      }
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось выбрать случайную игру",
      });
    }
  };

  const handleAddUser = async (username: string) => {
    try {
      await inviteParticipantMutation.mutateAsync(username);
      toast({
        title: "Пользователь добавлен",
        description: `Пользователь "${username}" добавлен в комнату`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось добавить пользователя",
      });
    }
  };

  const handleRenameRoom = async (newName: string) => {
    try {
      await updateRoomMutation.mutateAsync({ name: newName });
      toast({
        title: "Комната переименована",
        description: `Комната переименована в "${newName}"`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось переименовать комнату",
      });
    }
  };

  const handleDeleteRoom = async () => {
    if (!activeRoomId) return;
    const roomName = currentRoom?.name;

    try {
      await deleteRoomMutation.mutateAsync(activeRoomId);

      // Выбор другой комнаты после удаления
      const remainingRooms = rooms.filter((room) => room.id !== activeRoomId);
      if (remainingRooms.length > 0) {
        setSelectedRoomId(remainingRooms[0].id);
      } else {
        setSelectedRoomId(null);
      }

      toast({
        title: "Комната удалена",
        description: `Комната "${roomName}" успешно удалена`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось удалить комнату",
      });
    }
  };

  const handleLeaveRoom = async () => {
    const roomName = currentRoom?.name;
    try {
      await leaveRoomMutation.mutateAsync();
      
      // Выбор другой комнаты после выхода
      const remainingRooms = rooms.filter((room) => room.id !== activeRoomId);
      if (remainingRooms.length > 0) {
        setSelectedRoomId(remainingRooms[0].id);
      } else {
        setSelectedRoomId(null);
      }

      toast({
        title: "Вы покинули комнату",
        description: `Вы вышли из комнаты "${roomName}"`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось покинуть комнату",
      });
    }
  };

  if (roomsLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-background">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <RoomsList
        rooms={rooms}
        selectedRoomId={activeRoomId || ''}
        onRoomSelect={setSelectedRoomId}
      />
      {hasRooms && activeRoomId ? (
        <>
          <GamesList
            games={gamesWithVotes}
            userVotedGameIds={userVotedGameIds}
            onVote={handleVote}
            onRemoveVote={handleRemoveVote}
            onAddGame={handleAddGame}
            onDeleteGame={handleDeleteGame}
            onRandomSelect={handleRandomSelect}
            winner={winner}
            onCloseWinner={() => setWinner(null)}
            isLoading={gamesLoading}
            history={historyWithNames}
          />
          {currentRoom && (
            <RoomSettings
              roomName={currentRoom.name}
              isOwner={isOwner}
              onAddUser={handleAddUser}
              onRenameRoom={handleRenameRoom}
              onDeleteRoom={handleDeleteRoom}
              onLeaveRoom={handleLeaveRoom}
            />
          )}
        </>
      ) : (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <h2 className="text-2xl font-semibold text-foreground mb-2">
              Нет комнат
            </h2>
            <p className="text-muted-foreground">
              Создайте первую комнату, нажав "+" в боковой панели
            </p>
          </div>
        </div>
      )}
    </div>
  );
};

export default Index;
