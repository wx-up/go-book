package web

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func Get() any {
	return nil
}

func Test(t *testing.T) {
	err := Get()
	fmt.Println(err == nil)
	_, ok := err.(gin.Error)
	fmt.Println(ok)
}
