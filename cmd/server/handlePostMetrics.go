package main

import (
    "fmt"
	"net/http"
	"strconv"
	"strings"
)

func handlePostMetrics(w http.ResponseWriter, r *http.Request) {
	var dataType string
	var name string
	var value string

	// Извлекаем URL
	u := strings.Split(r.URL.Path, "/")
	// Извлекаем значения из URL
	dataType = u[2]
	name = u[3]
	value = u[4]

    fmt.Println("+++++++++++++++")
	fmt.Println(r.URL.Path)
	fmt.Println(u)
    fmt.Println(dataType, name, value)
    fmt.Println("+++++++++++++++")
    
	// Проверяем что данные не пустые
	if dataType == "" || name == "" || value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
		_, exists := Data.GetGauge(name)
		if !exists {
			value, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				Data.UpdateCounter(name, counter(value))
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
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
