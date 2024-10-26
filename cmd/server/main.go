package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
)

func main() {    
    router := mux.NewRouter()
    // Регистрация обработчика для метода POST
    router.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics).Methods("Post")
    // Запуск сервера
    log.Fatal(http.ListenAndServe(":8080", router))
}