package main

import (
	"context"
	"flag"

	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/quanxiang-cloud/message/api/restful"
	"github.com/quanxiang-cloud/message/internal/core"
	"github.com/quanxiang-cloud/message/internal/service"
	"github.com/quanxiang-cloud/message/package/config"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
)

func main() {
	var pubsubName string
	var tenant string
	var configPath string

	flag.StringVar(&tenant, "tenant", "default", "Tenant ID.")
	flag.StringVar(&pubsubName, "pubsub-name", "default", "The dapr pubsub component name.")
	flag.StringVar(&configPath, "config", "/configs/config.yml", "config file path")
	flag.Parse()

	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	err = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	bus, err := core.New(ctx,
		core.WithPubsubName(pubsubName),
		core.WithTenant(tenant),
	)
	if err != nil {
		panic(err)
	}

	cz, err := service.NewCacheZone(ctx, conf)
	if err != nil {
		panic(err)
	}

	manager, err := wm.NewManager(ctx, cz)
	if err != nil {
		panic(err)
	}
	ws, err := restful.NewWebsocket(ctx, conf, manager)
	if err != nil {
		panic(err)
	}

	client, err := restful.NewRouter(ctx, conf, []restful.RouterOption{
		restful.WithBus(bus),
		restful.WithWebSocket(ctx, ws),
		restful.WithSender(cz, manager),
	})

	if err != nil {
		panic(err)
	}

	client.Run()
	client.Close()
}
