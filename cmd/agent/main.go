package main

import (
	"github.com/UsachevDmitry/mcaas/internal/agent"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	internal.GetConfig()
	jobs := make(chan internal.Metrics, 10)
	defer close(jobs)

	wg.Add(4 + *internal.RateLimit)
	go func() {
		internal.UpdateData(time.Duration(*internal.PollInterval), false)
		defer wg.Done()
	}()
	go func() {
		internal.UpdateDataMemCPU(time.Duration(*internal.PollInterval))
		defer wg.Done()
	}()
	go func() {
		internal.CollectDataCounterListNewAPI(jobs, time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		internal.CollectDataGaugeListNewAPI(jobs, time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()

	for w := 1; w <= *internal.RateLimit; w++ {
		go func() {
			internal.SendMetrics(jobs, time.Duration(*internal.ReportInterval))
			defer wg.Done()
		}()
	}

	wg.Wait()
}
