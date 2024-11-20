package internal

import (
	"net/http"
	"encoding/json"
)

func HandlePostMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string
		var DataType string
		var Name string
		var metrics Metrics
		var ValueInt64 int64
		var ValueFloat64 float64

		ContentType = r.Header.Get("Content-Type")

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metrics)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			return
		}
		
		DataType = metrics.MType
		Name = metrics.ID

		if DataType == "" || Name == "" {
			WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
			return
		}

		if DataType == "counter" {
			if metrics.Delta == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				return
			} else {
				ValueInt64 = int64(*metrics.Delta)
			}
			_, exists := Data.GetCounter(Name)
			if !exists {
				Data.UpdateCounter(Name, counter(ValueInt64))				
				PostMetricAnswer(Name, DataType, w, r)
				return
			} else {
				Data.AddCounter(Name, counter(ValueInt64))
				PostMetricAnswer(Name, DataType, w, r)
				return
				}
		} else if DataType == "gauge" {
			if metrics.Value == nil {
				WriteHeaderAndSaveStatus(http.StatusNotFound, ContentType, w)
				return
			} else {
				ValueFloat64 = float64(*metrics.Value)
				Data.UpdateGauge(Name, gauge(ValueFloat64))
				PostMetricAnswer(Name, DataType, w, r)
				return
			}			
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}	
	return http.HandlerFunc(fn)
}
