package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn_go/webook/internal/web"
	"net/http"
)

// LoginJWTMiddlewareBuilder 该中间件用于检查短token
type LoginJWTMiddlewareBuilder struct {
	paths      []string
	jwtHandler *web.JWTHandler
}

func NewLoginJWTMiddlewareBuilder(handler *web.JWTHandler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		jwtHandler: handler,
	}
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

		tokenStr, err := m.jwtHandler.ExtractToken(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtHandler.ShortTokenKey, nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userAgent := c.GetHeader("User-Agent")
		if userAgent != claims.UserAgent {
			// TODO: 异常情况，需要记录日志。
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = m.jwtHandler.CheckSession(c, claims.SSid)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", claims)
	}
}
