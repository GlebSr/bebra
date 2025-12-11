import { useState } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { LogOut, Plus } from "lucide-react";
import { useCreateRoom } from "@/hooks/useRooms";
import { useToast } from "@/hooks/use-toast";
import CreateRoomDialog from "./CreateRoomDialog";
import type { Room } from "@/api";

interface RoomsListProps {
  rooms: Room[];
  selectedRoomId: string;
  onRoomSelect: (roomId: string) => void;
}

const RoomsList = ({ rooms, selectedRoomId, onRoomSelect }: RoomsListProps) => {
  const { user, logout } = useAuth();
  const { toast } = useToast();
  const [dialogOpen, setDialogOpen] = useState(false);
  const createRoomMutation = useCreateRoom();

  const handleCreateRoom = async (name: string) => {
    try {
      const newRoom = await createRoomMutation.mutateAsync({ name });
      setDialogOpen(false);
      onRoomSelect(newRoom.id);
      toast({
        title: "Комната создана",
        description: `Комната "${name}" успешно создана`,
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Ошибка",
        description: error.data?.error || "Не удалось создать комнату",
      });
    }
  };

  return (
    <aside className="w-60 bg-sidebar border-r border-sidebar-border flex flex-col h-screen">
      <div className="px-6 py-5 border-b border-sidebar-border flex items-center justify-between">
        <h2 className="text-lg font-semibold text-sidebar-foreground">Rooms</h2>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setDialogOpen(true)}
          className="h-8 w-8 text-muted-foreground hover:text-foreground"
        >
          <Plus className="h-4 w-4" />
        </Button>
      </div>
      
      <ScrollArea className="flex-1">
        <div className="px-3 py-2">
          {rooms.length === 0 ? (
            <p className="text-sm text-muted-foreground px-3 py-2">Нет комнат</p>
          ) : (
            rooms.map((room) => (
              <button
                key={room.id}
                onClick={() => onRoomSelect(room.id)}
                className={`w-full text-left px-3 py-2.5 rounded-lg mb-1 transition-colors ${
                  selectedRoomId === room.id
                    ? "bg-sidebar-accent text-sidebar-accent-foreground"
                    : "text-sidebar-foreground hover:bg-sidebar-accent/50"
                }`}
              >
                {room.name}
              </button>
            ))
          )}
        </div>
      </ScrollArea>

      <div className="px-4 py-4 border-t border-sidebar-border flex items-center gap-3">
        <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center">
          <span className="text-sm font-medium">
            {user?.name?.charAt(0).toUpperCase() || 'U'}
          </span>
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-sidebar-foreground truncate">
            {user?.name || 'User'}
          </p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={logout}
          className="h-8 w-8 text-muted-foreground hover:text-foreground"
        >
          <LogOut className="h-4 w-4" />
        </Button>
      </div>

      <CreateRoomDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        onCreateRoom={handleCreateRoom}
        isLoading={createRoomMutation.isPending}
      />
    </aside>
  );
};

export default RoomsList;
