import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { UserPlus, Edit, LogOut, AlertCircle, Trash2 } from "lucide-react";
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

interface RoomSettingsProps {
  roomName: string;
  isOwner: boolean;
  onAddUser: (userId: string) => void;
  onRenameRoom: (newName: string) => void;
  onDeleteRoom: () => void;
  onLeaveRoom: () => void;
}

const RoomSettings = ({
  roomName,
  isOwner,
  onAddUser,
  onRenameRoom,
  onDeleteRoom,
  onLeaveRoom,
}: RoomSettingsProps) => {
  const [username, setUsername] = useState("");
  const [newRoomName, setNewRoomName] = useState(roomName);
  const [isLeaveDialogOpen, setIsLeaveDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const handleAddUser = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      onAddUser(username.trim());
      setUsername("");
    }
  };

  const handleRenameRoom = (e: React.FormEvent) => {
    e.preventDefault();
    if (newRoomName.trim() && newRoomName.trim() !== roomName) {
      onRenameRoom(newRoomName.trim());
    }
  };

  return (
    <aside className="w-96 bg-background border-l border-border flex flex-col h-screen">
      <div className="px-8 py-5 border-b border-border">
        <h2 className="text-xl font-semibold text-foreground">Настройки комнаты</h2>
      </div>

      <ScrollArea className="flex-1 px-8 py-6">
        <div className="space-y-6">
          {/* Добавить пользователя */}
          <div className="space-y-3">
            <div className="flex items-center gap-2 text-foreground">
              <UserPlus className="w-5 h-5" />
              <h3 className="font-semibold">Добавить пользователя</h3>
            </div>
            <form onSubmit={handleAddUser} className="space-y-3">
              <div className="space-y-2">
                <Label htmlFor="username" className="text-sm text-muted-foreground">
                  Имя пользователя
                </Label>
                <Input
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="Введите имя пользователя"
                  className="bg-card border-border"
                />
              </div>
              <Button
                type="submit"
                disabled={!username.trim()}
                className="w-full"
                size="sm"
              >
                Добавить
              </Button>
            </form>
          </div>

          {/* Переименовать комнату - только для владельцев */}
          {isOwner && (
            <>
              <Separator className="bg-border" />

              <div className="space-y-3">
                <div className="flex items-center gap-2 text-foreground">
                  <Edit className="w-5 h-5" />
                  <h3 className="font-semibold">Переименовать комнату</h3>
                </div>
                <form onSubmit={handleRenameRoom} className="space-y-3">
                  <div className="space-y-2">
                    <Label htmlFor="room-name" className="text-sm text-muted-foreground">
                      Название комнаты
                    </Label>
                    <Input
                      id="room-name"
                      value={newRoomName}
                      onChange={(e) => setNewRoomName(e.target.value)}
                      placeholder="Введите новое название"
                      className="bg-card border-border"
                    />
                  </div>
                  <Button
                    type="submit"
                    disabled={!newRoomName.trim() || newRoomName.trim() === roomName}
                    className="w-full"
                    size="sm"
                  >
                    Переименовать
                  </Button>
                </form>
              </div>
            </>
          )}

          <Separator className="bg-border" />

          {/* Выйти из комнаты */}
          <div className="space-y-3">
            <div className="flex items-center gap-2 text-foreground">
              <LogOut className="w-5 h-5" />
              <h3 className="font-semibold">Выйти из комнаты</h3>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">
                Вы покинете эту комнату и больше не будете видеть её в списке.
              </p>
              <Button
                onClick={() => setIsLeaveDialogOpen(true)}
                variant="destructive"
                className="w-full"
                size="sm"
              >
                <LogOut className="w-4 h-4 mr-2" />
                Выйти из комнаты
              </Button>
            </div>
          </div>

          {/* Удалить комнату - только для владельцев */}
          {isOwner && (
            <>
              <Separator className="bg-border" />

              <div className="space-y-3">
                <div className="flex items-center gap-2 text-destructive">
                  <Trash2 className="w-5 h-5" />
                  <h3 className="font-semibold">Удалить комнату</h3>
                </div>
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    Комната будет удалена навсегда. Это действие нельзя отменить.
                  </p>
                  <Button
                    onClick={() => setIsDeleteDialogOpen(true)}
                    variant="destructive"
                    className="w-full"
                    size="sm"
                  >
                    <Trash2 className="w-4 h-4 mr-2" />
                    Удалить комнату
                  </Button>
                </div>
              </div>
            </>
          )}
        </div>
      </ScrollArea>

      {/* Диалог подтверждения выхода */}
      <AlertDialog open={isLeaveDialogOpen} onOpenChange={setIsLeaveDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-destructive" />
              Подтвердите действие
            </AlertDialogTitle>
            <AlertDialogDescription>
              Вы уверены, что хотите выйти из комнаты "{roomName}"? Вы больше не
              увидите эту комнату в списке.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction
              onClick={onLeaveRoom}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Выйти
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Диалог подтверждения удаления */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle className="flex items-center gap-2">
              <AlertCircle className="w-5 h-5 text-destructive" />
              Удалить комнату?
            </AlertDialogTitle>
            <AlertDialogDescription>
              Вы уверены, что хотите удалить комнату "{roomName}"? Все игры и голоса
              будут удалены. Это действие нельзя отменить.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction
              onClick={onDeleteRoom}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Удалить
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </aside>
  );
};

export default RoomSettings;
