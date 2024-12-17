package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type Metrics struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceID string
}

func (m *Metrics) Build() gin.HandlerFunc {
	labels := []string{"pattern", "method", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      "http_resp_time",
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

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      "http_active_reqs",
		Help:      "活跃请求数",
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
	})
	prometheus.MustRegister(gauge)

	return func(ctx *gin.Context) {
		start := time.Now()
		gauge.Inc()

		defer func() {
			gauge.Dec()
			// 准备上报 prometheus
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(
				pattern,
				method,
				strconv.Itoa(status)).
				Observe(float64(duration))
		}()

		ctx.Next()

		//duration := time.Since(now)
		//// 以route-method-status记录该请求花费的时间（毫秒数）
		//summary.WithLabelValues(c.FullPath(), c.Request.Method, strconv.Itoa(c.Writer.Status())).
		//	Observe(float64(duration.Milliseconds()))
	}
}
