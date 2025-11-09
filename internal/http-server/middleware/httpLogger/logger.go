package httpLogger

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

func New(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger = logger.With(
			slog.String("component", "middleware/httpLogger"),
		)

		logger.Info("httpLogger middleware enabled")

		handler := func(w http.ResponseWriter, r *http.Request) {
			rDesc := logger.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			wrapper := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				rDesc.Info("request completed",
					slog.Int("status", wrapper.Status()),
					slog.Int("bytes", wrapper.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			next.ServeHTTP(wrapper, r)

		}

		return http.HandlerFunc(handler)
	}
}
