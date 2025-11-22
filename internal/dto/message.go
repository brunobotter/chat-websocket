package dto

import "time"

// MessageInterface defines the contract for a message.
type MessageInterface interface {
	GetUser() string
	GetContent() string
	GetTimestamp() time.Time
	GetRoomID() string
	GetTarget() string
}

type Message struct {
	User      string    `json:"user"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RoomID    string    `json:"room_id,omitempty"`
	Target    string    `json:"target,omitempty"`
}

func (m *Message) GetUser() string      { return m.User }
func (m *Message) GetContent() string   { return m.Content }
func (m *Message) GetTimestamp() time.Time { return m.Timestamp }
func (m *Message) GetRoomID() string    { return m.RoomID }
func (m *Message) GetTarget() string    { return m.Target }

type IncomingInterface interface {
	GetContent() string
	GetTarget() string
}

type Incoming struct {
	Content string `json:"content"`
	Target  string `json:"target"`
}

func (i *Incoming) GetContent() string { return i.Content }
func (i *Incoming) GetTarget() string  { return i.Target }
