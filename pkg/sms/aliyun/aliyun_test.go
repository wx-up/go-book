package aliyun

import (
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
)

func TestService_Send(t *testing.T) {
	client, err := sms.NewClient(&openapi.Config{
		AccessKeyId:     nil,
		AccessKeySecret: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(client, "")
	_ = svc
}
