package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors   = make(map[string]*rate.Limiter)
	mu         sync.Mutex
	rateLimit  = rate.Every(12 * time.Second) // 5 requests per minute (60/5 = 12 sec)
	burstLimit = 1
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burstLimit)
		visitors[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware applies the rate limiter to each request by IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please slow down",
			})
			return
		}

		c.Next()
	}
}
