package internal

import (
	"context"
	"database/sql"
	"sync"
)

type gauge float64
type counter int64
type MemStorage struct {
	MetricsGauge   map[string]gauge
	MetricsCounter map[string]counter
	Mutex          sync.Mutex
}

var Data = &MemStorage{
	MetricsGauge:   map[string]gauge{},
	MetricsCounter: map[string]counter{},
	Mutex:          sync.Mutex{},
}

var DB *sql.DB
var FlagUsePosgresSQL bool

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

func UpdateGauge(key string, value gauge) {
	if FlagUsePosgresSQL {
		UpdateGaugeSQL(context.Background(), key, value)
	} else {
		Data.UpdateGauge(key, value)
	}
}

func (ms *MemStorage) UpdateCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] = value
}

func UpdateCounter(key string, value counter) {
	if FlagUsePosgresSQL {
		UpdateCounterSQL(context.Background(), key, value)
	} else {
		Data.UpdateCounter(key, value)
	}
}

func (ms *MemStorage) AddCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] += value
}

func AddCounter(key string, value counter) {
	if FlagUsePosgresSQL {
		AddCounterSQL(context.Background(), key, value)
	} else {
		Data.AddCounter(key, value)
	}
}

func (ms *MemStorage) GetGauge(key string) (gauge, bool) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	value, ok := ms.MetricsGauge[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func GetGauge(key string) (gauge, bool) {
	var Value gauge
	var Ok bool
	if FlagUsePosgresSQL {
		Value, Ok = GetGaugeSQL(context.Background(), key)
	} else {
		Value, Ok = Data.GetGauge(key)
	}
	if !Ok {
		return 0, false
	}
	return Value, true
}

func (ms *MemStorage) GetCounter(key string) (counter, bool) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	value, ok := ms.MetricsCounter[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func GetCounter(key string) (counter, bool) {
	var Value counter
	var Ok bool
	if FlagUsePosgresSQL {
		Value, Ok = GetCounterSQL(context.Background(), key)
	} else {
		Value, Ok = Data.GetCounter(key)
	}
	if !Ok {
		return 0, false
	}
	return Value, true
}

func (ms *MemStorage) DeleteGauge(key string) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	delete(ms.MetricsGauge, key)
}

func DeleteGauge(key string) {
	Data.DeleteGauge(key)
}

func (ms *MemStorage) DeleteCounter(key string) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	delete(ms.MetricsCounter, key)
}

func DeleteCounter(key string) {
	Data.DeleteCounter(key)
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики key
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter value
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge value
}
