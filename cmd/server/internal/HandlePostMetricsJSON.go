package internal

import (
	"net/http"
	"encoding/json"
	// "compress/gzip"
	// "io"
	// "fmt"
	// "bytes"
)

func HandlePostMetricsJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ContentType = "application/json"
		var DataType string
		var Name string
		var metrics Metrics
		var ValueInt64 int64
		var ValueFloat64 float64

		// // Чтение сжатых данных из тела запроса
		// body, err2 := io.ReadAll(r.Body)
		// if err2 != nil {
		// 	fmt.Println("Error reading request body:", err2)
		// 	return
		// }

		// // Распаковка данных
		// //var data map[string]interface{}
		// reader, err2 := gzip.NewReader(bytes.NewReader(body))
		// if err2 != nil {
		// 	fmt.Println("Error creating gzip reader:", err2)
		// 	return
		// }
		// defer reader.Close()
		
		decoder := json.NewDecoder(Decompress(r.Body))
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
