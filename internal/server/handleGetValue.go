package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleGetValue() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var dataType string
		var name string

		dataType = mux.Vars(r)["type"]
		name = mux.Vars(r)["name"]

		if dataType == "gauge" {
			value, exists := GetGauge(context.Background(), name)
			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%v", value)
			w.WriteHeader(http.StatusOK)
			return
		} else if dataType == "counter" {
			value, exists := GetCounter(context.Background(), name)
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
