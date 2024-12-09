package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func SendMetrics(jobs chan Metrics, reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		var Metrics MetricsList
		for value := range jobs {
			Metrics.MetricsList = append(Metrics.MetricsList, value)
			if len(Metrics.MetricsList) > 10 {
				Send(Metrics.MetricsList)
				Metrics.ClearMetrics()
				break
			}
		}
	}
}


func Send(DataMetrics []Metrics) {
		url := "http://" + *Addr + "/updates/"

		if len(DataMetrics) == 0 {
			return
		}

		jsonBody, err := json.Marshal(DataMetrics)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		compressedJSONBody, err := Compress(jsonBody)
		if err != nil {
			fmt.Println("Error compress jsonBody", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedJSONBody))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		if *Key != "" {
			keyAndData := make([]byte, len(jsonBody)+len([]byte(*Key)))
			copy(keyAndData, jsonBody)
			copy(keyAndData, []byte(*Key))
			hash := sha256.Sum256(keyAndData)
			req.Header.Set("HashSHA256", fmt.Sprintf("%x", hash))
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		for i := 1; i < 6; i += 2 {
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
