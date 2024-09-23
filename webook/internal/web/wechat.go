package web

import (
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/service"
	"learn_go/webook/internal/service/oauth2/wechat"
	"net/http"
)

type OAuth2WechatHandler struct {
	oauth2Svc  wechat.OAuth2Service
	userSvc    service.UserService
	jwtHandler *JWTHandler
}

func NewOAuth2WechatHandler(oauth2Svc wechat.OAuth2Service, userSvc service.UserService, jwtHandler *JWTHandler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		oauth2Svc:  oauth2Svc,
		userSvc:    userSvc,
		jwtHandler: jwtHandler,
	}
}

func (w *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/oauth2/wechat/auth_url", w.Auth2Url)
	server.Any("/oauth2/wechat/verify", w.Callback)
}

// Auth2Url 返回微信OAuth认证地址
func (w *OAuth2WechatHandler) Auth2Url(c *gin.Context) {
	url, err := w.oauth2Svc.Auth2Url(c)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "failed to build auth url.",
		})
		return
	}

	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: url,
	})
}

// Callback 触发回调时的认证
func (w *OAuth2WechatHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	// state := c.Query("state")

	wechatInfo, err := w.oauth2Svc.VerifyCode(c, code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "授权码错误.",
		})
		return
	}

	// 通过微信的openid来查询用户是否存在，不存在就创建。
	user, err := w.userSvc.FindOrCreateByWechat(c, wechatInfo)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
		})
		return
	}

	userAgent := c.GetHeader("user-agent")
	err = w.jwtHandler.SetLoginToken(c, user.ID, userAgent)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 5})
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
	})
}
