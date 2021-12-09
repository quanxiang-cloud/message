package restful

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
)

type Websocket struct {
	manager *wm.Manager
}

func NewWebsocket(ctx context.Context, manager *wm.Manager) (*Websocket, error) {
	return &Websocket{
		manager: manager,
	}, nil
}

//Handler  Handler
func (w *Websocket) Handler(c *gin.Context) {
	// token := c.Query("token")
	// if token == "" {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	// ctx := logger.CTXTransfer(c)
	// profile, err := w.oauth2s.CheckToken(ctx, token, config.Conf.AUth.CheckToken)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	wsConn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	ctx := context.Background()

	client := w.manager.Register(ctx, "1", wsConn)
	// err = client.SendOut(&comet.Message{
	// 	Content: comet.Content{
	// 		Type:    comet.ConnetSucess,
	// 		Message: map[string]string{"uuid": client.UUID},
	// 	},
	// })

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
		ID:      "1",
		Content: pong,
	})
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
