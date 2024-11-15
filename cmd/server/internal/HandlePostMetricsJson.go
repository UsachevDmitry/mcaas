package internal

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

func HandlePostMetricsJson() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string = "application/json"
		var dataType string
		var name string
		var value string
		var metrics Metrics


		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			message := Message{Message: "provided json file is invalid."}
			json.NewEncoder(w).Encode(message)
			return
		}
		
		dataType = metrics.MType
		name = metrics.ID
		if dataType == "counter" {
			if metrics.Delta == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				message := Message{Message: "counter value == nil."}
				json.NewEncoder(w).Encode(message)
				return
			} else {
				WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				value = strconv.Itoa(int(*metrics.Delta))
				//json.NewEncoder(w).Encode(metrics)
			}
		}
		if dataType == "gauge" {
			if metrics.Value == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				message := Message{Message: "gauge value == nil."}
				json.NewEncoder(w).Encode(message)
				return
			} else {
				WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				value = fmt.Sprintf("%f", *metrics.Value)
				//json.NewEncoder(w).Encode(metrics)
			}			
		}


		if dataType == "" || name == "" || value == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			message := Message{Message: "provided json file is invalid."}
			json.NewEncoder(w).Encode(message)
			return
		}

		if dataType == "gauge" {
			value, err := strconv.ParseFloat(value, 64)
			if err != nil {
				WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				return
			} else {
				Data.UpdateGauge(name, gauge(value))
				WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
				return
			}
		} else if dataType == "counter" {
			_, exists := Data.GetCounter(name)
			if !exists {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					Data.UpdateCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					return
				}
			} else {
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					Data.AddCounter(name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					return
				}
			}
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}

		CounterValue, _ := Data.GetCounter(name)
		GaugeValue, _ := Data.GetGauge(name)
		
		var CounterValueInt64 int64 = int64(CounterValue)
		var GaugeValueFloat64 float64 = float64(GaugeValue)

		var metrics2 = Metrics{
			ID: name,    
			MType: dataType,
			Delta: &CounterValueInt64,
			Value: &GaugeValueFloat64,
		}
		
		fmt.Println(metrics2)
	}
	return http.HandlerFunc(fn)
}

// type Metrics struct {
// 	ID    string   `json:"id"`              // имя метрики key
// 	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter 
// 	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter value
// 	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge value
//  } 