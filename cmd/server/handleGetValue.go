package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func handleGetValue(w http.ResponseWriter, r *http.Request) {
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
