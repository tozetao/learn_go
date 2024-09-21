package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"learn_go/webook/internal/domain"
	"net/http"
	"net/url"
)

var redirectUrl = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type OAuth2Service interface {
	Auth2Url(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type OAuth2WechatService struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewOAuth2WechatService(appId string) OAuth2Service {
	return &OAuth2WechatService{
		appId:  appId,
		client: http.DefaultClient,
	}
}

// VerifyCode 使用code去换取token
func (o *OAuth2WechatService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		o.appId, o.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, "GET", accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// json解码
	var res Result
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信授权失败，err_code: %v, err_msg: %v", res.ErrCode, res.ErrMsg)
	}

	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

// Auth2Url 生成微信授权的回调地址
func (o *OAuth2WechatService) Auth2Url(c context.Context) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	state := shortuuid.New()
	return fmt.Sprintf(authURLPattern, o.appId, redirectUrl, state), nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
