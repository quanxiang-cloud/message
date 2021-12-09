package websocket

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	sync.RWMutex

	pool map[string][]*Connect
}

func NewManager(ctx context.Context) (*Manager, error) {
	return &Manager{
		pool: make(map[string][]*Connect),
	}, nil
}

func (m *Manager) Register(ctx context.Context, id string, ws *websocket.Conn) *Connect {
	conns := m.getConns(id)
	if len(conns) == 0 {
		conns = make([]*Connect, 0, 1)
	}
	conn := NewConn(ctx, id, ws)
	conns = append(conns, conn)

	m.upsetConns(id, conns)

	ws.SetCloseHandler(func(code int, text string) error {
		m.UnRegister(ctx, conn)
		return nil
	})

	ws.SetPingHandler(func(appData string) error {
		return nil
	})
	return conn
}

func (m *Manager) UnRegister(ctx context.Context, conn *Connect) {
	id := conn.id

	conns := m.getConns(id)
	for i, elem := range conns {
		if elem.uuid == conn.id {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	m.upsetConns(id, conns)
}

func (m *Manager) Renewal(ctx context.Context, conn *Connect) {
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
		// FIXME
		_ = err
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
