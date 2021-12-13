package websocket

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/websocket"
)

type Manager struct {
	sync.RWMutex
	pool   map[string][]*Connect
	affair Affair
	log    logr.Logger
}

func NewManager(ctx context.Context, affair Affair, log logr.Logger) (*Manager, error) {
	return &Manager{
		pool:   make(map[string][]*Connect),
		affair: affair,
		log:    log.WithName("manager"),
	}, nil
}

func (m *Manager) Register(ctx context.Context, id string, ws *websocket.Conn) (*Connect, error) {
	ctx, cacel := context.WithCancel(context.Background())
	conns := m.getConns(id)
	if len(conns) == 0 {
		conns = make([]*Connect, 0, 1)
	}
	conn := NewConn(ctx, id, ws)
	conns = append(conns, conn)

	m.upsetConns(id, conns)

	ws.SetCloseHandler(func(code int, text string) error {
		cacel()
		m.UnRegister(ctx, conn)
		m.log.Info("UnRegister", "id", id, "uuid", conn.GetUUID())
		return nil
	})

	ws.SetPingHandler(func(appData string) error {
		err := m.affair.Renewal(ctx, CopyFromConnect(conn))
		if err != nil {
			return err
		}
		err = ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	err := m.affair.Create(ctx, CopyFromConnect(conn))

	m.read(ctx, *conn)

	m.log.Info("Register", "id", id, "uuid", conn.GetUUID())
	return conn, err
}

func (m *Manager) read(ctx context.Context, conn Connect) {
	go func(ctx context.Context, conn Connect) {
		rc := conn.Read()
		for {
			select {
			case _, ok := <-rc:
				if !ok {
					return
				}
			case <-ctx.Done():
				conn.socket.Close()
				return
			}

		}
	}(ctx, conn)
}

func (m *Manager) UnRegister(ctx context.Context, conn *Connect) error {
	id := conn.id

	conns := m.getConns(id)
	for i, elem := range conns {
		if elem.uuid == conn.id {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	m.upsetConns(id, conns)

	return m.affair.Delete(ctx, CopyFromConnect(conn))
}

func (m *Manager) Renewal(ctx context.Context, conn *Connect) error {
	return m.affair.Renewal(ctx, CopyFromConnect(conn))
}

type SendReq struct {
	ID      string
	UUID    []string
	Content []byte
}

type SendResp struct {
}

func (m *Manager) Send(ctx context.Context, req *SendReq) (*SendResp, error) {
	conns := m.getConns(req.ID, req.UUID...)
	for _, conn := range conns {
		err := conn.Write(websocket.TextMessage, req.Content)
		if err != nil {
			m.log.Error(err, "writh message", "id", conn.id, "uuid", conn.uuid)
			_ = conn.socket.Close()
		}
	}
	return &SendResp{}, nil
}

func (m *Manager) getConns(id string, uuids ...string) []*Connect {
	m.RLock()
	defer m.RUnlock()
	if l := len(uuids); l != 0 {
		if l == 1 && len(m.pool[id]) == 1 &&
			m.pool[id][0].uuid == uuids[0] {
			return []*Connect{m.pool[id][0]}
		}

		conns := make([]*Connect, 0, l)
		for _, uuid := range uuids {
			for _, conn := range m.pool[id] {
				if conn.uuid == uuid {
					conns = append(conns, conn)
					break
				}
			}
		}
		return conns
	}
	return m.pool[id]
}

func (m *Manager) upsetConns(id string, conns []*Connect) {
	m.Lock()
	defer m.Unlock()

	if len(conns) == 0 {
		delete(m.pool, id)
		return
	}
	m.pool[id] = conns
}
