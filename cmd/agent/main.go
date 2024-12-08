package main

import (
	"github.com/UsachevDmitry/mcaas/internal/agent"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	internal.GetConfig()
	wg.Add(4)
	go func() {
		internal.UpdateData(time.Duration(*internal.PollInterval))
		defer wg.Done()
	}()
	go func() {
		internal.CollectDataCounterListNewAPI(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		internal.CollectDataGaugeListNewAPI(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		internal.SendMetrics(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	wg.Wait()
}
