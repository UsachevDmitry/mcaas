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
    // Проверяем данные
    if name == "" {
        w.WriteHeader(http.StatusNotFound)  
        return
    }
    if _, err := strconv.ParseInt(value, 10, 64); err != nil && dataType == "counter" {
        w.WriteHeader(http.StatusBadRequest)  
        return
    }
    if _, err := strconv.ParseFloat(value, 64); err != nil && dataType == "gauge" {
        w.WriteHeader(http.StatusBadRequest)  
        return
    }
    // Проверяем type данных
    if dataType == "gauge" {
        if name == "" {
            w.WriteHeader(http.StatusNotFound)    
            return
        } else {                
            // f, _ := strconv.ParseFloat(value, 64)
            // Data.MetricsGauge[name] = gauge(f)
            value, err := strconv.ParseFloat(value, 64)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)  
                return
            } else {
                Data.AddGauge(name, gauge(value))
                w.WriteHeader(http.StatusOK)
                return
            }

        }
    } else if dataType == "counter" {
        _, exists := Data.MetricsCounter[name]
        if !exists { 
            f, _ := strconv.ParseInt(value, 10, 64)
            Data.MetricsCounter[name] = counter(f)
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