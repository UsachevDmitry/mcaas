package internal

import (
	"encoding/json"
	"net/http"
	"context"
	"fmt"
)

func HandlePostMetricsListJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string
		var DataType string
		var Name string
		var ValueInt64 int64
		var ValueFloat64 float64

		var metricsList []Metrics

		ContentType = r.Header.Get("Content-Type")

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metricsList)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			return
		}
		WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)

		for _, metrics := range metricsList {
			DataType = metrics.MType
			Name = metrics.ID

			if DataType == "" || Name == "" {
				continue
			}

			if DataType == "counter" {
				if metrics.Delta == nil {
					continue
				} else {
					ValueInt64 = int64(*metrics.Delta)
				}
				_, exists := GetCounter(context.TODO(), Name)
				if !exists {
					fmt.Println("!update")
					UpdateCounter(context.TODO(), Name, counter(ValueInt64))
					continue
				} else {
					AddCounter(context.TODO(), Name, counter(ValueInt64))
					continue
				}
			} else if DataType == "gauge" {
				if metrics.Value == nil {
					continue
				} else {
					ValueFloat64 = float64(*metrics.Value)
					UpdateGauge(context.TODO(), Name, gauge(ValueFloat64))
					continue
				}
			} else {
				continue
			}
		}
	}
	return http.HandlerFunc(fn)
}
