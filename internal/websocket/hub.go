package websocket

import (
	"context"
	"encoding/json"

	"github.com/brunobotter/chat-websocket/internal/dto"
	"go.uber.org/zap"
)

type IClient interface {
	GetRoomID() string
	SetRoomID(string)
	GetUser() string
	SendMessage([]byte)
	GetSendChannel() chan []byte
}

type Hub struct {
	Rooms      map[string]map[IClient]bool
	Broadcast  chan dto.Message
	Register   chan IClient
	Unregister chan IClient
	logger     *zap.Logger
	ChatStore  ChatStore
}

func NewHub(logger *zap.Logger, chatStore ChatStore) *Hub {
	return &Hub{
		Rooms:      make(map[string]map[IClient]bool),
		Broadcast:  make(chan dto.Message),
		Register:   make(chan IClient),
		Unregister: make(chan IClient),
		logger:     logger,
		ChatStore:  chatStore,
	}
}

func (h *Hub) Run() {
	ctx := context.Background()
	for {
		select {
		//registro clientes
		case client := <-h.Register:
			rid := client.GetRoomID()
			if rid == "" {
				rid = "default"
				client.SetRoomID(rid)
			}
			if _, ok := h.Rooms[rid]; !ok {
				h.Rooms[rid] = make(map[IClient]bool)
			}
			h.Rooms[rid][client] = true
			if h.ChatStore != nil {
				go func(c IClient) {

					unread, err := h.ChatStore.GetUnreadMessages(ctx, c.GetUser())
					if err != nil {
						h.logger.Error("Falha ao buscar mensagens não lidas", zap.String("user", c.GetUser()), zap.Error(err))
						return
					}
					for _, msg := range unread {
						payload, _ := json.Marshal(msg)
						c.SendMessage(payload)
					}
					_ = h.ChatStore.ClearUnread(ctx, c.GetUser())
				}(client)
			}

//deregistro de clientes
case client := <-h.Unregister:
	rid := client.GetRoomID()
	if clients, ok := h.Rooms[rid]; ok {
	    delete(clients, client)
        close(client.GetSendChannel())
        if len(clients) == 0 {
            delete(h.Rooms, rid)
        }
    }
//recebimento de mensagens
case msg := <-h.Broadcast:
    //mensagens privadas
    if msg.Target != "" {
        for _, clients := range h.Rooms {
            for client := range clients {
                if client.GetUser() == msg.Target {
                    select {
                    case client.GetSendChannel() <- []byte(msg.Content):
                    default:
                        close(client.GetSendChannel())
                        delete(clients, client)
                    }
                }
            }
        }
        if h.ChatStore != nil {
            _ = h.ChatStore.SaveUnread(ctx, msg.Target, msg)
        }
        continue
    }
    //broadcast por sala
    if clients, ok := h.Rooms[msg.RoomID]; ok {
        for client := range clients {
            select {
            case client.GetSendChannel() <- []byte(msg.Content):
            default:
                delete(clients, client)
                close(client.GetSendChannel())
            }
        }
    }
    if h.ChatStore != nil {
        _ = h.ChatStore.SaveMessage(ctx, msg.RoomID, msg, 50)
    }
}	}	}	}	