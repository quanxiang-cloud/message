package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/quanxiang-cloud/message/pkg/component"
	"github.com/quanxiang-cloud/message/pkg/component/letter"
	"go.uber.org/zap"
)

var (
	log logr.Logger
)

func main() {
	var host string
	var port string

	flag.StringVar(&host, "message-server", "", "message server host")
	flag.StringVar(&port, "port", ":80", "")
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	ctx := context.Background()
	sender, err := letter.New(ctx, host, log)
	if err != nil {
		log.Error(err, "new sender")
		panic(err)
	}

	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	_ = component.New(
		context.Background(),
		sender,
		component.WithRouter(e.Group("")),
	)

	log.Info("start...")
	err = e.Run(port)
	if err != nil {
		panic(err)
	}
}
