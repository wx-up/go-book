package tencent

import (
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func TestService_Send(t *testing.T) {
	client, err := sms.NewClient(common.NewCredential("", ""), regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}
	svc := NewService(client, "", "")
	_ = svc
}
