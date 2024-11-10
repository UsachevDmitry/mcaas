package internal

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var StatusCode int
var Size string

func WriteHeaderAndSaveStatus(statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	StatusCode = statusCode
}

func HandlePostMetrics() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var dataType string
		var name string
		var value string

		dataType = mux.Vars(r)["type"]
		name = mux.Vars(r)["name"]
		value = mux.Vars(r)["value"]

		if dataType == "" || name == "" || value == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, w)
			return
		}

		if dataType == "gauge" {
			value, err := strconv.ParseFloat(value, 64)
			if err != nil {
				WriteHeaderAndSaveStatus(http.StatusBadRequest, w)
				return
			} else {
				Data.UpdateGauge(name, gauge(value))
				WriteHeaderAndSaveStatus(http.StatusOK, w)
				return
			}
		} else if dataType == "counter" {
			_, exists := Data.GetCounter(name)
			if !exists {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, w)
					return
				} else {
					Data.UpdateCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, w)
					return
				}
			} else {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, w)
					return
				} else {
					Data.AddCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, w)
					return
				}
			}
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, w)
		}
	}
	return http.HandlerFunc(fn)
}

func WithLoggingHandlePostMetrics(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		Size = r.Header.Get("Content-Length")
		GlobalSugar.Infoln(
			"statusCode", StatusCode,
			"size", Size,
		)
	}
	return logFn
}