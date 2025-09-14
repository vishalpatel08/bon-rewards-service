package models

import "time"

type Reward struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Description string    `json:"description"`
	IssuedAt    time.Time `json:"issued_at"`
}
