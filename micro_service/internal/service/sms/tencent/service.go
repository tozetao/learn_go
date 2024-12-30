package tencent

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111" // 引入sms
	"log"
)

type Service struct {
	client        *sms.Client
	clientProfile *profile.ClientProfile
	appId         string
	signName      string
}

func NewService(secretId string, secretKey string) *Service {
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	cpf.HttpProfile.ReqTimeout = 10
	cpf.SignMethod = "HmacSHA1"

	credential := common.NewCredential(secretId, secretKey)
	client, err := sms.NewClient(credential, "ap-guangzhou", cpf)

	// 这里的错误要怎么处理? 记录日志，发出警告
	if err != nil {
		log.Fatal(err)
		return &Service{}
	}

	return &Service{
		clientProfile: cpf,
		client:        client,
	}
}

func (service *Service) Send(ctx context.Context, templateId string, params []string, phones []string) error {
	if service.client == nil {
		return errors.New("SMSService.client not initialized")
	}

	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = common.StringPtr(service.appId)
	request.SignName = common.StringPtr(service.signName)
	request.TemplateId = common.StringPtr(templateId)
	request.TemplateParamSet = common.StringPtrs(params)
	request.PhoneNumberSet = common.StringPtrs(phones)

	response, err := service.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("an API error has returned: %s", err)
	}

	// 更好的做法是返回每个手机对应的错误，怎么封装?
	for _, statusPtr := range response.Response.SendStatusSet {
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("短信发送失败 %s, %s", *status.Code, *status.Message)
		}
	}

	return nil
}
