package main

import (
	"flag"
	"github.com/UsachevDmitry/mcaas/cmd/server/internal"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

const defaultAddr = "localhost:8080"

var addr = flag.String("a", defaultAddr, "Адрес HTTP-сервера")

func main() {
	flag.Parse()
	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		*addr = addrEnv
	}
	router := mux.NewRouter()
	router.HandleFunc("/", internal.WithLoggingGet(internal.HandleIndex())).Methods(http.MethodGet)
	router.HandleFunc("/update/", internal.WithLoggingPost(internal.HandlePostMetricsJson())).Methods(http.MethodPost)
	router.HandleFunc("/update/{type}/{name}/{value}", internal.WithLoggingPost(internal.HandlePostMetrics())).Methods(http.MethodPost)
	//router.HandleFunc("/value/", internal.WithLoggingGet(internal.HandleGetMetricsJson())).Methods(http.MethodGet)
	router.HandleFunc("/value/{type}/{name}", internal.WithLoggingGet(internal.HandleGetValue())).Methods(http.MethodGet)

	internal.Logger()
	internal.GlobalSugar.Infow(
		"Starting server",
		"addr", *addr,
	)
	if err := http.ListenAndServe(*addr, router); err != nil {
		internal.GlobalSugar.Fatalw(err.Error(), "event", "start server")
	}
}
