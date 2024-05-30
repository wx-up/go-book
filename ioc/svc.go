package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wx-up/go-book/api/proto/gen/inter"
	"github.com/wx-up/go-book/interactive/service"
	"github.com/wx-up/go-book/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"

	"github.com/wx-up/go-book/internal/service/oauth2/wechat"
)

func CreateOAuth2WechatService() wechat.Service {
	// 从环境变量中获取
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		// panic("WECHAT_APP_ID not found in environment variables")
	}

	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		// panic("WECHAT_APP_SECRET not found in environment variables")
	}

	return wechat.NewService(appId, appSecret)
}

func CreateInterGRPCClient(
	svc service.InteractiveService,
) inter.InteractiveServiceClient {
	type Config struct {
		Addr      string `json:"addr"`
		Secure    bool   `json:"secure"`
		Threshold int    `json:"threshold"`
	}
	var c Config
	err := viper.UnmarshalKey("grpc.client.inter", &c)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if c.Secure {
		// 加载证书、启动 https
	} else {

		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(c.Addr, opts...)
	if err != nil {
		panic(err)
	}
	local := client.NewInteractiveServiceAdapter(svc)
	remote := inter.NewInteractiveServiceClient(cc)
	res := client.NewGreyScaleInteractiveServiceClient(remote, local)
	viper.OnConfigChange(func(in fsnotify.Event) {
		type Config struct {
			Addr      string `json:"addr"`
			Secure    bool   `json:"secure"`
			Threshold int32  `json:"threshold"`
		}
		var c Config
		err = viper.UnmarshalKey("grpc.client.inter", &c)
		if err != nil {
			// 打印日志
			return
		}
		res.UpdateThreshold(c.Threshold)
	})
	return res
}
