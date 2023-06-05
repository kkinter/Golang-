package http

import (
	"net/http"

	"go.uber.org/zap"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handled request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		next.ServeHTTP(w, r)
	})
}

// query d.Client.ExecContext(ctx, "SELECT pg_sleep(16)")
// func TimeoutMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
// 		defer cancel()

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
