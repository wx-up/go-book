package user

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test(t *testing.T) {
	var err error
	fmt.Println(err == nil)
	_, ok := err.(gin.Error)
	fmt.Println(ok)
}
