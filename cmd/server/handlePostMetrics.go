package main

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func handlePostMetrics(w http.ResponseWriter, r *http.Request) {   
    var dataType string
    var name string
    var value string
        
	// Получаем данные из запроса
	dataType = mux.Vars(r)["type"]
	name = mux.Vars(r)["name"]
	value = mux.Vars(r)["value"]
    if name == "yes" {
        w.WriteHeader(http.StatusNotFound)  
        return
    }
    // Проверяем type данных
    if dataType == "gauge" {
        // fmt.Println("key= ",key)
        //fmt.Println("name empty=",name)
        if name == "" {
            w.WriteHeader(http.StatusNotFound)
            //fmt.Println("bad")        
            return
        } else {                
            f, _ := strconv.ParseFloat(value, 64)
            Data.MetricsGauge[name] = gauge(f)
            fmt.Println("GOOD! ",name, Data.MetricsGauge[name])
            w.WriteHeader(http.StatusOK)
            return
        }
        //w.WriteHeader(http.StatusBadRequest)
    } else if dataType == "counter" {
        _, exists := Data.MetricsCounter[name]
        if !exists { 
            //fmt.Println("Not exist ", name, Data.Metrics[name])
            f, _ := strconv.ParseInt(value, 10, 64)
            Data.MetricsCounter[name] = counter(f)
            //w.WriteHeader(http.StatusNotFound)
            w.WriteHeader(http.StatusOK)
            return
        } else { 
            f, _ := strconv.ParseInt(value, 10, 64)
            Data.MetricsCounter[name] += counter(f)
            fmt.Println("Exist! ", name, Data.MetricsCounter[name])
            w.WriteHeader(http.StatusOK)
            return
        } 
    } else {
        w.WriteHeader(http.StatusBadRequest)
    }    
}