// package letter

package letter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/pkg/client"
	"github.com/quanxiang-cloud/message/pkg/component/event"
)

type Letter struct {
	host   string
	client http.Client

	log logr.Logger
}

func New(ctx context.Context, host string, log logr.Logger) (*Letter, error) {
	return &Letter{
		host: host,
		client: client.New(client.Config{
			Timeout:      time.Second * 20,
			MaxIdleConns: 10,
		}),
		log: log.WithName("letter"),
	}, nil
}

func (l *Letter) Scaffold(ctx context.Context, data event.Data) error {
	if data.LetterSpec == nil {
		return event.ErrDataIsNil
	}

	return l.Send(ctx, data.LetterSpec)
}

func (l *Letter) Send(ctx context.Context, data *event.LetterSpec) error {
	req := map[string]interface{}{
		"userID":  data.ID,
		"uuid":    data.UUID,
		"content": data.Content,
	}

	err := client.POST(ctx, &l.client, fmt.Sprintf("%s/api/v1/message/publish", l.host), req, nil)
	if err != nil {
		l.log.Error(err, "publish", "userID", data.ID, "uuid", data.UUID)
		return err
	}
	return nil
}
