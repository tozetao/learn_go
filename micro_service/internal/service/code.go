package service

import (
	"context"
	"learn_go/webook/internal/repository"
	"learn_go/webook/internal/service/sms"
)

// 短信服务 => 验证码服务 => 登录

// 验证码服务安全分析：
// 1.每个手机号每间隔1分钟发一次。
// 2.验证码的有效期为10分钟
// 3.验证码不能被暴力破解

var (
	ErrTooManyVerify = repository.ErrTooManyVerify
	ErrTooManySend   = repository.ErrTooManySend
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo       *repository.CodeRepository
	smsSvc     sms.Service
	templateId string
}

/*
什么时候需要定义错误?
当所实现的接口需要对外暴漏明确的错误时就可以定义新的错误类型。简单来说需要调用方处理处理时就可以考虑定义新的错误类型。

关于Service（业务层）的错误返回值。
	对于业务层的接口，只有在需要给调用方暴漏明确的错误时，service才重新需要重新定义错误。
	比如下面Service的Send接口，调用方只需要知道发送结果是成功或失败的；而Verify接口不一样，调用方需要知道
*/

func NewCodeService(templateId string, smsSvc sms.Service, repo *repository.CodeRepository) CodeService {
	return &codeService{
		templateId: templateId,
		repo:       repo,
		smsSvc:     smsSvc,
	}
}

// Send
// biz string: 表示业务模块。
func (c *codeService) Send(ctx context.Context, biz string, phone string, code string) error {
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = c.smsSvc.Send(ctx, c.templateId, []string{code}, []string{phone})
	if err != nil {
		// 服务异常，需要记录日志并提醒系统。
		return err
	}
	return nil
}

func (c *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	// 对调用方频闭了验证过多的错误
	if err == ErrTooManyVerify {
		// 异常点，正常来说不会有过多的验证错误的，可考虑记录日志。
		return false, nil
	}
	return ok, err
}
