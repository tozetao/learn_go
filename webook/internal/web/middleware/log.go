package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// 入口：指你的系统收到了请求，并返回了响应。
// 出口：即第三方调用，比如gorm对mysql的调用，redis的调用，其他接入的服务都统称为出口。

// 对入口进行日志记录

// AccessLog 要记录的日志数据
type AccessLog struct {
	ReqBody  string
	RespBody string
	Duration time.Duration
	Path     string
	Method   string
}

type LogMiddleware struct {
	logFn func(c *gin.Context, accessLog AccessLog)

	allowReqBody  bool
	allowRespBody bool
}

func NewLogMiddleware(logFn func(c *gin.Context, log AccessLog)) *LogMiddleware {
	return &LogMiddleware{
		logFn:         logFn,
		allowReqBody:  true,
		allowRespBody: true,
	}
}

// Build 构建一个gin中间件函数
func (m *LogMiddleware) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 记录的日志数据
		accessLog := AccessLog{}

		path := c.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}
		accessLog.Path = path
		accessLog.Method = c.Request.Method

		if m.allowReqBody {
			// Request.Body是一个stream，只能读取一次。读取完毕后需要重新赋值。
			body, _ := c.GetRawData()
			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			if len(body) > 2048 {
				body = body[0:2048]
			}
			accessLog.ReqBody = string(body)
		}

		if m.allowRespBody {
			c.Writer = responseWriter{
				accessLog:      &accessLog,
				ResponseWriter: c.Writer,
			}
		}

		defer func() {
			accessLog.Duration = time.Since(start)
			m.logFn(c, accessLog)
		}()

		c.Next()
	}
}

type responseWriter struct {
	accessLog *AccessLog
	gin.ResponseWriter
}

func (r *responseWriter) Writer(body []byte) (int, error) {
	r.accessLog.RespBody = string(body)
	return r.Writer(body)
}
