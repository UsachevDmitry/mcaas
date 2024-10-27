package main

import (
	"fmt"
	"runtime"
)

func main() {    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    Data.UpdateGauge("Alloc", gauge(m.Alloc))
    Data.UpdateGauge("BuckHashSys", gauge(m.BuckHashSys))
    Data.UpdateGauge("Frees",  gauge(m.Frees))
    Data.UpdateGauge("GCCPUFraction",  gauge(m.GCCPUFraction))
    Data.UpdateGauge("GCSys",  gauge(m.GCSys))
    Data.UpdateGauge("HeapAlloc",  gauge(m.HeapAlloc))
    Data.UpdateGauge("HeapIdle",  gauge(m.HeapIdle))
    Data.UpdateGauge("HeapInuse",  gauge(m.HeapInuse))
    Data.UpdateGauge("HeapObjects",  gauge(m.HeapObjects))
    Data.UpdateGauge("HeapReleased",  gauge(m.HeapReleased))
    Data.UpdateGauge("HeapSys",  gauge(m.HeapSys))
    Data.UpdateGauge("LastGC",  gauge(m.LastGC))
    Data.UpdateGauge("Lookups",  gauge(m.Lookups))
    Data.UpdateGauge("MCacheInuse",  gauge(m.MCacheInuse))
    Data.UpdateGauge("MCacheSys",  gauge(m.MCacheSys))
    Data.UpdateGauge("MSpanInuse",  gauge(m.MSpanInuse))
    Data.UpdateGauge("MSpanSys",  gauge(m.MSpanSys))
    Data.UpdateGauge("Mallocs",  gauge(m.Mallocs))
    Data.UpdateGauge("NextGC",  gauge(m.NextGC))
    Data.UpdateGauge("NumForcedGC",  gauge(m.NumForcedGC))
    Data.UpdateGauge("OtherSys",  gauge(m.OtherSys))
    Data.UpdateGauge("PauseTotalNs",  gauge(m.PauseTotalNs))
    Data.UpdateGauge("StackInuse",  gauge(m.StackInuse))
    Data.UpdateGauge("StackSys",  gauge(m.StackSys))
    Data.UpdateGauge("Sys",  gauge(m.Sys))
    Data.UpdateGauge("TotalAlloc",  gauge(m.TotalAlloc))

    fmt.Println(Data)
}
