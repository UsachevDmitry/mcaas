package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
)

func PostMetricAnswer(name string, dataType string, w http.ResponseWriter) {
	var GaugeValue gauge
	var CounterValue counter
	var CounterValueInt64 int64
	var GaugeValueFloat64 float64

	if dataType == "counter" {
		CounterValue, _ = Data.GetCounter(name)
		CounterValueInt64 = int64(CounterValue)
		var metrics = Metrics{
			ID:    name,
			MType: dataType,
			Delta: &CounterValueInt64,
		}
		json.NewEncoder(w).Encode(metrics)
	}
	if dataType == "gauge" {
		GaugeValue, _ = Data.GetGauge(name)
		GaugeValueFloat64 = float64(GaugeValue)
		var metrics = Metrics{
			ID:    name,
			MType: dataType,
			Value: &GaugeValueFloat64,
		}
		json.NewEncoder(w).Encode(metrics)
	}
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// Создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w := gzip.NewWriter(&b)
	// Запись данных
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// Обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// Переменная b содержит сжатые данные
	return b.Bytes(), nil
}
