package main

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/message/pkg/component"
	"github.com/quanxiang-cloud/message/pkg/component/letter"
)

func main() {
	sender, err := letter.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c := component.New(context.Background(), sender)
	fmt.Println("start...")
	err = c.Start(":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
}
