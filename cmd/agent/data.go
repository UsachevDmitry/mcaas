package main

type gauge float64
type counter int64
type MemStorage struct {
	MetricsGauge   map[string]gauge
	MetricsCounter map[string]counter
}

var Data = MemStorage{
	MetricsGauge:   map[string]gauge{},
	MetricsCounter: map[string]counter{},
}

// Определение интерфейса для MemStorage
type MemStorageInterface interface {
	UpdateGauge(key string, value gauge)
	UpdateCounter(key string, value counter)
	AddCounter(key string, value counter)
	GetGauge(key string) (gauge, bool)
	GetCounter(key string) (counter, bool)
	DeleteGauge(key string)
	DeleteCounter(key string)
}

// Реализация методов интерфейса для MemStorage
func (ms MemStorage) UpdateGauge(key string, value gauge) {
	ms.MetricsGauge[key] = value
}

func (ms MemStorage) UpdateCounter(key string, value counter) {
	ms.MetricsCounter[key] = value
}

func (ms MemStorage) AddCounter(key string, value counter) {
	ms.MetricsCounter[key] += value
}

func (ms MemStorage) GetGauge(key string) (gauge, bool) {
	value, ok := ms.MetricsGauge[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms MemStorage) GetCounter(key string) (counter, bool) {
	value, ok := ms.MetricsCounter[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms MemStorage) DeleteGauge(key string) {
	delete(ms.MetricsGauge, key)
}

func (ms MemStorage) DeleteCounter(key string) {
	delete(ms.MetricsCounter, key)
}
