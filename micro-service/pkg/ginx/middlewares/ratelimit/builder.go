package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	limter "learn_go/webook/pkg/limiter/ratelimit"
	"log"
	"net/http"
	"strings"
)

type Builder struct {
	limiter limter.Limiter
	prefix  string
}

// NewBuilder 可以考虑暴漏Limiter接口，作为参数提供。
func NewBuilder(limiter limter.Limiter, prefix string) *Builder {
	return &Builder{
		limiter: limiter,
		prefix:  prefix,
	}
}

// Build 返回中间件的实现
func (builder *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		var b strings.Builder

		// limiter:ip_
		b.WriteString(builder.prefix)
		b.WriteString(c.ClientIP())
		limited, err := builder.limiter.Limit(c, b.String())

		// 这边出错了要怎么办?
		// 依赖的服务报错了，需要记录日志并告警。
		if err != nil {
			fmt.Printf("limiter error: %v\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if limited {
			log.Println("too many request.")
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
