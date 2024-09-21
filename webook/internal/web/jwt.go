package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn_go/webook/internal/domain"
	"time"
)

type jwtHandler struct{}

func (j *jwtHandler) setJWTToken(c *gin.Context, user domain.User, userAgent string) error {
	signingKey := []byte("ihbwtj3dGZvDKmgE2gyQL8gBU2saIE")
	claims := UserClaims{
		Uid:       user.ID,
		UserAgent: userAgent,
		// token过期时间, 1分钟有效期。
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (j *jwtHandler) setRefreshToken(c *gin.Context, user domain.User, userAgent string) error {
	signingKey := []byte("Ahbwtj74udvDKmgE2gyQL8gBU2saxF")
	claims := RefreshClaims{
		UserClaims{
			Uid:       user.ID,
			UserAgent: userAgent,
			// refresh token过期时间, 7天有效期。
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

type RefreshClaims struct {
	UserClaims
}

// 改造点:
// 1. 短token增加过期时间，但是中间件取消过期时间的验证。
// 2. 生成长短token
// 3. 增加一个刷新短token的接口。UserHandler:RefreshToken

/*
problems:
1. 长token过期了怎么办?
2. 短token过期了。
	前端通过长token换取短token。
3.
*/
