package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics)
    log.Fatal(http.ListenAndServe(":8080", mux))
}
