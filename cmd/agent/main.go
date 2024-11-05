package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/agent/internal"
	"time"
)

func main() {
	internal.GetConfig()
	go internal.UpdateData(time.Duration(*internal.PollInterval))
	go internal.SendDataCounter(time.Duration(*internal.ReportInterval))
	internal.SendDataGauge(time.Duration(*internal.ReportInterval))
}
