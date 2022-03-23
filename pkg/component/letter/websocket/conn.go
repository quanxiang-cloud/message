package websocket

import (
	"context"

	"github.com/gorilla/websocket"
	id2 "github.com/quanxiang-cloud/cabin/id"
)

type Connect struct {
	ctx    context.Context
	id     string
	uuid   string
	socket *websocket.Conn
}

func NewConn(ctx context.Context, id string, conn *websocket.Conn) *Connect {
	uuid := id2.BaseUUID()
	return &Connect{
		ctx:    ctx,
		id:     id,
		uuid:   uuid,
		socket: conn,
	}
}

func (c *Connect) GetUUID() string {
	return c.uuid
}

func (c *Connect) Read() <-chan []byte {
	mc := make(chan []byte)
	go func(mc chan<- []byte) {
		defer close(mc)
		for {
			_, message, err := c.socket.ReadMessage()
			if err != nil {
				break
			}
			mc <- message
		}
	}(mc)

	return mc
}

func (c *Connect) Write(messageType int, data []byte) error {
	return c.socket.WriteMessage(messageType, data)
}
