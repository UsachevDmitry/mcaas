package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func HandleGetValue() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var dataType string
		var name string

		dataType = mux.Vars(r)["type"]
		name = mux.Vars(r)["name"]

		if dataType == "gauge" {
			value, exists := Data.GetGauge(name)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%v", value)
			w.WriteHeader(http.StatusOK)
			return
		} else if dataType == "counter" {
			value, exists := Data.GetCounter(name)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%v", value)
			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	return http.HandlerFunc(fn)
}

func WithLoggingHandleGetValue(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
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
