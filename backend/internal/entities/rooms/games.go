package rooms

import "time"

type Game struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func (g Game) IsValid() bool {
	return g.ID != "" && g.RoomID != "" && g.Title != ""
}

type Vote struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	GameID    string    `json:"game_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Result struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	GameID    string    `json:"game_id"`
	ChosenBy  string    `json:"chosen_by"`
	CreatedAt time.Time `json:"created_at"`
}
