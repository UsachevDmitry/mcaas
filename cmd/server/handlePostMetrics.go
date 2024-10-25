package main

import (
	// "io/ioutil"
	"net/http"
	// "encoding/json"
	"fmt"
    "strconv"
	"github.com/gorilla/mux"
    "strings"
	//"golang.org/x/tools/godoc/vfs/gatefs"
)

// type Metric struct {
//     Type string
//     Name string
//     Value float64
// }

// func handlePostMetrics(w http.ResponseWriter, r *http.Request) {
// // Получаем данные из запроса
// body, err := ioutil.ReadAll(r.Body)
// if err != nil {
// // Обработка ошибки
// w.WriteHeader(http.StatusInternalServerError)
// return
// }

// // Разбираем данные
// metric := Metric{}
// err = json.Unmarshal(body, &metric)
// if err != nil {
// // Обработка ошибки
// w.WriteHeader(http.StatusBadRequest)
// return
// }

// // Проверяем формат данных
// if metric.Type == "" || metric.Name == "" || metric.Value == 0 {
// // Обработка ошибки
// w.WriteHeader(http.StatusBadRequest)
// return
// }

// // Обработка метрики
// // Здесь можно добавить логику для сохранения метрики в базе данных или отправки в систему мониторинга

// // Ответ клиенту
// w.WriteHeader(http.StatusOK)
// }
func handlePostMetrics(w http.ResponseWriter, r *http.Request) {
    type gauge float64
    //type counter int64

    MemStorage := map[string]gauge {
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
    }
    var type2 string
    var value2 string
    
	// Получаем данные из запроса
	type2 = mux.Vars(r)["type"]
	//name2 := mux.Vars(r)["name"]
	value2 = mux.Vars(r)["value"]
    fmt.Println("type=",type2)
	type2 = strings.TrimLeft(strings.TrimRight(type2, "}"), "{")
    fmt.Println("type=",type2)
    fmt.Println("=======")
    // Проверяем type данных
	for key := range MemStorage {
        fmt.Println("key= ",key)
        if strings.Compare(type2, key) != 0 {
            //w.WriteHeader(http.StatusBadRequest)
            fmt.Println("bad")        
            //return
        } else {
            f, _ := strconv.ParseFloat(strings.TrimLeft(strings.TrimRight(value2, "}"), "{"), 64)
            MemStorage[key] = gauge(f)
            fmt.Println("GOOD! ", MemStorage[key])
            w.WriteHeader(http.StatusOK)
            return
        }
    }
    w.WriteHeader(http.StatusBadRequest)

	// Обработка метрики
	// Здесь можно добавить логику для сохранения метрики в базе данных или отправки в систему мониторинга

	// Ответ клиенту
	//w.WriteHeader(http.StatusOK)
    //fmt.Println("test")
}