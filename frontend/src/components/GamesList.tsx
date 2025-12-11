import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import AddGameDialog from "./AddGameDialog";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Trophy, Loader2, History, Trash2, Check } from "lucide-react";

interface Game {
  id: string;
  name: string;
  votes: number;
}

interface HistoryItem {
  gameId: string;
  gameName: string;
}

interface GamesListProps {
  games: Game[];
  userVotedGameIds: string[];
  onVote: (gameId: string) => void;
  onRemoveVote: (gameId: string) => void;
  onAddGame: (gameName: string) => void;
  onDeleteGame: (gameId: string) => void;
  onRandomSelect: () => void;
  winner: string | null;
  onCloseWinner: () => void;
  isLoading?: boolean;
  history?: HistoryItem[];
}

const GamesList = ({
  games,
  userVotedGameIds,
  onVote,
  onRemoveVote,
  onAddGame,
  onDeleteGame,
  onRandomSelect,
  winner,
  onCloseWinner,
  isLoading,
  history = [],
}: GamesListProps) => {
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [isHistoryOpen, setIsHistoryOpen] = useState(false);
  const [gameToDelete, setGameToDelete] = useState<Game | null>(null);

  const handleGameClick = (game: Game) => {
    if (userVotedGameIds.includes(game.id)) {
      onRemoveVote(game.id);
    } else {
      onVote(game.id);
    }
  };

  const reversedHistory = [...history].reverse();

  return (
    <main className="flex-1 flex flex-col h-screen">
      <div className="px-8 py-5 border-b border-border">
        <h1 className="text-xl font-semibold text-foreground">Голосование</h1>
      </div>

      <ScrollArea className="flex-1 px-8 py-6">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : games.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Нет игр в этой комнате</p>
            <p className="text-sm text-muted-foreground mt-1">Добавьте первую игру</p>
          </div>
        ) : (
          <div className="space-y-3 max-w-3xl p-1">
            {games.map((game) => {
              const hasVoted = userVotedGameIds.includes(game.id);
              return (
                <div
                  key={game.id}
                  className="flex items-center gap-2"
                >
                  <button
                    onClick={() => handleGameClick(game)}
                    className={`flex-1 rounded-xl px-6 py-4 transition-colors flex items-center justify-between ${
                      hasVoted
                        ? "bg-primary/20 hover:bg-primary/30 ring-2 ring-primary"
                        : "bg-game-card hover:bg-game-card-hover"
                    }`}
                  >
                    <div className="flex items-center gap-3">
                      {hasVoted && <Check className="w-5 h-5 text-primary" />}
                      <span className="text-foreground font-medium text-lg">{game.name}</span>
                    </div>
                    <span className="text-muted-foreground">Голосов: {game.votes}</span>
                  </button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-12 w-12 text-muted-foreground hover:text-destructive"
                    onClick={() => setGameToDelete(game)}
                  >
                    <Trash2 className="w-5 h-5" />
                  </Button>
                </div>
              );
            })}
          </div>
        )}
      </ScrollArea>

      <div className="px-8 py-6 border-t border-border flex gap-4">
        <Button
          onClick={() => setIsAddDialogOpen(true)}
          size="lg"
          className="flex-1 h-14 text-base rounded-full bg-secondary hover:bg-secondary/80"
        >
          Добавить
        </Button>
        <Button
          onClick={onRandomSelect}
          size="lg"
          disabled={games.length === 0}
          className="flex-1 h-14 text-base rounded-full bg-secondary hover:bg-secondary/80"
        >
          Выбрать случайно
        </Button>
        <Button
          onClick={() => setIsHistoryOpen(true)}
          size="lg"
          variant="outline"
          className="h-14 px-4"
        >
          <History className="w-5 h-5" />
        </Button>
      </div>

      <AddGameDialog
        open={isAddDialogOpen}
        onOpenChange={setIsAddDialogOpen}
        onAddGame={(name) => {
          onAddGame(name);
          setIsAddDialogOpen(false);
        }}
      />

      <Dialog open={!!winner} onOpenChange={() => onCloseWinner()}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-2xl">
              <Trophy className="w-6 h-6 text-yellow-500" />
              Победитель выбран!
            </DialogTitle>
          </DialogHeader>
          <div className="py-6 text-center">
            <p className="text-lg text-muted-foreground mb-3">Победа:</p>
            <p className="text-3xl font-bold text-primary">{winner}</p>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={isHistoryOpen} onOpenChange={setIsHistoryOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <History className="w-5 h-5" />
              История выборов
            </DialogTitle>
          </DialogHeader>
          <ScrollArea className="max-h-80">
            {history.length === 0 ? (
              <p className="text-muted-foreground text-center py-4">История пуста</p>
            ) : (
              <div className="space-y-2">
                {reversedHistory.map((item, index) => (
                  <div
                    key={index}
                    className="flex items-center gap-3 p-3 bg-muted rounded-lg"
                  >
                    <span className="text-muted-foreground text-sm">#{index + 1}</span>
                    <span className="text-foreground font-medium">{item.gameName}</span>
                  </div>
                ))}
              </div>
            )}
          </ScrollArea>
        </DialogContent>
      </Dialog>

      <AlertDialog open={!!gameToDelete} onOpenChange={() => setGameToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Удалить игру?</AlertDialogTitle>
            <AlertDialogDescription>
              Вы уверены, что хотите удалить "{gameToDelete?.name}"? Это действие нельзя отменить.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (gameToDelete) {
                  onDeleteGame(gameToDelete.id);
                  setGameToDelete(null);
                }
              }}
              className="bg-destructive hover:bg-destructive/90"
            >
              Удалить
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </main>
  );
};

export default GamesList;
