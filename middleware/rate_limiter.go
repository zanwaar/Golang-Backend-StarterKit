package middleware

import (
	"net/http"
	"sync"

	"golang-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips sync.Map
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

type UserRateLimiter struct {
	users sync.Map
	mu    sync.Mutex
	r     rate.Limit
	b     int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		r: r,
		b: b,
	}
}

func NewUserRateLimiter(r rate.Limit, b int) *UserRateLimiter {
	return &UserRateLimiter{
		r: r,
		b: b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips.Load(ip)
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips.Store(ip, limiter)
	}

	return limiter.(*rate.Limiter)
}

func (u *UserRateLimiter) GetLimiter(userID uint) *rate.Limiter {
	u.mu.Lock()
	defer u.mu.Unlock()

	limiter, exists := u.users.Load(userID)
	if !exists {
		limiter = rate.NewLimiter(u.r, u.b)
		u.users.Store(userID, limiter)
	}

	return limiter.(*rate.Limiter)
}

func RateLimiterMiddleware() gin.HandlerFunc {
	// IP Limit: 5 requests per second, burst: 10
	ipLimiter := NewIPRateLimiter(5, 10)
	// User Limit: 10 requests per second, burst: 15
	userLimiter := NewUserRateLimiter(10, 15)

	return func(c *gin.Context) {
		// Check IP Limit
		ip := c.ClientIP()
		if !ipLimiter.GetLimiter(ip).Allow() {
			utils.ErrorResponse(c, "Too Many Requests", http.StatusTooManyRequests, "IP rate limit exceeded")
			c.Abort()
			return
		}

		// Check User Limit if authenticated
		userID, exists := c.Get("user_id")
		if exists {
			if !userLimiter.GetLimiter(userID.(uint)).Allow() {
				utils.ErrorResponse(c, "Too Many Requests", http.StatusTooManyRequests, "User rate limit exceeded")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
