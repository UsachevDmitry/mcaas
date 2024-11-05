package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func SendDataCounter(reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval * time.Second)
		for key, value := range Data.MetricsCounter {
			url := "http://" + *Addr + "/update/counter/" + key + "/" + strconv.FormatInt(int64(value), 10)

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
