package ioc

import (
	"github.com/spf13/viper"
	"github.com/wx-up/go-book/interactive/grpc"
	"github.com/wx-up/go-book/pkg/grpcx"
	grpc2 "google.golang.org/grpc"
)

func CreateGRPCServer(inter *grpc.InteractiveServiceServer) *grpcx.Server {
	s := grpc2.NewServer()
	inter.RegisterServer(s)
	addr := viper.GetString("grpc.server.addr")
	return &grpcx.Server{
		Addr:   addr,
		Server: s,
	}
}
