package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Logger(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         200,
			}

			logger.Printf("INFO: --> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Printf("INFO: <-- %s %s %d %d bytes in %v", r.Method, r.URL.Path, rw.status, rw.size, duration)
		})
	}
}
