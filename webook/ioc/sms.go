package ioc

import "learn_go/webook/internal/service/sms"

func NewSMSService() sms.Service {
	return sms.NewMockSMSService()
}
