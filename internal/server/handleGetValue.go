package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleGetValue() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var dataType string
		var name string

		dataType = mux.Vars(r)["type"]
		name = mux.Vars(r)["name"]

		if dataType == "gauge" {
			value, exists := GetGauge(name)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%v", value)
			w.WriteHeader(http.StatusOK)
			return
		} else if dataType == "counter" {
			value, exists := GetCounter(name)
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
