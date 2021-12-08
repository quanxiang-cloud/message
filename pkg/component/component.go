package component

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/pkg/component/dapr"
)

type Sender interface {
	Send(context.Context, dapr.Data) error
}

type Component struct {
	e *gin.Engine

	sender Sender
}

func New(ctx context.Context, sender Sender) *Component {
	c := &Component{
		sender: sender,
	}
	c.init(ctx, sender)
	return c
}

func (c *Component) Start(port string) error {
	return c.e.Run(port)
}

func (c *Component) init(ctx context.Context, sender Sender) {
	c.e = gin.New()
	c.e.Use(gin.Logger(), gin.Recovery())

	c.e.POST("/send", func(ctx *gin.Context) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		fmt.Println(string(body))
		ctx.AbortWithStatus(http.StatusOK)
		return
		event := new(dapr.DaprEvent)
		err = json.Unmarshal(body, event)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		err = c.sender.Send(ctx, event.Data)
		if err != nil {
			errHandle(ctx, err)
			return
		}
	})
}

func errHandle(c *gin.Context, err error) {
	c.AbortWithError(http.StatusBadRequest, err)
}
