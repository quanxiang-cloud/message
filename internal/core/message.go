package core

import (
	"context"
	"fmt"

	daprd "github.com/dapr/go-sdk/client"
	"github.com/quanxiang-cloud/message/pkg/component/dapr"
)

//go:generate stringer -type Channel
type Channel int

const (
	Letter Channel = iota
	Email
)

type Message struct {
	Channel Channel   `json:"channel,omitempty"`
	Data    dapr.Data `json:"data,omitempty"`
}

type SendResp struct{}

type Bus struct {
	daprClient daprd.Client

	pubsubName string
	tenant     string
}

func New(ctx context.Context, opts ...Option) (*Bus, error) {
	client, err := daprd.NewClient()
	if err != nil {
		return nil, err
	}
	bus := &Bus{
		daprClient: client,
	}

	for _, fn := range opts {
		fn(bus)
	}
	return bus, nil
}

type Option func(*Bus) error

func WithPubsubName(pubsubName string) Option {
	return func(b *Bus) error {
		b.pubsubName = pubsubName
		return nil
	}
}

func WithTenant(tenant string) Option {
	return func(b *Bus) error {
		b.tenant = tenant
		return nil
	}
}

func (b *Bus) Send(ctx context.Context, req *Message) (*SendResp, error) {
	topic := fmt.Sprintf("%s.%s", b.tenant, req.Channel)
	fmt.Println("topic:", topic)
	fmt.Println("pubsubName:", b.pubsubName)
	if err := b.daprClient.PublishEvent(context.Background(), b.pubsubName, topic, req.Data); err != nil {
		return &SendResp{}, err
	}
	return &SendResp{}, nil
}

func (b *Bus) Close() error {
	b.daprClient.Close()
	return nil
}
