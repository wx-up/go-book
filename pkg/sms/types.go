package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplId string, params []string, phones ...string) error
}
