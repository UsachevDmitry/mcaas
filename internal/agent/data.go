package internal

import (
	"sync"
)

type gauge float64
type counter int64
type MemStorage struct {
	MetricsGauge   map[string]gauge
	MetricsCounter map[string]counter
	Mutex          sync.RWMutex
}

var Data = &MemStorage{
	MetricsGauge:   map[string]gauge{},
	MetricsCounter: map[string]counter{},
	Mutex:          sync.RWMutex{},
}

type MemStorageInterface interface {
	UpdateGauge(key string, value gauge)
	UpdateCounter(key string, value counter)
	AddCounter(key string, value counter)
	GetGauge(key string) (gauge, bool)
	GetCounter(key string) (counter, bool)
	DeleteGauge(key string)
	DeleteCounter(key string)
}

func (ms *MemStorage) UpdateGauge(key string, value gauge) {
	ms.Mutex.Lock()
	ms.MetricsGauge[key] = value
	ms.Mutex.Unlock()
}

func (ms *MemStorage) UpdateCounter(key string, value counter) {
	ms.Mutex.Lock()
	ms.MetricsCounter[key] = value
	ms.Mutex.Unlock()
}

func (ms *MemStorage) AddCounter(key string, value counter) {
	ms.Mutex.Lock()
	ms.MetricsCounter[key] += value
	ms.Mutex.Unlock()
}

func (ms *MemStorage) GetGauge(key string) (gauge, bool) {
	ms.Mutex.RLock()
	value, ok := ms.MetricsGauge[key]
	ms.Mutex.RUnlock()
	if !ok {
		return 0, false
	}
	return value, true
}

func GetGauge(key string) (gauge, bool) {
	value, ok := Data.GetGauge(key)
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms *MemStorage) GetCounter(key string) (counter, bool) {
	ms.Mutex.RLock()
	value, ok := ms.MetricsCounter[key]
	ms.Mutex.RUnlock()
	if !ok {
		return 0, false
	}
	return value, true
}

func GetCounter(key string) (counter, bool) {
	value, ok := Data.GetCounter(key)
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms *MemStorage) DeleteGauge(key string) {
	ms.Mutex.Lock()
	delete(ms.MetricsGauge, key)
	ms.Mutex.Unlock()
}

func (ms *MemStorage) DeleteCounter(key string) {
	ms.Mutex.Lock()
	delete(ms.MetricsCounter, key)
	ms.Mutex.Unlock()
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики key
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter value
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge value
}

// var MetricsList []Metrics
// var mutexForMetricsList sync.Mutex

type MetricsList struct {
	MetricsList []Metrics
	Mutex       sync.RWMutex
}

var DataMetricsList = &MetricsList{
	MetricsList: []Metrics{},
	Mutex:       sync.RWMutex{},
}

func (ml *MetricsList) AppendMetrics(metric Metrics) {
	ml.Mutex.Lock()
	DataMetricsList.MetricsList = append(DataMetricsList.MetricsList, metric)
	ml.Mutex.Unlock()
}

func (ml *MetricsList) ClearMetrics() {
	ml.Mutex.Lock()
	DataMetricsList.MetricsList = nil
	ml.Mutex.Unlock()
}

func (ms *MemStorage) GetMetricsCounter() map[string]counter {
	ms.Mutex.RLock()
	copiedMetrics := Data.MetricsCounter
	ms.Mutex.RUnlock()
	return copiedMetrics
}

func (ms *MemStorage) GetMetricsGauge() map[string]gauge {
	ms.Mutex.RLock()
	copiedMetrics := Data.MetricsGauge
	ms.Mutex.RUnlock()
	return copiedMetrics
}

func (ms *MemStorage) GetMetrics() *MemStorage {
	ms.Mutex.RLock()
	copiedData := Data
	ms.Mutex.RUnlock()
	return copiedData
}

func (ms *MemStorage) SetMetrics(updatedData *MemStorage) {
	ms.Mutex.Lock()
	Data = updatedData
	ms.Mutex.Unlock()
}
