package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"runtime"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	runtime.Gosched()
	cc, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := client.CreateUser(ctx, &CreateUserRequest{})
	require.NoError(t, err)
	t.Log(resp)
}
