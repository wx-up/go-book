package grpc

import (
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	server := grpc.NewServer()
	userServer := &Server{}
	defer func() {
		// 优雅退出
		server.GracefulStop()
	}()

	RegisterUserServiceServer(server, userServer)

	// 创建一个监听器
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
	//server.ServeHTTP()
}
