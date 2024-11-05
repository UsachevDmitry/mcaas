package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/agent/internal"
	"time"
	"sync"
)

func main() {
	// Создание группы ожидания
	var wg sync.WaitGroup

	// Получение конфигурации
	internal.GetConfig()

	// Добавление количества подпрограмм в группу ожидания
	wg.Add(3)

	// Запуск подпрограмм
	go func() {
		defer wg.Done()
		internal.UpdateData(time.Duration(*internal.PollInterval))
	}()

	go func() {
		defer wg.Done()
		internal.SendDataCounter(time.Duration(*internal.ReportInterval))
	}()

	go func() {
		defer wg.Done()
		internal.SendDataGauge(time.Duration(*internal.ReportInterval))
	}()

	// Ожидание завершения всех подпрограмм
	wg.Wait()

}
