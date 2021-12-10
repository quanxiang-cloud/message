package component

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/pkg/component/dapr"
)

type Sender interface {
	Scaffold(context.Context, dapr.Data) error
}

type Component struct {
	e *gin.RouterGroup

	sender Sender
}

func New(ctx context.Context, sender Sender, opts ...Option) *Component {
	c := &Component{
		sender: sender,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.init(ctx, sender)
	return c
}

func (c *Component) init(ctx context.Context, sender Sender) {
	c.e.POST("/send", func(ctx *gin.Context) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		event := new(dapr.DaprEvent)
		err = json.Unmarshal(body, event)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		err = c.sender.Scaffold(ctx, event.Data)
		if err != nil {
			errHandle(ctx, err)
			return
		}
	})
}

func errHandle(c *gin.Context, err error) {
	log.Println(err.Error())
	c.JSON(http.StatusOK, nil)
}

type Option func(*Component)

func WithRouter(group *gin.RouterGroup) Option {
	return func(c *Component) {
		c.e = group
	}
}
