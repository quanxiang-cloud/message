package restful

import (
	"context"
	"net/http"

	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/internal/core"
	"github.com/quanxiang-cloud/message/package/config"
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
	engine.Use(logger.GinLogger(), logger.GinRecovery())

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
				if err != nil {
					c.JSON(http.StatusBadRequest, err)
					return
				}
				return
			}

			c.JSON(http.StatusOK, resp)
		})
		return nil
	}
}
