package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
	"flag"
)

const (
	defaultAddr = "localhost:8080"
	defaultPollInterval = 2
	defaultReportInterval = 10
)

var (
	addr = flag.String("a", defaultAddr, "Адрес HTTP-сервера")
	pollInterval = flag.Int("p", defaultPollInterval, "pollInterval")
	reportInterval = flag.Int("r", defaultReportInterval, "reportInterval")
)

func main() {
	flag.Parse()
	fmt.Println("Адрес эндпоинта:", *addr)
	fmt.Println("pollInterval:", *pollInterval)
	fmt.Println("reportInterval:", *reportInterval)

	go updateData(time.Duration(*pollInterval))
	sendDataCounter(time.Duration(*reportInterval))
	go sendDataGauge(time.Duration(*reportInterval))
}

func sendDataGauge(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.MetricsGauge {
			// Собираем строку с данными для отправки
			url := "http://" + *addr + "/update/gauge/" + key + "/" + fmt.Sprintf("%.2f", float64(value))
			fmt.Println(url)

			// Отправляем POST-запрос
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				return
			}
		}
	}
}

func sendDataCounter(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.MetricsCounter {
			// Собираем строку с данными для отправки
			url := "http://localhost:8080/update/counter/" + key + "/" + fmt.Sprintf("%v", int64(value))
			fmt.Println(url)

			// Отправляем POST-запрос
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				return
			}
		}
	}
}

func updateData(pollInterval time.Duration) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var i = 0
	for {
		Data.UpdateGauge("Alloc", gauge(m.Alloc))
		Data.UpdateGauge("BuckHashSys", gauge(m.BuckHashSys))
		Data.UpdateGauge("Frees", gauge(m.Frees))
		Data.UpdateGauge("GCCPUFraction", gauge(m.GCCPUFraction))
		Data.UpdateGauge("GCSys", gauge(m.GCSys))
		Data.UpdateGauge("HeapAlloc", gauge(m.HeapAlloc))
		Data.UpdateGauge("HeapIdle", gauge(m.HeapIdle))
		Data.UpdateGauge("HeapInuse", gauge(m.HeapInuse))
		Data.UpdateGauge("HeapObjects", gauge(m.HeapObjects))
		Data.UpdateGauge("HeapReleased", gauge(m.HeapReleased))
		Data.UpdateGauge("HeapSys", gauge(m.HeapSys))
		Data.UpdateGauge("LastGC", gauge(m.LastGC))
		Data.UpdateGauge("Lookups", gauge(m.Lookups))
		Data.UpdateGauge("MCacheInuse", gauge(m.MCacheInuse))
		Data.UpdateGauge("MCacheSys", gauge(m.MCacheSys))
		Data.UpdateGauge("MSpanInuse", gauge(m.MSpanInuse))
		Data.UpdateGauge("MSpanSys", gauge(m.MSpanSys))
		Data.UpdateGauge("Mallocs", gauge(m.Mallocs))
		Data.UpdateGauge("NextGC", gauge(m.NextGC))
		Data.UpdateGauge("NumForcedGC", gauge(m.NumForcedGC))
		Data.UpdateGauge("OtherSys", gauge(m.OtherSys))
		Data.UpdateGauge("PauseTotalNs", gauge(m.PauseTotalNs))
		Data.UpdateGauge("StackInuse", gauge(m.StackInuse))
		Data.UpdateGauge("StackSys", gauge(m.StackSys))
		Data.UpdateGauge("Sys", gauge(m.Sys))
		Data.UpdateGauge("TotalAlloc", gauge(m.TotalAlloc))
		i++
		Data.AddCounter("PollCount", counter(i))
		Data.UpdateGauge("RandomValue", gauge(rand.Float64()))
		time.Sleep(pollInterval * time.Second)
	}
}
