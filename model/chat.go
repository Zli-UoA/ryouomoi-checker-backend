package model

import "time"

type ChatRoom struct {
	ID        int64
	Couple    *Couple
	CreatedAt *time.Time
}

type Chat struct {
	ID        int64
	Room      *ChatRoom
	User      *User
	Message   string
	CreatedAt *time.Time
}
