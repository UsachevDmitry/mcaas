package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func SendDataGauge(reportInterval time.Duration) {
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.MetricsGauge {
			url := "http://" + *Addr + "/update/gauge/" + key + "/" + fmt.Sprintf("%.2f", float64(value))

			req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
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

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				return
			}
		}
	}
}
