package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
)

func main() {
    // Инициализация маршрутизатора
    //mux := http.NewServeMux()
    router := mux.NewRouter()

    // Регистрация обработчика для метода POST
    router.HandleFunc("/update/{type}/{name}/{value}", handlePostMetrics).Methods("Post")
    //?type=cpu&name=usage&value=100
    // Запуск сервера
    log.Fatal(http.ListenAndServe(":8080", router))

    // mux := http.NewServeMux()
    // mux.handlePostMetrics(`/`, handler)
    // type gauge float64
    // type counter int64
    // MemStorage := map[string]gauge{
    //     "Alloc": 0,
    //     "BuckHashSys": 0,
    //     "Frees": 0,
    //     "GCCPUFraction": 0,
    //     "GCSys": 0,
    //     "HeapAlloc": 0,
    //     "HeapIdle": 0,
    //     "HeapInuse": 0,
    //     "HeapObjects": 0,
    //     "HeapReleased": 0,
    //     "HeapSys": 0,
    //     "LastGC": 0,
    //     "Lookups": 0,
    //     "MCacheInuse": 0,
    //     "MCacheSys": 0,
    //     "MSpanInuse": 0,
    //     "MSpanSys": 0,
    //     "Mallocs": 0,
    //     "NextGC": 0,
    //     "NumForcedGC": 0,
    //     "OtherSys": 0,
    //     "PauseTotalNs": 0,
    //     "StackInuse": 0,
    //     "StackSys": 0,
    //     "Sys": 0,
    //     "TotalAlloc": 0,
    // }

    // err := http.ListenAndServe(`:8080`, mux)
    // if err != nil {
    //     panic(err)
    // }
} 