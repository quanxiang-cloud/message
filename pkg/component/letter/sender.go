// package letter

package letter

import (
	"context"

	"github.com/quanxiang-cloud/message/pkg/component/dapr"
	ws "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
)

type Letter struct {
	manager ws.Manager
}

func New(ctx context.Context) (*Letter, error) {
	return &Letter{}, nil
}

func (l *Letter) Scaffold(ctx context.Context, data dapr.Data) error {
	if data.LetterSpec == nil {
		return dapr.ErrDataIsNil
	}

	return l.Send(ctx, data.LetterSpec)
}

func (l *Letter) Send(ctx context.Context, data *dapr.LetterSpec) error {
	_, err := l.manager.Send(ctx, &ws.SendReq{
		ID:      data.ID,
		UUID:    data.UUID,
		Content: data.Content,
	})
	return err
}
