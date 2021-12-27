package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/quanxiang-cloud/message/api/restful"
	"github.com/quanxiang-cloud/message/internal/core"
	"github.com/quanxiang-cloud/message/internal/service"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
	"go.uber.org/zap"
)

var (
	log logr.Logger
)

func main() {
	var pubsubName string
	var tenant string
	var configPath string

	flag.StringVar(&tenant, "tenant", "default", "Tenant ID.")
	flag.StringVar(&pubsubName, "pubsub-name", "default", "The dapr pubsub component name.")
	flag.StringVar(&configPath, "config", "/configs/config.yml", "config file path")
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Error(err, "get config")
		panic(err)
	}

	ctx := context.Background()

	bus, err := core.New(ctx, log,
		core.WithPubsubName(pubsubName),
		core.WithTenant(tenant),
	)
	if err != nil {
		log.Error(err, "new bus")
		panic(err)
	}

	cz, err := service.NewCacheZone(ctx, conf, log)
	if err != nil {
		log.Error(err, "new cache zone")
		panic(err)
	}

	manager, err := wm.NewManager(ctx, cz, log)
	if err != nil {
		log.Error(err, "new manager")
		panic(err)
	}
	ws, err := restful.NewWebsocket(ctx, conf, manager, log)
	if err != nil {
		log.Error(err, "new webSocket")
		panic(err)
	}

	client, err := restful.NewRouter(ctx, conf, []restful.RouterOption{
		restful.WithBus(bus),
		restful.WithWebSocket(ctx, ws),
		restful.WithSender(cz, manager),
	}, log)

	if err != nil {
		panic(err)
	}

	client.Run()
	client.Close()
}
