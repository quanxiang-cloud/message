package websocket

import (
	"context"
	"time"
)

type Object struct {
	ID   string
	UUID string
	Time time.Time
}

type Affair interface {
	Create(ctx context.Context, obj Object) error
	Renewal(ctx context.Context, obj Object) error
	Delete(ctx context.Context, obj Object) error
}

func CopyFromConnect(conn *Connect) Object {
	return Object{
		ID:   conn.id,
		UUID: conn.uuid,
		Time: time.Now(),
	}
}
