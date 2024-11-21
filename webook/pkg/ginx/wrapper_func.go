package ginx

import (
	"github.com/gin-gonic/gin"
	"learn_go/webook/pkg/logger"
	"net/http"
)

/*
	handler层模块化代码的抽象。
*/

var L logger.LoggerV2 = logger.NewNopLogger()

func WrapBody[Request any](fn func(req Request) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if c.Bind(&req) != nil {
			L.Info("binding error during article publication.")
			return
		}

		res, err := fn(req)
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		c.JSON(200, res)
	}
}

func WrapBodyAndClaims[Request any, Claims any](fn func(c *gin.Context, req Request, claims Claims) (Result, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		if c.Bind(&req) != nil {
			L.Info("binding error during article publication.")
			return
		}

		claimsVal := c.MustGet("user")
		userClaims, ok := claimsVal.(Claims)
		if !ok {
			L.Warn("get user claims error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		res, err := fn(c, req, userClaims)
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		c.JSON(200, res)
	}
}
