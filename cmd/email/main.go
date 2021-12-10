package main

import (
	"context"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/message/pkg/component"
	"github.com/quanxiang-cloud/message/pkg/component/email"
)

func main() {
	var port string

	flag.StringVar(&port, "port", ":80", "")
	email.Prepare()
	flag.Parse()

	ctx := context.Background()
	sender, err := email.New(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	_ = component.New(
		context.Background(),
		sender,
		component.WithRouter(e.Group("")),
	)

	log.Println("start...")
	err = e.Run(port)
	if err != nil {
		panic(err)
	}
}
