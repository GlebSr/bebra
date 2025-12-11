package rooms

import "time"

type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Room) IsValid() bool {
	return r.ID != "" && r.Name != "" && r.OwnerID != ""
}

type RoomParticipant struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (rp RoomParticipant) IsValid() bool {
	return rp.ID != "" && rp.RoomID != "" && rp.UserID != "" && (rp.Role == "member" || rp.Role == "owner")
}
