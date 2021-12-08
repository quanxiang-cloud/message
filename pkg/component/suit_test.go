package component

import (
	"context"
	"testing"

	"github.com/quanxiang-cloud/message/pkg/component/letter"
)

func TestComponent(t *testing.T) {
	sender, err := letter.New()
	if err != nil {
		t.Fatal(err)
		return
	}

	c := New(context.Background(), sender)
	t.Log("start...")
	err = c.Start(":8080")
	if err != nil {
		t.Fatal(err)
	}
}
