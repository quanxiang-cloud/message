package restful

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
)

type Websocket struct {
	manager *wm.Manager
}

func NewWebsocket(ctx context.Context, conf *config.Config, manager *wm.Manager) (*Websocket, error) {
	return &Websocket{
		manager: manager,
	}, nil
}

//Handler Handler
func (w *Websocket) Handler(c *gin.Context) {
	ctx := context.Background()

	wsConn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	id := c.Request.Header.Get("Id")
	client, err := w.manager.Register(ctx, id, wsConn)
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pong, err := json.Marshal(struct {
		UUID string `json:"uuid"`
	}{
		UUID: client.GetUUID(),
	})
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = w.manager.Send(ctx, &wm.SendReq{
		ID:      id,
		Content: pong,
	})
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
