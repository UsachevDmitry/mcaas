package internal

import (
	"testing"
	"time"
	"sync"
)

var Mutex sync.Mutex

func TestUpdateData(t *testing.T) {
	Data.UpdateGauge("Alloc", gauge(1))
	Data.UpdateGauge("BuckHashSys", gauge(1))
	Data.UpdateGauge("Frees", gauge(1))
	Data.UpdateGauge("GCCPUFraction", gauge(1))
	Data.UpdateGauge("GCSys", gauge(1))
	Data.UpdateGauge("HeapAlloc", gauge(1))
	Data.UpdateGauge("HeapIdle", gauge(1))
	Data.UpdateGauge("HeapInuse", gauge(1))
	Data.UpdateGauge("HeapObjects", gauge(1))
	Data.UpdateGauge("HeapReleased", gauge(1))
	Data.UpdateGauge("HeapSys", gauge(1))
	Data.UpdateGauge("LastGC", gauge(1))
	Data.UpdateGauge("Lookups", gauge(1))
	Data.UpdateGauge("MCacheInuse", gauge(1))
	Data.UpdateGauge("MCacheSys", gauge(1))
	Data.UpdateGauge("MSpanInuse", gauge(1))
	Data.UpdateGauge("MSpanSys", gauge(1))
	Data.UpdateGauge("Mallocs", gauge(1))
	Data.UpdateGauge("NextGC", gauge(1))
	Data.UpdateGauge("NumForcedGC", gauge(1))
	Data.UpdateGauge("OtherSys", gauge(1))
	Data.UpdateGauge("PauseTotalNs", gauge(1))
	Data.UpdateGauge("StackInuse", gauge(1))
	Data.UpdateGauge("StackSys", gauge(1))
	Data.UpdateGauge("Sys", gauge(1))
	Data.UpdateGauge("TotalAlloc", gauge(1))
	Data.AddCounter("PollCount", counter(1))
	Data.UpdateGauge("RandomValue", gauge(1))

	go UpdateData(time.Duration(2))
	time.Sleep(2 * time.Second)

	Mutex.Lock()
    metricsGauge := Data.MetricsGauge
	Mutex.Unlock()

	for key, value := range metricsGauge {
		if metricsGauge[key] == 1 {
			t.Errorf("Expected %v for key %s, got %v", metricsGauge[key], key, value)
		}
	}
	Mutex.Lock()
	metricsCounter := Data.MetricsCounter
	Mutex.Unlock()
	
	for key, value := range metricsCounter {
		if metricsCounter[key] == 1 {
			t.Errorf("Expected %v for key %s, got %v", metricsCounter[key], key, value)
		}
	}
}
