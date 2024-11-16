package internal
import (
	"net/http"
	"time"
	"encoding/json"
)

type Message struct {
	Message string `json:"message"`
}

var GlobalStatusCode int

func WithLoggingPost(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		GlobalSugar.Infoln(
			"statusCode", GlobalStatusCode,
			"size", r.Header.Get("Content-Length"),
		)
	}
	return logFn
}

func WithLoggingGet(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r)
		duration := time.Since(start)
		GlobalSugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
	}
	return logFn
}

func WriteHeaderAndSaveStatus(statusCode int, ContentType string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(statusCode)
	GlobalStatusCode = statusCode
}

func PostMetricAnswer(name string, dataType string, w http.ResponseWriter){
	var GaugeValue gauge
	var CounterValue counter
	var CounterValueInt64 int64
	var GaugeValueFloat64 float64

	if dataType == "counter" {
		CounterValue, _ = Data.GetCounter(name)
		CounterValueInt64 = int64(CounterValue)
		var metrics = Metrics{
			ID: name,    
			MType: dataType,
			Delta: &CounterValueInt64,
		}
		requestBody, err := json.Marshal(metrics)
		if err != nil {
			GlobalSugar.Errorln("Error marshaling JSON:", err)
			return
		}
		// w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(requestBody)
		//json.NewEncoder(w).Encode(metrics)
	}
	if dataType == "gauge" {
		GaugeValue, _ = Data.GetGauge(name)
		GaugeValueFloat64 = float64(GaugeValue)
		var metrics = Metrics{
			ID: name,    
			MType: dataType,
	        Value: &GaugeValueFloat64,
		}
		requestBody, err := json.Marshal(metrics)
		if err != nil {
			GlobalSugar.Errorln("Error marshaling JSON:", err)
			return
		}
		
		// w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(requestBody)
		//json.NewEncoder(w).Encode(metrics)
	}
	

	//ToDo почему omitempty не работает ? пришлось занести этот код в условия и убрать Delta или Value
	// var metrics = Metrics{
	// 	ID: name,    
	// 	MType: dataType,
	// 	Delta: &CounterValueInt64,
	// 	Value: &GaugeValueFloat64,
	// }
	// json.NewEncoder(w).Encode(metrics)
}