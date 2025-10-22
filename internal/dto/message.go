package dto

import "time"

type Message struct {
	User      string    `json:"user"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RoomID    string    `json:"room_id,omitempty"`
	Target    string    `json:"target,omitempty"`
}

type Incoming struct {
	Content string `json:"content"`
	Target  string `json:"target"`
}
