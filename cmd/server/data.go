package main

type gauge float64
type counter int64
type MemStorage struct {
    MetricsGauge map[string]gauge
    MetricsCounter map[string]counter
}

var Data = MemStorage{
        MetricsGauge: map[string]gauge{},
        MetricsCounter: map[string]counter{},
}

// Определение интерфейса для MemStorage
type MemStorageInterface interface {
	AddGauge(key string, value gauge)
	AddCounter(key string, value counter)
	GetGauge(key string) (gauge, bool)
	GetCounter(key string) (counter, bool)
	DeleteGauge(key string)
	DeleteCounter(key string)
}

// Реализация методов интерфейса для MemStorage
func (ms MemStorage) AddGauge(key string, value gauge) {
	ms.MetricsGauge[key] = value
}

func (ms MemStorage) AddCounter(key string, value counter) {
	ms.MetricsCounter[key] = value
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

    // Data := MemStorage{
    //     Metrics: map[string]gauge{},
    //     //     "Alloc": 0,
    //     //     "BuckHashSys": 0,
    //     //     "Frees": 0,
    //     //     "GCCPUFraction": 0,
    //     //     "GCSys": 0,
    //     //     "HeapAlloc": 0,
    //     //     "HeapIdle": 0,
    //     //     "HeapInuse": 0,
    //     //     "HeapObjects": 0,
    //     //     "HeapReleased": 0,
    //     //     "HeapSys": 0,
    //     //     "LastGC": 0,
    //     //     "Lookups": 0,
    //     //     "MCacheInuse": 0,
    //     //     "MCacheSys": 0,
    //     //     "MSpanInuse": 0,
    //     //     "MSpanSys": 0,
    //     //     "Mallocs": 0,
    //     //     "NextGC": 0,
    //     //     "NumForcedGC": 0,
    //     //     "OtherSys": 0,
    //     //     "PauseTotalNs": 0,
    //     //     "StackInuse": 0,
    //     //     "StackSys": 0,
    //     //     "Sys": 0,
    //     //     "TotalAlloc": 0,
    //     // },
    // }