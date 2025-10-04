package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	service "product-service/Service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Claims struct {
	Username string `json:"sub"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("onHxm0yleYqLUcT6289RGI8tnAYxS9EY"), nil // Should match your JWT secret
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

func RateLimitMiddleware(ps *service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP
		// Check current count
		val, err := ps.Redis.Get(context.Background(), key).Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		count := 0
		if val != "" {
			count, _ = strconv.Atoi(val)
		}

		if count >= 100 { // 100 requests per minute
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Increment counter
		ps.Redis.Incr(context.Background(), key)
		ps.Redis.Expire(context.Background(), key, time.Minute)
		// Add headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", strconv.Itoa(100-count-1))

		c.Next()
	}
}
