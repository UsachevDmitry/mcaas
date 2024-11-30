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

		jsonBody, err := json.Marshal(MetricsList)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		ClearMetrics() // Очишаем список
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