package main

import (
	"flag"
	"fmt"
	"github.com/UsachevDmitry/mcaas/cmd/server/internal"
	"github.com/gorilla/mux"
	"log"
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
	fmt.Println("Адрес эндпоинта:", *addr)
	router := mux.NewRouter()
	router.HandleFunc("/update/{type}/{name}/{value}", internal.HandlePostMetrics).Methods("Post")
	router.HandleFunc("/", internal.HandleIndex).Methods("Get")
	router.HandleFunc("/value/{type}/{name}", internal.HandleGetValue).Methods("Get")
	log.Fatal(http.ListenAndServe(*addr, router))
}
