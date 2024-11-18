package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"encoding/json"
	"bytes"
)

func SendDataCounter(reportInterval time.Duration) {
	var mutex sync.Mutex
	
	mutex.Lock()
	defer mutex.Unlock()
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.MetricsCounter {
			url := "http://" + *Addr + "/update/counter/" + key + "/" + strconv.FormatInt(int64(value), 10)
			req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Content-Type", "text/plain")
			//req.Header.Set("Status-Code", "200")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}
			//defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				continue
			}
		}
	}
}

func SendDataCounterNewAPI(reportInterval time.Duration) {
	var mutex sync.Mutex
	
	mutex.Lock()
	defer mutex.Unlock()
	for {
		time.Sleep(reportInterval * time.Second)
		for name, value := range Data.MetricsCounter {
			url := "http://" + *Addr + "/update/"
			CounterValueInt64 := int64(value)
			var metrics = Metrics{
				ID: name,    
				MType: "counter",
				Delta: &CounterValueInt64,
				Value: nil,
			}
			jsonBody, err := json.Marshal(metrics)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			//req.Header.Set("Status-Code", "200")

			client := &http.Client{}
			resp, err := client.Do(req)
			
			if err != nil {
				fmt.Println("Error sending request:", err)
				continue
			}
			//defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				continue
			}
		}
	}
}
