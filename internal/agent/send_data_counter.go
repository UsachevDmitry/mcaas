package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func SendDataCounter(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.GetMetricsCounter() {
			url := "http://" + *Addr + "/update/counter/" + key + "/" + strconv.FormatInt(int64(value), 10)
			req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Content-Type", "text/plain")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				continue
			}
		}
	}
}

func SendDataCounterNewAPI(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for name, value := range Data.GetMetricsCounter() {
			url := "http://" + *Addr + "/update/"
			CounterValueInt64 := int64(value)
			var metrics = Metrics{
				ID:    name,
				MType: "counter",
				Delta: &CounterValueInt64,
				Value: nil,
			}
			jsonBody, err := json.Marshal(metrics)
			if err != nil {
				fmt.Println("Error:", err)
				continue

			}

			compressedJSONBody, _ := Compress(jsonBody)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedJSONBody))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				continue
			}
		}
	}
}
