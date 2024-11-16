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
				PostMetricAnswer(name, dataType, w)
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
					PostMetricAnswer(name, dataType, w)
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
					PostMetricAnswer(name, dataType, w)
					return
				}
			}
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}
	return http.HandlerFunc(fn)
}
