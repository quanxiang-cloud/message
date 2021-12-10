// package letter

package letter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/message/pkg/client"
	"github.com/quanxiang-cloud/message/pkg/component/dapr"
)

type Letter struct {
	host   string
	client http.Client
}

func New(ctx context.Context, host string) (*Letter, error) {
	return &Letter{
		host: host,
		client: client.New(client.Config{
			Timeout:      time.Second * 20,
			MaxIdleConns: 10,
		}),
	}, nil
}

func (l *Letter) Scaffold(ctx context.Context, data dapr.Data) error {
	if data.LetterSpec == nil {
		return dapr.ErrDataIsNil
	}

	return l.Send(ctx, data.LetterSpec)
}

func (l *Letter) Send(ctx context.Context, data *dapr.LetterSpec) error {
	req := map[string]interface{}{
		"userID":  data.ID,
		"uuid":    data.UUID,
		"content": data.Content,
	}

	err := client.POST(ctx, &l.client, fmt.Sprintf("%s/api/v1/message/publish", l.host), req, nil)
	if err != nil {
		log.Printf("publish fail: %s \n", err.Error())
		return err
	}
	return nil
}
