package main

import (
	"net/http"
	"fmt"
    "strconv"
	"github.com/gorilla/mux"
    "strings"
)

func handlePostMetrics(w http.ResponseWriter, r *http.Request) {
    type gauge float64
    //type counter int64
    type MemStorage struct {
        metrics map[string]gauge
    }
    data := MemStorage{
        metrics: map[string]gauge{
            "Alloc": 0,
            "BuckHashSys": 0,
            "Frees": 0,
            "GCCPUFraction": 0,
            "GCSys": 0,
            "HeapAlloc": 0,
            "HeapIdle": 0,
            "HeapInuse": 0,
            "HeapObjects": 0,
            "HeapReleased": 0,
            "HeapSys": 0,
            "LastGC": 0,
            "Lookups": 0,
            "MCacheInuse": 0,
            "MCacheSys": 0,
            "MSpanInuse": 0,
            "MSpanSys": 0,
            "Mallocs": 0,
            "NextGC": 0,
            "NumForcedGC": 0,
            "OtherSys": 0,
            "PauseTotalNs": 0,
            "StackInuse": 0,
            "StackSys": 0,
            "Sys": 0,
            "TotalAlloc": 0,
        },
    }

    var type2 string
    var value2 string
    
	// Получаем данные из запроса
	//type2 = mux.Vars(r)["type"]
	name2 := mux.Vars(r)["name"]
	value2 = mux.Vars(r)["value"]
    fmt.Println("name=",type2)
	name2 = strings.TrimLeft(strings.TrimRight(name2, "}"), "{")
    fmt.Println("name=",name2)
    fmt.Println("=======")
    // Проверяем type данных
    //gauge
    // Проверяем Name
	for key := range data.metrics {
        fmt.Println("key= ",key)
        if strings.Compare(name2, key) != 0 {
            //w.WriteHeader(http.StatusBadRequest)
            fmt.Println("bad")        
            //return
        } else {
            f, _ := strconv.ParseFloat(strings.TrimLeft(strings.TrimRight(value2, "}"), "{"), 64)
            data.metrics[key] = gauge(f)
            fmt.Println("GOOD! ", data.metrics[key])
            w.WriteHeader(http.StatusOK)
            return
        }
    }
    w.WriteHeader(http.StatusBadRequest)
}
