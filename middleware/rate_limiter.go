package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/pkg/response"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	clients   = make(map[string]*Client)
	mu        sync.Mutex
	limit     = rate.Every(time.Minute / 10)
	burstSize = 10
)

func cleanupClient() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, cl := range clients {
			if time.Since(cl.LastSeen) > time.Minute*5 {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func init() {
	go cleanupClient()
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, exist := clients[ip]; exist {
		client.LastSeen = time.Now()
		return client.Limiter
	}

	limiter := rate.NewLimiter(limit, burstSize)
	clients[ip] = &Client{limiter, time.Now()}
	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, "Too many request, please slow down")
			c.Abort()
			return
		}

		c.Next()
	}
}
