package internal

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandlePostMetricsListJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType string
		var HashSHA256Value string
		var DataType string
		var Name string
		var ValueInt64 int64
		var ValueFloat64 float64
		var metricsList []Metrics

		ContentType = r.Header.Get("Content-Type")
		HashSHA256Value = r.Header.Get("HashSHA256")

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&metricsList)
		if err != nil {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
			return
		}

		if *Key != "" {
			metricsListByte, err := json.Marshal(&metricsList)
			if err != nil {
				GlobalSugar.Errorf("Error marshaling JSON:", err)
				return
			}
			keyAndData := make([]byte, len(metricsListByte)+len([]byte(*Key)))
			copy(keyAndData, metricsListByte)
			copy(keyAndData, []byte(*Key))
			hash := sha256.Sum256(keyAndData)
			hashString := fmt.Sprintf("%x", hash)
			w.Header().Set("HashSHA256", hashString)
			if hashString != HashSHA256Value {
				WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				return
			}
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
				_, exists := GetCounter(r.Context(), Name)
				if !exists {
					fmt.Println("!update")
					UpdateCounter(r.Context(), Name, counter(ValueInt64))
					continue
				} else {
					AddCounter(r.Context(), Name, counter(ValueInt64))
					continue
				}
			} else if DataType == "gauge" {
				if metrics.Value == nil {
					continue
				} else {
					ValueFloat64 = float64(*metrics.Value)
					UpdateGauge(r.Context(), Name, gauge(ValueFloat64))
					continue
				}
			} else {
				continue
			}
		}
	}
	return http.HandlerFunc(fn)
}
