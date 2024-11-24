package main

import (
	"github.com/UsachevDmitry/mcaas/internal/agent"
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
		internal.SendDataCounterNewAPI(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	go func() {
		internal.SendDataGaugeNewAPI(time.Duration(*internal.ReportInterval))
		defer wg.Done()
	}()
	wg.Wait()
}
