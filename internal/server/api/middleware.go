package api

import (
	"go-boilerplate/internal/lib"
	"net/http"
	"strings"

	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Correlation-ID")
		if requestID == "" {
			requestID = lib.GetNewID()
		}
		// Use the unexported type for the key to avoid collisions
		ctx := WithCorrelationID(r.Context(), requestID)
		// Set on request for downstream handlers that might not have direct context access (rare)
		r.Header.Set("X-Correlation-ID", requestID)
		// Also common to use X-Request-ID
		r.Header.Set("X-Request-ID", requestID)
		// Set on response
		w.Header().Set("X-Correlation-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SaveRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithRequest(r.Context(), r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default-src 'self' is a good starting point.
		// Unpkg is often used for Huma's OpenAPI UI, so it's included. Adjust as needed.
		// 'unsafe-inline' for styles is common for simple setups but try to avoid if possible.
		var cspDirectives = []string{
			"default-src 'self'",
			"script-src 'self' https://unpkg.com", // For Huma docs UI
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com https://unpkg.com", // For Huma docs UI & fonts
			"connect-src 'self' https://unpkg.com",                                            // For Huma docs UI
			"img-src 'self' data: https://unpkg.com",                                          // For Huma docs UI
			"font-src 'self' https://fonts.gstatic.com",
			"worker-src 'self' blob:", // Huma docs UI might use workers
			"form-action 'self'",
			"frame-ancestors 'none'",
		}
		w.Header().Set("Content-Security-Policy", strings.Join(cspDirectives, "; "))
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin") // A common, reasonable default
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")                                          // Prevent clickjacking
		w.Header().Set("X-XSS-Protection", "0")                                            // Modern browsers have better built-in XSS protection; CSP is preferred.
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains") // If serving over HTTPS

		next.ServeHTTP(w, r)
	})
}

func (wrapper *APIWrapper) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip        = r.RemoteAddr
			proto     = r.Proto
			method    = r.Method
			uri       = r.URL.RequestURI()
			requestID = GetCorrelationID(r.Context())
		)
		wrapper.logger.Info("received request",
			"ip", ip,
			"protocol", proto,
			"method", method,
			"uri", uri,
			"correlation_id", requestID,
		)
		next.ServeHTTP(w, r)
	})
}

func (w *APIWrapper) SetupMiddleware(mux *http.ServeMux) http.Handler {
	w.logger.Info("Setting up middleware...")

	allowedHeaders := []string{"Authorization", "Content-Type"}
	allowedHeaders = append(allowedHeaders, w.config.Server.AllowedHeaders...)

	allowedOrigins := w.config.Server.AllowedOrigins
	if w.config.Server.AllowedOrigins == nil {
		allowedOrigins = []string{}
	}
	corsOptions := cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   allowedHeaders,
		AllowedMethods:   []string{"DELETE", "GET", "OPTIONS", "POST", "PUT"},
		AllowedOrigins:   allowedOrigins,
	}
	w.logger.Info("CORS allowed origins", "origins", w.config.Server.AllowedOrigins)
	corsHandler := cors.New(corsOptions)

	middlewareChain := alice.New(
		corsHandler.Handler,
		CorrelationIDMiddleware,
		w.logRequest,
		CommonHeaders,
		SaveRequest,
	)

	return middlewareChain.Then(mux)
}
