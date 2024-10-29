package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func handlePostMetrics(w http.ResponseWriter, r *http.Request) {
	var dataType string
	var name string
	var value string

	dataType = mux.Vars(r)["type"]
	name = mux.Vars(r)["name"]
	value = mux.Vars(r)["value"]

	if dataType == "" || name == "" || value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if dataType == "gauge" {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			Data.UpdateGauge(name, gauge(value))
			w.WriteHeader(http.StatusOK)
			return
		}
	} else if dataType == "counter" {
		_, exists := Data.GetCounter(name)
		if !exists {
			value, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				Data.UpdateCounter(name, counter(value))
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			value, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				Data.AddCounter(name, counter(value))
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
