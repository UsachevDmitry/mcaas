package main

type gauge float64
//type counter int64
type MemStorage struct {
    Metrics map[string]gauge
}

var Data = MemStorage{
        Metrics: map[string]gauge{},
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