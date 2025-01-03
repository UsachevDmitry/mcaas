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
    defer ms.Mutex.Unlock()
	ms.MetricsGauge[key] = value
}

func (ms *MemStorage) UpdateCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] = value
}

func (ms *MemStorage) AddCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] += value
}

func (ms *MemStorage) GetGauge(key string) (gauge, bool) {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	value, ok := ms.MetricsGauge[key]
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
	defer ms.Mutex.RUnlock()
	value, ok := ms.MetricsCounter[key]
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
	defer ms.Mutex.Unlock()
	delete(ms.MetricsGauge, key)
}

func (ms *MemStorage) DeleteCounter(key string) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	delete(ms.MetricsCounter, key)
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
	defer ml.Mutex.Unlock()
	DataMetricsList.MetricsList = append(DataMetricsList.MetricsList, metric)
}

func (ml *MetricsList) ClearMetrics() {
	ml.Mutex.Lock()
	defer ml.Mutex.Unlock()
	DataMetricsList.MetricsList = nil
}

func (ms *MemStorage) GetMetricsCounter() map[string]counter {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	copiedMetrics := Data.MetricsCounter
	return copiedMetrics
}

func (ms *MemStorage) GetMetricsGauge() map[string]gauge {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	copiedMetrics := Data.MetricsGauge
	return copiedMetrics
}

func (ms *MemStorage) GetMetrics() *MemStorage {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	copiedData := Data
	return copiedData
}

func (ms *MemStorage) SetMetrics(updatedData *MemStorage) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	Data = updatedData
}
