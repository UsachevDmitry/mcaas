package main

import (
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

    
    // Проверяем что данные не пустые
    if dataType == "" || name == "" || value == ""  {
        w.WriteHeader(http.StatusNotFound)  
        return
    }
    // if _, err := strconv.ParseInt(value, 10, 64); err != nil && dataType == "counter" {
    //     w.WriteHeader(http.StatusBadRequest)  
    //     return
    // }
    // if _, err := strconv.ParseFloat(value, 64); err != nil && dataType == "gauge" {
    //     w.WriteHeader(http.StatusBadRequest)  
    //     return
    // }
    // Проверяем type данных
    if dataType == "gauge" {
            value, err := strconv.ParseFloat(value, 64)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)  
                return
            } else {
                Data.UpdateGauge(name, gauge(value))
                w.WriteHeader(http.StatusOK)
                return
            }
        } else if dataType == "counter" {
            _, exists := Data.MetricsCounter[name]
            if !exists { 
                value, err := strconv.ParseInt(value, 10, 64)
                if err != nil {
                    w.WriteHeader(http.StatusBadRequest)  
                    return
                } else {
                    Data.UpdateCounter(name, counter(value))
                    //Data.MetricsCounter[name] = counter(f)
                    w.WriteHeader(http.StatusOK)
                    return
                }
            } else { 
                value, err := strconv.ParseInt(value, 10, 64)
                if err != nil {
                    w.WriteHeader(http.StatusBadRequest)  
                    return
                } else {
                    Data.AddCounter(name, counter(value))
                    //Data.MetricsCounter[name] += counter(value)
                    //fmt.Println("Exist! ", name, Data.MetricsCounter[name])
                    w.WriteHeader(http.StatusOK)
                    return
                }
            }     
    } else {
        w.WriteHeader(http.StatusBadRequest)
    }    
}