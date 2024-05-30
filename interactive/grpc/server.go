package grpc

import (
	"context"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"github.com/wx-up/go-book/interactive/domain"
	"github.com/wx-up/go-book/interactive/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InteractiveServiceServer 它不会处理业务逻辑，业务逻辑的处理都是交给 service.InteractiveService 去处理
// 它只负责和grpc客户端通信，并返回相应的结果
// 定位和 http 中的 handler 一致
type InteractiveServiceServer struct {
	inter.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewInteractiveServiceServer(svc service.InteractiveService) *InteractiveServiceServer {
	return &InteractiveServiceServer{
		svc: svc,
	}
}

func (i *InteractiveServiceServer) RegisterServer(server *grpc.Server) {
	inter.RegisterInteractiveServiceServer(server, i)
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, req *inter.IncrReadCntRequest) (*inter.IncrReadCntResponse, error) {
	// 可以考虑增加参数验证
	if req.GetBizId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "biz_id 为空")
	}
	err := i.svc.IncrReadCnt(ctx, req.GetBiz(), req.GetBizId())
	return &inter.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, req *inter.GetRequest) (*inter.GetResponse, error) {
	res, err := i.svc.Get(ctx, req.GetBiz(), req.GetId(), req.GetUid())
	if err != nil {
		return nil, err
	}
	return &inter.GetResponse{Inter: i.toDTO(res)}, nil
}

// toDTO DTO数据传输对象
func (i *InteractiveServiceServer) toDTO(interactive domain.Interactive) *inter.Interactive {
	return &inter.Interactive{}
}
