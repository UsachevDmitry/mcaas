package main

import (
	"flag"
	"fmt"
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
	router.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics).Methods("Post")
	router.HandleFunc("/", handleIndex).Methods("Get")
	router.HandleFunc("/value/{type}/{name}", handleGetValue).Methods("Get")
	log.Fatal(http.ListenAndServe(*addr, router))
}
