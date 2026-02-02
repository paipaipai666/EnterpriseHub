package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	userServiceURL, _ := url.Parse("http://localhost:8000/api/v1/users")
	authServiceURL, _ := url.Parse("http://localhost:9000/api/v1/auth")

	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)

	router.Any("/api/v1/users/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/api/v1/auth/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		authProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Run(":8080")
}
