package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMetrics(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)

		url := "http://" + *Addr + "/updates/"

		if len(DataMetricsList.MetricsList) == 0 {
			continue
		} 

		jsonBody, err := json.Marshal(DataMetricsList.MetricsList)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		DataMetricsList.ClearMetrics() // Очишаем список
		compressedJSONBody, err := Compress(jsonBody)
		if err != nil {
			fmt.Println("Error compress jsonBody", err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedJSONBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		for i := 1; i < 6; i+=2 {
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				time.Sleep(time.Second * time.Duration(i)) // Задержка перед следующей попыткой
				continue
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				time.Sleep(5 * time.Second) // Задержка перед следующей попыткой
				continue
			}
			if i == 5 {
				fmt.Println("All retries exhausted, exiting...")
				break
			}
			break
		}
	}
}