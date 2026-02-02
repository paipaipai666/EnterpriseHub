package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte(os.Getenv("SECRET_KEY"))

var publicPaths = []string{
	"/api/v1/auth/login",
	"/api/v1/users/register",
}

func JWTAuth(c *gin.Context) {
	path := c.Request.URL.Path

	if isPublicPath(path) {
		c.Next()
		return
	}

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"data":    "Authorization header missing",
		})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"data":    "Authorization header format must be Bearer {token}",
		})
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSampleSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.JSON(401, gin.H{
			"message": "failed",
			"data":    "Failed to parse token:" + err.Error(),
		})
		c.Abort()
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"error":   "Invalid token",
		})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"error":   "Invalid token claims",
		})
		c.Abort()
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		fmt.Println("Invalid or missing expiration time")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"data":    "Invalid expiration time",
		})
		c.Abort()
		return
	}

	if exp < float64(time.Now().Unix()) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"data":    "Token has expired",
		})
		c.Abort()
		return
	}

	// 修复：安全地获取用户名
	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "failed",
			"data":    "Invalid username in token",
		})
		c.Abort()
		return
	}

	c.Set("username", username)
	fmt.Println("JWTAuth: Token validated successfully for user:", username)
	c.Next()
}

func isPublicPath(path string) bool {
	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
