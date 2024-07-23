package log

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Write implements http.ResponseWriter.
func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

// WriteHeader implements http.ResponseWriter.
func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

// Handler returns a middleware that logs HTTP requests.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trans := "-"
		if value := r.Header.Get("X-Request-ID"); value != "" {
			trans = value
		}
		proxy := &wrappedResponseWriter{w, http.StatusOK, 0}
		start := time.Now()
		client := r.RemoteAddr
		if c := strings.LastIndex(client, ":"); c != -1 {
			client = client[:c]
		}
		line := fmt.Sprintf("%s %s - \"%s %s %s\"", client, trans, r.Method, r.RequestURI, r.Proto)
		defer func() {
			Debugf("%s %d %d %vms\n", line, proxy.status, proxy.size, time.Since(start).Milliseconds())
		}()
		h.ServeHTTP(proxy, r)
	})
}
