package ioc

import "learn_go/webook/internal/service/sms"

func InitSMSService() sms.Service {
	return sms.NewMockSMSService()
}
