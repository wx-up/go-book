package grpc

import (
	"context"
	"fmt"
)

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	fmt.Println("我收到了")
	return &CreateUserResponse{}, nil
}
