package internal

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

func HandlePostMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType = "application/json"
		var DataType string
		var Name string
		var Value string
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
		
		DataType = metrics.MType
		Name = metrics.ID
		if DataType == "counter" {
			if metrics.Delta == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				message := Message{Message: "counter value == nil."}
				json.NewEncoder(w).Encode(message)
				return
			} else {
				//WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				Value = strconv.Itoa(int(*metrics.Delta))
			}
		}
		if DataType == "gauge" {
			if metrics.Value == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				message := Message{Message: "gauge value == nil."}
				json.NewEncoder(w).Encode(message)
				return
			} else {
				//WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				Value = fmt.Sprintf("%f", *metrics.Value)
			}			
		}


		if DataType == "" || Name == "" || Value == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			message := Message{Message: "provided json file is invalid."}
			json.NewEncoder(w).Encode(message)
			return
		}

		if DataType == "gauge" {
			// value, err := strconv.ParseFloat(Value, 64)
			value, err := ParseFloat10(Value, 10)
			if err != nil {
				WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				return
			} else {
				Data.UpdateGauge(Name, gauge(value))
				WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
				PostMetricAnswer(Name, DataType, w)
				return
			}
		} else if DataType == "counter" {
			_, exists := Data.GetCounter(Name)
			if !exists {
				value, err := strconv.ParseInt(Value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					Data.UpdateCounter(Name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					PostMetricAnswer(Name, DataType, w)
					return
				}
			} else {
				value, err := strconv.ParseInt(Value, 10, 64)
				if err != nil {
					WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					return
				} else {
					Data.AddCounter(Name, counter(value))
					WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					PostMetricAnswer(Name, DataType, w)
					return
				}
			}
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}
	return http.HandlerFunc(fn)
}
