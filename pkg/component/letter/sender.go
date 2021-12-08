// package letter

package letter

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/message/pkg/component/dapr"
)

type Letter struct {
}

func New() (*Letter, error) {
	return &Letter{}, nil
}

func (l *Letter) Send(ctx context.Context, data dapr.Data) error {
	fmt.Println("----")
	return nil
}
