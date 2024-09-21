package sms

import "context"

// 短信服务: 支持发送各种内容。考虑后续可能适配不同的服务商

type Service interface {
	//Send
	//params: 模板占位符号对应的参数。
	Send(ctx context.Context, templateId string, params []string, phones []string) error
}
