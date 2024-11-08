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
	// fmt.Println("Адрес эндпоинта:", *addr)
	router := mux.NewRouter()
	router.HandleFunc("/update/{type}/{name}/{value}", internal.HandlePostMetrics).Methods("Post")
	router.HandleFunc("/", internal.WithLoggingHandleIndex(internal.HandleIndex())).Methods("Get")
	router.HandleFunc("/value/{type}/{name}", internal.WithLoggingHandleGetValue(internal.HandleGetValue())).Methods("Get")

	internal.Logger()
    // записываем в лог, что сервер запускается
    internal.GlobalSugar.Infow(
        "Starting server",
        "addr", *addr,
    )
    if err := http.ListenAndServe(*addr, router); err != nil {
        // записываем в лог ошибку, если сервер не запустился
        internal.GlobalSugar.Fatalw(err.Error(), "event", "start server")
    }
}
