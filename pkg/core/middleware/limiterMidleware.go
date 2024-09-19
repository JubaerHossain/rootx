package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/config"
	"github.com/JubaerHossain/rootx/pkg/core/limiter"
	"github.com/JubaerHossain/rootx/pkg/utils"
	"golang.org/x/time/rate"
)

func LimiterMiddleware(next http.Handler) http.Handler {
	is_limit_enabled := config.GlobalConfig.RateLimitEnabled
	if !is_limit_enabled {
		return next
	}
	limit := config.GlobalConfig.RateLimit
	duration, err := time.ParseDuration(os.Getenv("RATE_LIMIT_DURATION"))
	if err != nil {
		duration = time.Second * 2
	}
	var limiter = limiter.NewIPRateLimiter(rate.Every(duration), limit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			utils.WriteJSONError(w, http.StatusTooManyRequests, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}
