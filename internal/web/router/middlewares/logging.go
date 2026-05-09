package middlewares

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

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	size, err := rw.ResponseWriter.Write(b)

	rw.size += size

	return size, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"%d %s %s %s %dB",
			rw.status,
			r.Method,
			r.URL.Path,
			duration,
			rw.size,
		)
	})
}
