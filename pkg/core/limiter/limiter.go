package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter handles rate limiting by IP address
type IPRateLimiter struct {
	ips    map[string]*limiterEntry
	mu     *sync.RWMutex
	r      rate.Limit
	b      int
	maxAge time.Duration
}

// NewIPRateLimiter creates a new IP rate limiter with cleanup routine
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:    make(map[string]*limiterEntry),
		mu:     &sync.RWMutex{},
		r:      r,
		b:      b,
		maxAge: time.Hour, // Cleanup entries older than 1 hour
	}

	// Start cleanup routine
	go i.cleanupRoutine()

	return i
}

// AddIP creates a new rate limiter and adds it to the ips map
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	entry := &limiterEntry{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	i.ips[ip] = entry

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	entry, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	// Update last seen time
	i.mu.Lock()
	entry.lastSeen = time.Now()
	i.mu.Unlock()

	return entry.limiter
}

// cleanupRoutine removes old entries periodically
func (i *IPRateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.cleanup()
	}
}

// cleanup removes entries that haven't been used recently
func (i *IPRateLimiter) cleanup() {
	i.mu.Lock()
	defer i.mu.Unlock()

	for ip, entry := range i.ips {
		if time.Since(entry.lastSeen) > i.maxAge {
			delete(i.ips, ip)
		}
	}
}

// GetIPCount returns the current number of IP limiters
func (i *IPRateLimiter) GetIPCount() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.ips)
}