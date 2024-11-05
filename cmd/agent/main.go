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
		mutex.Lock()
		defer mutex.Unlock()
		internal.UpdateData(time.Duration(*internal.PollInterval))
		defer wg.Done()
	}()
	go func() {		
		mutex.Lock()
		defer mutex.Unlock()
		internal.SendDataCounter(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		internal.SendDataGauge(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	wg.Wait()
}
