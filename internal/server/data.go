package internal

import (
	"context"
	"database/sql"
	"fmt"
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

type PostgresStorage struct {
	db *sql.DB
}

type Storage interface {
	UpdateGauge(key string, value gauge)
	UpdateCounter(key string, value counter)
	AddCounter(key string, value counter)
	GetGauge(key string) (gauge, bool)
	GetCounter(key string) (counter, bool)
	Close() error
	PingContext(ctx context.Context) error
}

type DatabaseConfig struct {
	Type string
}

var Config DatabaseConfig

func SelectStorage(config DatabaseConfig) (Storage, error) {
	//var db *MemStorage
	switch config.Type {
	case "mem":
		//db = Data
		return Data, nil
		// db := &MemStorage{
		// 	MetricsGauge:   map[string]gauge{},
		// 	MetricsCounter: map[string]counter{},
		// 	Mutex:          sync.Mutex{},
		// }
		//return db, nil
	case "postgres":
		db := &PostgresStorage{}
		errdb := db.Connect()
		if errdb != nil {
			panic(errdb)
		}
		db.CreateTables()
		return db, nil
	default:
		return nil, fmt.Errorf("неизвестная база данных: %s", config.Type)
	}
}

var FlagUsePosgresSQL bool
var Testdb *sql.DB

func (ms *MemStorage) UpdateGauge(key string, value gauge) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	//fmt.Println("UpdateGauge VALUE!!",value)
	ms.MetricsGauge[key] = value
	// fmt.Println("UpdateGauge!!",ms.MetricsGauge[key])
	// r,_ := GetGauge(key)
	// fmt.Println("R!!",r)
}

func UpdateGauge(key string, value gauge) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.UpdateGauge(key, value)
}

func (ms *MemStorage) UpdateCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] = value
	// fmt.Println("UpdateCounter value!!",value)
	// fmt.Println("UpdateCounter!!",ms.MetricsGauge[key])
	// r,_ := GetCounter(key)
	// fmt.Println("R!!",r)
}

func UpdateCounter(key string, value counter) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.UpdateCounter(key, value)
}

func (ms *MemStorage) AddCounter(key string, value counter) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.MetricsCounter[key] += value
}

func AddCounter(key string, value counter) {
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return
	}
	db.AddCounter(key, value)
}

func (ms *MemStorage) GetGauge(key string) (gauge, bool) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	// db, err := SelectStorage(Config)
	// if err != nil {
	// 	fmt.Println("Ошибка выбора базы данных:", err)
	// 	return 0, false
	// }
	//fmt.Println("INFOVALUE!!",ms.MetricsGauge[key])

	value, ok := Data.MetricsGauge[key]
	if !ok {
		return 0, false
	}
	return value, true
}

func GetGauge(key string) (gauge, bool) {
	var value gauge
	var ok bool
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return 0, false
	}
	// fmt.Println("INFOVALUEBeforeGetGAuge!!",Data.MetricsGauge[key])
	// fmt.Println("INFOVALUEBeforeGetGAuge!!",Data)
	// fmt.Println("INFOVALUEBeforeGetGAuge!!",Data.MetricsGauge[key])
	// fmt.Println("INFOVALUEBeforeGetGAuge!!",db)
	value, ok = db.GetGauge(key)
	if !ok {
		return 0, false
	}
	return value, true
}

func (ms *MemStorage) GetCounter(key string) (counter, bool) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	//fmt.Println("VALUEADDED!!",ms.MetricsCounter[key])
	value, ok := Data.MetricsCounter[key]
	if !ok {
		//fmt.Println("YYYYYYYYYYYYYYYYYYY!!!!!!!!!!!!!!!!!!",ok)
		return 0, false
	}
	return value, true
}

func GetCounter(key string) (counter, bool) {
	var Value counter
	var Ok bool
	db, err := SelectStorage(Config)
	if err != nil {
		fmt.Println("Ошибка выбора базы данных:", err)
		return 0, false
	}
	
	Value, Ok = db.GetCounter(key)
	if !Ok {
		fmt.Println("YYYYYYYYYYYYYYYYYYY!!!!!!!!!!!!!!!!!!")
		return 0, false
	}
	return Value, true
}

func (ms *MemStorage) Close() error {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	//Data.MetricsGauge = nil
	return nil
}

func (ms *MemStorage) PingContext(ctx context.Context) error {
	return fmt.Errorf("ping not work for with DB")
}

// func (ms *MemStorage) DeleteGauge(key string) {
// 	ms.Mutex.Lock()
// 	defer ms.Mutex.Unlock()
// 	delete(ms.MetricsGauge, key)
// }

// func DeleteGauge(key string) {
// 	Data.DeleteGauge(key)
// }

// func (ms *MemStorage) DeleteCounter(key string) {
// 	ms.Mutex.Lock()
// 	defer ms.Mutex.Unlock()
// 	delete(ms.MetricsCounter, key)
// }

// func DeleteCounter(key string) {
// 	Data.DeleteCounter(key)
// }

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики key
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter value
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge value
}
