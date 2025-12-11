import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface AddGameDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAddGame: (gameName: string) => void;
}

const AddGameDialog = ({ open, onOpenChange, onAddGame }: AddGameDialogProps) => {
  const [gameName, setGameName] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (gameName.trim()) {
      onAddGame(gameName.trim());
      setGameName("");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Добавить игру</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="game-name">Название игры</Label>
            <Input
              id="game-name"
              value={gameName}
              onChange={(e) => setGameName(e.target.value)}
              placeholder="Введите название игры"
              autoFocus
            />
          </div>
          <div className="flex justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Отмена
            </Button>
            <Button type="submit" disabled={!gameName.trim()}>
              Добавить
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddGameDialog;
