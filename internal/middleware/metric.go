package middleware

import (
	"github.com/biryanim/wb_tech_L0/internal/metric"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func MetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		metric.IncRequestCounter()

		timeStart := time.Now()

		c.Next()

		difftime := time.Since(timeStart)
		status := strconv.Itoa(c.Writer.Status())

		if c.Writer.Status() >= http.StatusBadRequest {
			metric.IncResponseCounter("error", status, c.Request.Method, c.Request.RequestURI)
			metric.HistogramResponseTimeObserve("error", difftime.Seconds())
		} else {
			metric.IncResponseCounter("success", status, c.Request.Method, c.Request.RequestURI)
			metric.HistogramResponseTimeObserve("success", difftime.Seconds())
		}
	}
}
