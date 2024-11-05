package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/agent/internal"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	internal.GetConfig()
	wg.Add(3)

	go func() {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		internal.UpdateData(time.Duration(*internal.PollInterval))
	}()
	go func() {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		internal.SendDataCounter(time.Duration(*internal.ReportInterval))
	}()
	go func() {
		defer wg.Done()
		mutex.Lock()
		defer mutex.Unlock()
		internal.SendDataGauge(time.Duration(*internal.ReportInterval))
	}()
	wg.Wait()
}
