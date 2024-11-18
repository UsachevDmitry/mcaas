package internal

import (
	"net/http"
	"encoding/json"
)

func HandlePostMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType = "application/json"
		var DataType string
		var Name string
		var metrics Metrics
		var ValueInt64 int64
		var ValueFloat64 float64
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			// message := Message{Message: "provided json file is invalid."}
			// json.NewEncoder(w).Encode(message)
			return
		}
		
		DataType = metrics.MType
		Name = metrics.ID

		if DataType == "" || Name == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			// message := Message{Message: "provided json file is invalid."}
			// json.NewEncoder(w).Encode(message)
			return
		}

		if DataType == "counter" {
			if metrics.Delta == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				// message := Message{Message: "counter value == nil."}
				// json.NewEncoder(w).Encode(message)
				return
			} else {
				ValueInt64 = int64(*metrics.Delta)
			}
			_, exists := Data.GetCounter(Name)
			if !exists {
				Data.UpdateCounter(Name, counter(ValueInt64))				
				PostMetricAnswer(Name, DataType, w)
				return
			} else {
				Data.AddCounter(Name, counter(ValueInt64))
				PostMetricAnswer(Name, DataType, w)
				return
				}
		} else if DataType == "gauge" {
			if metrics.Value == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				// message := Message{Message: "gauge value == nil."}
				// json.NewEncoder(w).Encode(message)
				return
			} else {
				ValueFloat64 = float64(*metrics.Value)
				Data.UpdateGauge(Name, gauge(ValueFloat64))
				PostMetricAnswer(Name, DataType, w)
				return
			}			
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}	
	return http.HandlerFunc(fn)
}
