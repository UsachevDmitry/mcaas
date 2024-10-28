package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func handleGetValue(w http.ResponseWriter, r *http.Request) {
	var dataType string
	var name string

	dataType = mux.Vars(r)["type"]
	name = mux.Vars(r)["name"]
	// Проверяем type данных
	if dataType == "gauge" {
		value, exists := Data.GetGauge(name)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		//fmt.Fprintf(w, "<html><body><h1>Gauge</h1><p>name: %s<br>value: %v</p></body></html>", name, value)
		fmt.Fprintf(w, "%v", value)
		w.WriteHeader(http.StatusOK)
		return
	} else if dataType == "counter" {
		value, exists := Data.GetCounter(name)
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		//fmt.Fprintf(w, "<html><body><h1>Counter</h1><p>name: %s<br>value: %v</p></body></html>", name, value)
		fmt.Fprintf(w, "%v", value)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
