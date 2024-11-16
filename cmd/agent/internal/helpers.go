package internal

import (
	"net/http"
	"encoding/json"
)

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
		json.NewEncoder(w).Encode(metrics)
	}
	if dataType == "gauge" {
		GaugeValue, _ = Data.GetGauge(name)
		GaugeValueFloat64 = float64(GaugeValue)
		var metrics = Metrics{
			ID: name,    
			MType: dataType,
	        Value: &GaugeValueFloat64,
		}
		json.NewEncoder(w).Encode(metrics)
	}
}