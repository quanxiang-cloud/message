package restful

import (
	"context"
	"net/http"

	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/internal/core"
	"github.com/quanxiang-cloud/message/internal/service"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router 路由
type Router struct {
	c      *config.Config
	engine *gin.Engine
}

// NewRouter 开启路由
func NewRouter(ctx context.Context, c *config.Config, rf []RouterOption) (*Router, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	v1 := engine.Group("/api/v1/message")
	for _, fn := range rf {
		fn(v1)
	}
	// 创建跟消息相关的路由
	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

// Run 启动服务
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close 关闭服务
func (r *Router) Close() {

}

type RouterOption func(*gin.RouterGroup) error

func WithSender(cz *service.CacheZone, manager *wm.Manager) RouterOption {
	return func(g *gin.RouterGroup) error {
		g.POST("/publish", func(c *gin.Context) {
			req := new(service.PublishReq)
			err := c.ShouldBind(req)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			_, err = cz.Publish(context.Background(), req)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			resp.Format(nil, nil)
		})

		g.POST("/write", func(c *gin.Context) {
			req := new(wm.SendReq)
			err := c.ShouldBind(req)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			_, err = manager.Send(context.Background(), req)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			resp.Format(nil, nil)
		})
		return nil
	}
}

func WithBus(bus *core.Bus) RouterOption {
	return func(g *gin.RouterGroup) error {
		g.POST("/send", func(c *gin.Context) {
			message := new(core.Message)
			err := c.ShouldBindJSON(message)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			resp, err := bus.Send(c.Request.Context(), message)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}

			c.JSON(http.StatusOK, resp)
		})

		g.POST("/send/batch", func(c *gin.Context) {
			type batch []*core.Message

			batchData := make(batch, 0)
			err := c.ShouldBindJSON(batchData)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}

			for _, message := range batchData {
				_, _ = bus.Send(c.Request.Context(), message)
			}
			c.JSON(http.StatusOK, nil)
		})
		return nil
	}
}

func WithWebSocket(ctx context.Context, ws *Websocket) RouterOption {
	return func(rg *gin.RouterGroup) error {
		rg.GET("/ws", ws.Handler)
		return nil
	}
}
