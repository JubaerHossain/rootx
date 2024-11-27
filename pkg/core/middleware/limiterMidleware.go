// /Users/sookh.com/go/pkg/mod/github.com/!jubaer!hossain/rootx@v1.4.7/pkg/core/middleware/limiterMidleware.go
package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/JubaerHossain/rootx/pkg/core/limiter"
	"github.com/JubaerHossain/rootx/pkg/utils"
	"golang.org/x/time/rate"
)

// RateLimitConfig holds the configuration for rate limiting
type RateLimitConfig struct {
	Enabled      bool
	Limit        int
	Duration     time.Duration
	WhitelistIPs []string
	BlacklistIPs []string
}

var (
	whitelistedIPs = make(map[string]bool)
	blacklistedIPs = make(map[string]bool)
)

// isIPWhitelisted checks if an IP is in the whitelist
func isIPWhitelisted(ip string) bool {
	return whitelistedIPs[ip]
}

// isIPBlacklisted checks if an IP is in the blacklist
func isIPBlacklisted(ip string) bool {
	return blacklistedIPs[ip]
}

// getClientIP extracts the real client IP considering various headers
func getClientIP(r *http.Request) string {
	// Check CF-Connecting-IP (Cloudflare)
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	// Check X-Real-IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Check X-Forwarded-For
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// LimiterMiddleware creates a new rate limiting middleware
func LimiterMiddleware(next http.Handler) http.Handler {
	// Initialize configuration
	config := loadRateLimitConfig()
	if !config.Enabled {
		log.Println("Rate limiting is disabled.")
		return next
	}

	// Initialize IP lists
	initializeIPLists(config)

	var rateLimiter = limiter.NewIPRateLimiter(rate.Every(config.Duration), config.Limit)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		// Check whitelist
		if isIPWhitelisted(clientIP) {
			next.ServeHTTP(w, r)
			return
		}

		// Check blacklist
		if isIPBlacklisted(clientIP) {
			utils.WriteJSONError(w, http.StatusForbidden, "Access denied")
			return
		}

		// Get limiter for this IP
		limiter := rateLimiter.GetLimiter(clientIP)
		
		// Set standard rate limit headers
		setRateLimitHeaders(w, limiter, config)

		// Check if rate limit is exceeded
		if !limiter.Allow() {
			handleRateLimitExceeded(w, r, clientIP, config)
			return
		}

		// Add custom headers for debugging/monitoring
		w.Header().Set("X-Rate-Limit-IP", clientIP)
		w.Header().Set("X-Rate-Limit-Active-IPs", fmt.Sprintf("%d", rateLimiter.GetIPCount()))

		next.ServeHTTP(w, r)
	})
}

func loadRateLimitConfig() RateLimitConfig {
	duration, err := time.ParseDuration(os.Getenv("RATE_LIMIT_DURATION"))
	if err != nil {
		duration = time.Minute
	}

	limit := config.GlobalConfig.RateLimit
	if limit <= 0 {
		limit = 100
	}

	return RateLimitConfig{
		Enabled:      config.GlobalConfig.RateLimitEnabled,
		Limit:        limit,
		Duration:     duration,
		WhitelistIPs: strings.Split(os.Getenv("RATE_LIMIT_WHITELIST"), ","),
		BlacklistIPs: strings.Split(os.Getenv("RATE_LIMIT_BLACKLIST"), ","),
	}
}

func initializeIPLists(config RateLimitConfig) {
	// Initialize whitelist
	for _, ip := range config.WhitelistIPs {
		if ip = strings.TrimSpace(ip); ip != "" {
			whitelistedIPs[ip] = true
		}
	}

	// Initialize blacklist
	for _, ip := range config.BlacklistIPs {
		if ip = strings.TrimSpace(ip); ip != "" {
			blacklistedIPs[ip] = true
		}
	}
}

func setRateLimitHeaders(w http.ResponseWriter, limiter *rate.Limiter, config RateLimitConfig) {
	tokens := limiter.Tokens()
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.Limit))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", tokens))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(config.Duration).Unix()))
}

func handleRateLimitExceeded(w http.ResponseWriter, r *http.Request, clientIP string, config RateLimitConfig) {
	log.Printf("Rate limit exceeded for IP: %s, UA: %s, Path: %s", 
		clientIP, 
		r.UserAgent(),
		r.URL.Path,
	)

	w.Header().Set("Retry-After", fmt.Sprintf("%d", config.Duration/time.Second))
	utils.WriteJSONError(w, 
		http.StatusTooManyRequests, 
		fmt.Sprintf("Rate limit exceeded. Please try again in %d seconds", 
			config.Duration/time.Second,
		),
	)
}
