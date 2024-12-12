package middleware

import (
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
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
	}, labels)
	prometheus.MustRegister(summary)
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		duration := time.Since(now)
		summary.WithLabelValues(c.FullPath(), c.Request.Method, strconv.Itoa(c.Writer.Status())).
			Observe(float64(duration.Milliseconds()))
	}
}
