package ioc

import (
	"learn_go/webook/internal/service/oauth2/wechat"
	"os"
)

func InitOAuth2Service() wechat.OAuth2Service {
	appId := os.Getenv("app_id")
	return wechat.NewOAuth2WechatService(appId)
}
