package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/wx-up/go-book/internal/integration/startup"
)

func Test(t *testing.T) {
	redisClient := startup.InitTestRedis()
	fmt.Println(redisClient.Exists(context.Background(), "test1").Result())
}
