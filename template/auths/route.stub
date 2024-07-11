package authHttp

import (
	"net/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
    "github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// AuthRouter registers routes for API endpoints
func AuthRouter(router *http.ServeMux, application *app.App) http.Handler {
	
	handler := NewHandler(application)
	// Register auth routes

	router.Handle("POST /auth/sign-in", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetSignIn)))
	router.Handle("POST /auth/refresh-token", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetRefreshToken)))
   

	return router
}
