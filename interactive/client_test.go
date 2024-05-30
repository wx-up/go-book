package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

func Test_GRPCClient(t *testing.T) {
	cc, err := grpc.Dial("localhost:8081",
		// 不使用 https
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	client := inter.NewInteractiveServiceClient(cc)
	fmt.Println(client.Get(context.TODO(), &inter.GetRequest{
		Biz: "",
		Id:  0,
		Uid: 0,
	}))
}
