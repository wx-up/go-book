package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisSlideWindowLimiter_Limit_E2E(t *testing.T) {
	v := time.Second
	fmt.Println(v.Nanoseconds())
}
