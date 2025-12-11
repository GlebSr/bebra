package profile

import "time"

type User struct {
	ID           string    `json:"id"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}
