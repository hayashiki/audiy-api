package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// logger
func middlewareLogger(logger Logger) func(next http.Handler) http.Handler {
	l, ok := logger.(*zap.SugaredLogger)
	if ok {
		log := l.Desugar()
		return func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

				t := time.Now()
				defer func() {
					log.Info("request",
						zap.String("proto", r.Proto),
						zap.String("path", r.URL.Path),
						zap.String("host", r.Host),
						zap.String("method", r.Method),
						zap.String("remote", r.RemoteAddr),
						zap.String("reqID", middleware.GetReqID(r.Context())),
						zap.Duration("duration", time.Since(t)),
						zap.Int("status", ww.Status()),
						zap.Int("size", ww.BytesWritten()),
						zap.String("user_agent", r.UserAgent()),
					)
				}()

				next.ServeHTTP(ww, r)
			}
			return http.HandlerFunc(fn)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
	}
}
