package grpc

import "context"

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, request *CreateUserRequest) (*CreateUserResponse, error) {
	// TODO implement me
	panic("implement me")
}
