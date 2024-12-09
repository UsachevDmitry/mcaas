package internal

import (
	"math/rand"
	"runtime"
	"time"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"fmt"
)

func UpdateData(pollInterval time.Duration) {
	var m runtime.MemStats
	UpdatedData := Data.GetMetrics()
	var i = 0
	for {
		runtime.ReadMemStats(&m)
		UpdatedData.UpdateGauge("Alloc", gauge(m.Alloc))
		UpdatedData.UpdateGauge("BuckHashSys", gauge(m.BuckHashSys))
		UpdatedData.UpdateGauge("Frees", gauge(m.Frees))
		UpdatedData.UpdateGauge("GCCPUFraction", gauge(m.GCCPUFraction))
		UpdatedData.UpdateGauge("GCSys", gauge(m.GCSys))
		UpdatedData.UpdateGauge("HeapAlloc", gauge(m.HeapAlloc))
		UpdatedData.UpdateGauge("HeapIdle", gauge(m.HeapIdle))
		UpdatedData.UpdateGauge("HeapInuse", gauge(m.HeapInuse))
		UpdatedData.UpdateGauge("HeapObjects", gauge(m.HeapObjects))
		UpdatedData.UpdateGauge("HeapReleased", gauge(m.HeapReleased))
		UpdatedData.UpdateGauge("HeapSys", gauge(m.HeapSys))
		UpdatedData.UpdateGauge("LastGC", gauge(m.LastGC))
		UpdatedData.UpdateGauge("Lookups", gauge(m.Lookups))
		UpdatedData.UpdateGauge("MCacheInuse", gauge(m.MCacheInuse))
		UpdatedData.UpdateGauge("MCacheSys", gauge(m.MCacheSys))
		UpdatedData.UpdateGauge("MSpanInuse", gauge(m.MSpanInuse))
		UpdatedData.UpdateGauge("MSpanSys", gauge(m.MSpanSys))
		UpdatedData.UpdateGauge("Mallocs", gauge(m.Mallocs))
		UpdatedData.UpdateGauge("NextGC", gauge(m.NextGC))
		UpdatedData.UpdateGauge("NumForcedGC", gauge(m.NumForcedGC))
		UpdatedData.UpdateGauge("NumGC", gauge(m.NumGC))
		UpdatedData.UpdateGauge("OtherSys", gauge(m.OtherSys))
		UpdatedData.UpdateGauge("PauseTotalNs", gauge(m.PauseTotalNs))
		UpdatedData.UpdateGauge("StackInuse", gauge(m.StackInuse))
		UpdatedData.UpdateGauge("StackSys", gauge(m.StackSys))
		UpdatedData.UpdateGauge("Sys", gauge(m.Sys))
		UpdatedData.UpdateGauge("TotalAlloc", gauge(m.TotalAlloc))
		i++
		UpdatedData.AddCounter("PollCount", counter(i))
		UpdatedData.UpdateGauge("RandomValue", gauge(rand.Float64()))
		Data.SetMetrics(UpdatedData)
		time.Sleep(pollInterval * time.Second)
	}
}

func UpdateDataMemCPU(pollInterval time.Duration) {
	UpdatedData := Data.GetMetrics()
	for {
		vm, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println(err)
		}
		cpuPercent, errcpu := cpu.Percent(0, false)
		if errcpu != nil {
			fmt.Println(errcpu)
		}
		UpdatedData.UpdateGauge("TotalMemory", gauge(vm.Total))
		UpdatedData.UpdateGauge("FreeMemory", gauge(vm.Free))
		for i, v := range cpuPercent {
			UpdatedData.UpdateGauge(fmt.Sprintf("CPUutilization%d\n", i), gauge(v))
		}
		Data.SetMetrics(UpdatedData)
		time.Sleep(pollInterval * time.Second)
	}
}

