package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/gateway-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/gateway-service/middleware"
	"go.uber.org/zap"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("gateway")
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(middleware.JWTAuth)

	userServiceURL, err := url.Parse("http://localhost:8000/api/v1/users")
	if err != nil {
		initializers.Log.Fatal("can not parse user service url : ", zap.Error(err))
	}
	authServiceURL, err := url.Parse("http://localhost:9000/api/v1/auth")
	if err != nil {
		initializers.Log.Fatal("can not parse auth service url : ", zap.Error(err))
	}
	orderServiceURL, err := url.Parse("http://localhost:10000/api/v1/order")
	if err != nil {
		initializers.Log.Fatal("can not parse order service url : ", zap.Error(err))
	}

	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)

	router.Any("/api/v1/users/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/api/v1/auth/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		authProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/api/v1/order/*path", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		orderProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Run(":8080")
}
