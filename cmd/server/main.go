package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const defaultAddr = "localhost:8080"

var addr = flag.String("a", defaultAddr, "Адрес HTTP-сервера")

func main() {
	flag.Parse()
	fmt.Println("Адрес эндпоинта:", *addr)
	router := mux.NewRouter()
	// Регистрация обработчика для метода POST
	router.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics).Methods("Post")
	// Регистрация обработчиков для методов GET
	router.HandleFunc("/", handleIndex).Methods("Get")
	router.HandleFunc("/value/{type}/{name}", handleGetValue).Methods("Get")
	//Запуск сервера
	log.Fatal(http.ListenAndServe(*addr, router))
}
