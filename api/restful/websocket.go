package restful

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/gorilla/websocket"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
)

type Websocket struct {
	manager *wm.Manager
	log     logr.Logger
}

func NewWebsocket(ctx context.Context, conf *config.Config, manager *wm.Manager, log logr.Logger) (*Websocket, error) {
	return &Websocket{
		manager: manager,
		log:     log.WithName("websocket"),
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
		w.log.Error(err, "register")
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
		w.log.Error(err, "json marshal")
		return
	}

	_, err = w.manager.Send(ctx, &wm.SendReq{
		ID:      id,
		UUID:    []string{client.GetUUID()},
		Content: pong,
	})
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		w.log.Error(err, "send first message")
		return
	}
}
