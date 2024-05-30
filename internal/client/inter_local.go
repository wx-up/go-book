package client

import (
	"context"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"github.com/wx-up/go-book/interactive/domain"
	"github.com/wx-up/go-book/interactive/service"
	"google.golang.org/grpc"
)

// InteractiveServiceAdapter 将本地实现伪装成一个 gRPC 客户端
type InteractiveServiceAdapter struct {
	svc service.InteractiveService
}

func NewInteractiveServiceAdapter(svc service.InteractiveService) *InteractiveServiceAdapter {
	return &InteractiveServiceAdapter{
		svc: svc,
	}
}

func (i *InteractiveServiceAdapter) IncrReadCnt(ctx context.Context, in *inter.IncrReadCntRequest, opts ...grpc.CallOption) (*inter.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	if err != nil {
		return nil, err
	}
	return &inter.IncrReadCntResponse{}, nil
}

func (i *InteractiveServiceAdapter) Get(ctx context.Context, in *inter.GetRequest, opts ...grpc.CallOption) (*inter.GetResponse, error) {
	res, err := i.svc.Get(ctx, in.GetBiz(), in.GetId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &inter.GetResponse{
		Inter: i.toDTO(res),
	}, nil
}

// toDTO DTO数据传输对象
func (i *InteractiveServiceAdapter) toDTO(interactive domain.Interactive) *inter.Interactive {
	return &inter.Interactive{}
}
