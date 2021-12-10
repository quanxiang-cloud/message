package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/pkg/component"
	"github.com/quanxiang-cloud/message/pkg/component/letter"
)

func main() {
	var host string
	var port string

	flag.StringVar(&host, "message-server", "", "message server host")
	flag.StringVar(&port, "port", ":80", "")
	flag.Parse()

	ctx := context.Background()
	sender, err := letter.New(ctx, host)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	_ = component.New(
		context.Background(),
		sender,
		component.WithRouter(e.Group("")),
	)

	fmt.Println("start...")
	err = e.Run(port)
	if err != nil {
		panic(err)
	}
}
