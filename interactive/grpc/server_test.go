package grpc

import (
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"testing"
)

func Test_Validate(t *testing.T) {
	req := &inter.GetRequest{
		Biz: "",
	}
	v, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(req); err != nil {

		fmt.Println("validation failed:")
	} else {
		fmt.Println("validation succeeded")
	}
}
