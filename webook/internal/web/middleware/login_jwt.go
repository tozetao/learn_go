package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn_go/webook/internal/web"
	"net/http"
	"strings"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (m *LoginJWTMiddlewareBuilder) IgnorePath(paths ...string) *LoginJWTMiddlewareBuilder {
	m.paths = append(m.paths, paths...)
	return m
}

func (m *LoginJWTMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, val := range m.paths {
			if path == val {
				return
			}
		}

		signingKey := []byte("ihbwtj3dGZvDKmgE2gyQL8gBU2saIE")
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		/*
			如果token被窃取了, 该怎么办?
			可以让前端携带用户登陆时的环境的特征，因为用户登录时的环境一般都是固定的，所以可以通过验证jwt token携带的用户特征，判断jwt token是否异常。
			为了方便，我们使用浏览器的user-agent来作为登录特征。
		*/
		userAgent := c.GetHeader("User-Agent")
		if userAgent != claims.UserAgent {
			// TODO: 异常情况，需要记录日志。
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("claims", claims)
		// fmt.Printf("parse successfully, UID: %v, user-agent: %s\n", claims.Uid, claims.UserAgent)
	}
}
