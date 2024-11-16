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
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", ContentType)
	GlobalStatusCode = statusCode
}

func PostMetricAnswer(name string, dataType string, w http.ResponseWriter){
	CounterValue, _ := Data.GetCounter(name)
	GaugeValue, _ := Data.GetGauge(name)
	
	var CounterValueInt64 int64 = int64(CounterValue)
	var GaugeValueFloat64 float64 = float64(GaugeValue)

	//ToDo разобраться с null и 0 при возврашение не заполненного значяения
	var metrics = Metrics{
		ID: name,    
		MType: dataType,
		Delta: &CounterValueInt64,
		Value: &GaugeValueFloat64,
	}
	json.NewEncoder(w).Encode(metrics)
}