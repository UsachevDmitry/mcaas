package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendDataGauge(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.GetMetricsGauge() {
			url := "http://" + *Addr + "/update/gauge/" + key + "/" + fmt.Sprintf("%.2f", float64(value))
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

func SendDataGaugeNewAPI(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for name, value := range Data.GetMetricsGauge() {
			url := "http://" + *Addr + "/update/"
			GaugeValueFloat64 := float64(value)
			var metrics = Metrics{
				ID:    name,
				MType: "gauge",
				Delta: nil,
				Value: &GaugeValueFloat64,
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
