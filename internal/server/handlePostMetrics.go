package internal

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var Size string

func HandlePostMetrics() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string
		var dataType string
		var name string
		var value string

		ContentType = r.Header.Get("Content-Type")

		dataType = mux.Vars(r)["type"]
		name = mux.Vars(r)["name"]
		value = mux.Vars(r)["value"]

		if dataType == "" || name == "" || value == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			return
		}

		if dataType == "gauge" {
			value, err := strconv.ParseFloat(value, 64)
			if err != nil {
				WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				return
			} else {
				UpdateGauge(name, gauge(value))
				WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
				return
			}
		} else if dataType == "counter" {
			_, exists := GetCounter(name)
			if !exists {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					UpdateCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					return
				}
			} else {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					AddCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					return
				}
			}
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}
	return http.HandlerFunc(fn)
}
