package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/cmilhench/x/exp/logger"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		client := r.RemoteAddr
		if c := strings.LastIndex(client, ":"); c != -1 {
			client = client[:c]
		}

		trans := "-"
		if value := r.Header.Get("X-Request-ID"); value != "" {
			trans = value
		}

		proxy := &wrappedResponseWriter{w, http.StatusOK, 0}

		defer func() {
			duration := time.Since(start)
			log := logger.Trace
			if proxy.status > 499 {
				log = logger.Error
			} else if proxy.status > 399 {
				log = logger.Debug
			}

			log("HTTP Request",
				"client", client,
				"request_id", trans,
				"method", r.Method,
				"uri", r.RequestURI,
				"proto", r.Proto,
				"status", proxy.status,
				"size", proxy.size,
				"duration_ms", duration.Milliseconds(),
			)
		}()
		h.ServeHTTP(proxy, r)
	})
}
