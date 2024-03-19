package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplId string, params []NameArg, phones ...string) error
	Type() string
}

// NameArg is a name-value pair for template parameters.
// 因为有些短信供应商的模板参数是以key-value形式传入的，还有一些是以切片形式传入的，所以这里统一用NameArg来表示。
type NameArg struct {
	Name string
	Val  string
}
