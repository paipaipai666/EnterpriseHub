package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"go.uber.org/zap"
)

func GinLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	c.Next()

	end := time.Now()
	latency := end.Sub(start)

	if query != "" {
		path = path + "?" + query
	}

	initializers.Log.Info("HTTP Request",
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("ip", c.ClientIP()),
		zap.Duration("latency", latency),
		zap.String("user-agent", c.Request.UserAgent()),
	)
}
