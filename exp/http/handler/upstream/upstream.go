package upstream

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/cmilhench/x/exp/logger"
)

type Location struct {
	Pattern string
	Path    string
	Host    string
	Headers map[string]string
	Handler func(http.Handler) http.Handler
}

var _ http.Handler = (*handler)(nil)

type handler struct {
	locations []Location
}

func New(locations []Location) http.Handler {
	return &handler{
		locations: locations,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, location := range h.locations {
		if regexp.MustCompile(location.Pattern).MatchString(path) {
			h.upstream(w, r, location)
			return
		}
	}
	http.NotFound(w, r)
}

func (h *handler) upstream(w http.ResponseWriter, r *http.Request, location Location) {
	start := time.Now()
	pattern := regexp.MustCompile(location.Pattern)
	path := r.URL.Path
	host := r.URL.Host
	if len(location.Host) > 0 {
		host = pattern.ReplaceAllString(path, location.Host)
	}
	if len(location.Path) > 0 {
		path = pattern.ReplaceAllString(path, location.Path)
	}

	target, err := url.Parse(host)
	if err != nil {
		logger.Error("failed to parse upstream", "error", err, "host", host)
		http.Error(w, "Invalid upstream host", http.StatusInternalServerError)
		return
	}

	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.URL.Path = path

	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Header.Set("X-Real-IP", r.RemoteAddr)
	r.Host = target.Host

	for key, value := range location.Headers {
		r.Header.Set(key, value)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error("failed to upstream request", "error", err, "url", r.URL.String())
		w.WriteHeader(http.StatusBadGateway)
	}
	proxy.ModifyResponse = func(r *http.Response) error {
		logResponse(r, start, target, path)
		r.Header.Del("Server")
		r.Header.Del("X-Powered-By")
		//r.Header.Del("Set-Cookie")
		r.Header.Del("Request-Context")
		r.Header.Del("Strict-Transport-Security")
		r.Header.Del("X-Miniprofiler-Ids")
		if r.StatusCode > 499 {
			r.StatusCode = http.StatusBadGateway
			r.Status = http.StatusText(http.StatusBadGateway)
		}
		if r.StatusCode == http.StatusForbidden {
			r.StatusCode = http.StatusNotFound
			r.Status = http.StatusText(http.StatusNotFound)
			payload := []byte("404 page not found")
			r.Body = io.NopCloser(bytes.NewBuffer(payload))
			r.ContentLength = int64(len(payload))
			r.Header.Set("Content-Length", fmt.Sprint(len(payload)))
			r.Header.Set("Content-Type", "text/plain")
		}
		return nil
	}
	if location.Handler != nil {
		location.Handler(proxy).ServeHTTP(w, r)
	} else {
		proxy.ServeHTTP(w, r)
	}
}

func logResponse(r *http.Response, start time.Time, target *url.URL, path string) {
	duration := time.Since(start)
	log := logger.Trace
	if r.StatusCode > 499 {
		log = logger.Error
	} else if r.StatusCode > 399 {
		log = logger.Debug
	}
	client := r.Request.RemoteAddr
	if c := strings.LastIndex(client, ":"); c != -1 {
		client = client[:c]
	}
	trans := "-"
	if value := r.Request.Header.Get("X-Request-ID"); value != "" {
		trans = value
	}
	log("Upstream Request",
		"client", client,
		"request_id", trans,
		"method", r.Request.Method,
		"uri", target.String()+path, //r.Request.RequestURI,
		"proto", r.Proto,
		"status", r.StatusCode,
		"size", r.ContentLength,
		"duration_ms", duration.Milliseconds(),
	)
}
