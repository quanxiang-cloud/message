package restful

import (
	"context"
	"encoding/json"
	"net/http"

	ct "git.internal.yunify.com/qxp/misc/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
	client "github.com/quanxiang-cloud/message/pkg/quanxiang"
)

type Websocket struct {
	manager *wm.Manager

	warden client.Warden
}

func NewWebsocket(ctx context.Context, conf *config.Config, manager *wm.Manager) (*Websocket, error) {
	return &Websocket{
		manager: manager,
		warden: client.NewOWarden(ct.Config{
			Timeout:      conf.InternalNet.Timeout,
			MaxIdleConns: conf.InternalNet.MaxIdleConns,
		}),
	}, nil
}

//Handler Handler
func (w *Websocket) Handler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx := context.Background()

	profile, err := w.warden.CheckToken(ctx, token, config.Conf.AUth.CheckToken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wsConn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	client, err := w.manager.Register(ctx, profile.UserID, wsConn)
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
		ID:      profile.UserID,
		Content: pong,
	})
	if err != nil {
		w.manager.UnRegister(ctx, client)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
