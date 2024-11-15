package internal
import (
	"net/http"
	"time"
)
var GlobalStatusCode int

func WithLoggingPost(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		GlobalSugar.Infoln(
			"statusCode", GlobalStatusCode,
			"size", r.Header.Get("Content-Length"),
		)
	}
	return logFn
}

func WithLoggingGet(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		GlobalSugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}
	return logFn
}

func WriteHeaderAndSaveStatus(statusCode int, ContentType string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", ContentType)
	GlobalStatusCode = statusCode
}

type Message struct {
	Message string `json:"message"`
}