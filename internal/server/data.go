package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

type PostgresStorage struct {
	db *pgxpool.Pool
}

var DataDB = &PostgresStorage{}

type Storage interface {
	UpdateGauge(ctx context.Context, key string, value gauge)
	UpdateCounter(ctx context.Context, key string, value counter)
	AddCounter(ctx context.Context, key string, value counter)
	GetGauge(ctx context.Context, key string) (gauge, bool)
	GetCounter(ctx context.Context, key string) (counter, bool)
	Close()
	Ping(ctx context.Context) error
	CreateTableGauge(ctx context.Context)
	CreateTableCounter(ctx context.Context)
}

type DatabaseConfig struct {
	Type string
}

var Config DatabaseConfig

func SelectStorage(config DatabaseConfig) (Storage, error) {
	switch config.Type {
	case "mem":
		return Data, nil

	case "postgres":
		return DataDB, nil
	default:
		return nil, fmt.Errorf("неизвестная база данных: %s", config.Type)
	}
}

func (ms *MemStorage) UpdateGauge(ctx context.Context, key string, value gauge) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsGauge[key] = value
}

func UpdateGauge(ctx context.Context, key string, value gauge) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.UpdateGauge(ctx, key, value)
}

func (ms *MemStorage) UpdateCounter(ctx context.Context, key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] = value
}

func UpdateCounter(ctx context.Context, key string, value counter) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.UpdateCounter(ctx, key, value)
}

func (ms *MemStorage) AddCounter(ctx context.Context, key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] += value
}

func AddCounter(ctx context.Context, key string, value counter) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.AddCounter(ctx, key, value)
}

func (ms *MemStorage) GetGauge(ctx context.Context, key string) (gauge, bool) {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	value, ok := Data.MetricsGauge[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func GetGauge(ctx context.Context, key string) (gauge, bool) {
	var value gauge
	var ok bool
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return 0, false
	}
	value, ok = db.GetGauge(ctx, key)
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms *MemStorage) GetCounter(ctx context.Context, key string) (counter, bool) {
	ms.Mutex.RLock()
	defer ms.Mutex.RUnlock()
	value, ok := Data.MetricsCounter[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func GetCounter(ctx context.Context, key string) (counter, bool) {
	var Value counter
	var Ok bool
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return 0, false
	}
	Value, Ok = db.GetCounter(ctx, key)
	if !Ok {
		return 0, false
	}
	return Value, true
}

func (ms *MemStorage) Close() {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	Data.MetricsGauge = nil
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("ping not work for with DB")
}

func (ms *MemStorage) CreateTableGauge(ctx context.Context)   {}
func (ms *MemStorage) CreateTableCounter(ctx context.Context) {}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики key
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter value
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge value
}
