package dto

import "time"

// Message represents a chat message with user, content, timestamp, room and target info.
type Message struct {
	User      string    `json:"user"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RoomID    string    `json:"room_id,omitempty"`
	Target    string    `json:"target,omitempty"`
}

// Incoming represents the structure for incoming messages from clients.
type Incoming struct {
	Content string `json:"content"`
	Target  string `json:"target"`
}
