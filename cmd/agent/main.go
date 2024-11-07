package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/agent/internal"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	internal.GetConfig()
	wg.Add(3)

	go func() {
		internal.UpdateData(time.Duration(*internal.PollInterval))
		defer wg.Done()
	}()
	go func() {
		internal.SendDataCounter(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		internal.SendDataGauge(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	wg.Wait()
}
