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
				//WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				//Value = strconv.Itoa(int(*metrics.Delta))

				ValueInt64 = int64(*metrics.Delta)
			}
			_, exists := Data.GetCounter(Name)
			if !exists {
					//value2, err := strconv.ParseInt(Value, 10, 64)
					//if err != nil {
					// WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
					// return
			//} else {
					Data.UpdateCounter(Name, counter(ValueInt64)) //value2
					//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					PostMetricAnswer(Name, DataType, w)
					return
			} else {
				// value2, err := strconv.ParseInt(Value, 10, 64)
				// if err != nil {
				// 	WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				// 	return
				// } else {
					Data.AddCounter(Name, counter(ValueInt64)) //value2
					//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
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
				//WriteHeaderAndSaveStatus(http.StatusCreated, ContentType, w)
				//Value = fmt.Sprintf("%10.15f", *metrics.Value)
				ValueFloat64 = float64(*metrics.Value)
				//value2, err := strconv.ParseFloat(Value, 64)			
				// if err != nil {
				// 	WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
				// 	return
				// } else {
					Data.UpdateGauge(Name, gauge(ValueFloat64)) //value2
					//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
					PostMetricAnswer(Name, DataType, w)
					return
				// }
			}			
		} else {
			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		}
	}	
	return http.HandlerFunc(fn)
}



		// if DataType == "gauge" {
		// 	value2, err := strconv.ParseFloat(Value, 64)			
		// 	if err != nil {
		// 		WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		// 		return
		// 	} else {
		// 		Data.UpdateGauge(Name, gauge(value2)) //value2
		// 		//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
		// 		PostMetricAnswer(Name, DataType, w)
		// 		return
		// 	}
		// } else if DataType == "counter" {
		// 	_, exists := Data.GetCounter(Name)
		// 	if !exists {
		// 			value2, err := strconv.ParseInt(Value, 10, 64)
		// 			if err != nil {
		// 			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		// 			return
		// 		} else {
		// 			Data.UpdateCounter(Name, counter(value2)) //value2
		// 			//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
		// 			PostMetricAnswer(Name, DataType, w)
		// 			return
		// 		}
		// 	} else {
		// 		value2, err := strconv.ParseInt(Value, 10, 64)
		// 		if err != nil {
		// 			WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		// 			return
		// 		} else {
		// 			Data.AddCounter(Name, counter(value2)) //value2
		// 			//WriteHeaderAndSaveStatus(http.StatusOK, ContentType, w)
		// 			PostMetricAnswer(Name, DataType, w)
		// 			return
		// 		}
		// 	}
		// } else {
		// 	WriteHeaderAndSaveStatus(http.StatusBadRequest, ContentType, w)
		// }