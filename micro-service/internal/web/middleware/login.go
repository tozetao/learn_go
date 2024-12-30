package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (m *LoginMiddlewareBuilder) IgnorePath(paths ...string) *LoginMiddlewareBuilder {
	m.paths = append(m.paths, paths...)
	return m
}

func (m *LoginMiddlewareBuilder) Builder() gin.HandlerFunc {
	return func(c *gin.Context) {
		gob.Register(time.Now())

		path := c.Request.URL.Path
		for _, val := range m.paths {
			if path == val {
				return
			}
		}
		sess := sessions.Default(c)
		id := sess.Get("id")
		if id == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每10秒钟刷新session
		now := time.Now()
		updateTime := sess.Get("update_time")
		sess.Options(sessions.Options{
			MaxAge: 60,
			Path:   "/",
		})
		sess.Set("id", id)

		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		fmt.Printf("update_time: %v\n", updateTimeVal.Format(time.DateTime))

		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
