package failover

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("父亲退出了")
		}
	}()

	ctx1, cancel1 := context.WithTimeout(ctx, time.Second*4)
	defer cancel1()
	go func() {
		select {
		case <-ctx1.Done():
			fmt.Println("儿子退出了")
		}
	}()

	time.Sleep(time.Second * 2)
	cancel1()

	time.Sleep(time.Second * 10)
}
