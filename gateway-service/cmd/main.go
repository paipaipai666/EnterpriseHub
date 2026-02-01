package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/gateway-service/middleware"
)

func main() {
	r := gin.Default()

	// JWT 校验中间件
	r.Use(middleware.JWTAuth)

	// 路由转发
	r.Any("/api/users/*path", proxy("http://user-service:8000"))
	r.Any("/api/auth/*path", proxy("http://auth-service:9000"))

	r.Run(":8080")
}

func proxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
