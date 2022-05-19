package core

import (
	"context"
	"fmt"

	daprd "github.com/dapr/go-sdk/client"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/pkg/component/event"
)

//go:generate stringer -type Channel
type Channel int

const (
	None Channel = iota
	Letter
	Email
)

type Message struct {
	event.Data `json:",omitempty"`
}

type SendResp struct{}

type Bus struct {
	daprClient daprd.Client
	log        logr.Logger

	pubsubName string
	tenant     string
}

func New(ctx context.Context, log logr.Logger, opts ...Option) (*Bus, error) {
	client, err := InitDaprClientIfNil()
	if err != nil {
		return nil, err
	}
	bus := &Bus{
		daprClient: client,
		log:        log.WithName("bus"),
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
	var topic string

	if req.Data.LetterSpec != nil {
		topic = fmt.Sprintf("%s.%s", b.tenant, Letter.String())
		if err := b.publish(ctx, topic, req.Data); err != nil {
			b.log.Error(err, "push letter", "userID", req.ID)
			return &SendResp{}, err
		}
	}

	if req.Data.EmailSpec != nil {
		topic = fmt.Sprintf("%s.%s", b.tenant, Email.String())
		if err := b.publish(ctx, topic, req.Data); err != nil {
			b.log.Error(err, "push email", "title", req.EmailSpec.Title)
			return &SendResp{}, err
		}
	}

	b.log.Info("publish success")
	return &SendResp{}, nil
}

func (b *Bus) publish(ctx context.Context, topic string, data interface{}) error {
	b.log.Info("send message", "topic", topic)
	if err := b.daprClient.PublishEvent(context.Background(), b.pubsubName, topic, data); err != nil {
		b.log.Error(err, "publishEvent", "topic", topic, "pubsubName", b.pubsubName)
		return err
	}
	return nil
}

func (b *Bus) Close() error {
	b.daprClient.Close()
	return nil
}
