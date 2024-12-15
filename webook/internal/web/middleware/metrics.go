package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type RequestSummary struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceID string
}

func (m *RequestSummary) Build() gin.HandlerFunc {
	labels := []string{"pattern", "method", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		Help:      "gin http请求统计",
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.95: 0.01,
			0.99: 0.005,
		},
	}, labels)
	prometheus.MustRegister(vector)

	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			// 准备上报 prometheus
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(pattern, method, strconv.Itoa(status)).
				Observe(float64(duration))
			fmt.Println("上报信息了")
		}()

		ctx.Next()

		//duration := time.Since(now)
		//// 以route-method-status记录该请求花费的时间（毫秒数）
		//summary.WithLabelValues(c.FullPath(), c.Request.Method, strconv.Itoa(c.Writer.Status())).
		//	Observe(float64(duration.Milliseconds()))
	}
}
