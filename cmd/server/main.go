package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Регистрация обработчика для метода POST
	router.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics).Methods("Post")
	// Регистрация обработчиков для методов GET
	router.HandleFunc("/", handleIndex).Methods("Get")
	router.HandleFunc("/value/{type}/{name}", handleGetValue).Methods("Get")
	//Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", router))
}
