package client

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyScaleInteractiveServiceClient struct {
	remote inter.InteractiveServiceClient
	local  inter.InteractiveServiceClient

	threshold *atomicx.Value[int32]
}

func NewGreyScaleInteractiveServiceClient(remote inter.InteractiveServiceClient, local inter.InteractiveServiceClient) *GreyScaleInteractiveServiceClient {
	return &GreyScaleInteractiveServiceClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValueOf[int32](0),
	}
}

func (g *GreyScaleInteractiveServiceClient) OnChange() {
	viper.OnConfigChange(func(in fsnotify.Event) {

	})
	return
}

func (g *GreyScaleInteractiveServiceClient) OnChangeV1(ch <-chan int) {
	go func() {
		for v := range ch {
			g.UpdateThreshold(int32(v))
		}
	}()
	return
}

func (g *GreyScaleInteractiveServiceClient) OnChangeV2() chan<- int {
	ch := make(chan int, 100)
	go func() {
		for v := range ch {
			g.UpdateThreshold(int32(v))
		}
	}()
	return ch
}

func (g *GreyScaleInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *inter.IncrReadCntRequest, opts ...grpc.CallOption) (*inter.IncrReadCntResponse, error) {
	return g.client().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) Get(ctx context.Context, in *inter.GetRequest, opts ...grpc.CallOption) (*inter.GetResponse, error) {
	return g.client().Get(ctx, in, opts...)
}

func (g *GreyScaleInteractiveServiceClient) UpdateThreshold(newThreshold int32) {
	g.threshold.Store(newThreshold)
}

func (g *GreyScaleInteractiveServiceClient) client() inter.InteractiveServiceClient {
	threshold := g.threshold.Load()
	num := rand.Int31n(100)
	if num < threshold {
		return g.remote
	}
	return g.local
}
