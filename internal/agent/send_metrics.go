// package internal

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func SendMetrics(reportInterval time.Duration) {
// 	for {
// 		time.Sleep(reportInterval * time.Second)

// 		url := "http://" + *Addr + "/updates/"

// 		if len(DataMetricsList.MetricsList) == 0 {
// 			continue
// 		}

// 		jsonBody, err := json.Marshal(DataMetricsList.MetricsList)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			continue
// 		}
// 		DataMetricsList.ClearMetrics() // Очишаем список
// 		compressedJSONBody, err := Compress(jsonBody)
// 		if err != nil {
// 			fmt.Println("Error compress jsonBody", err)
// 			continue
// 		}
// 		var countRetry = 1
// 		for i := 1; i < 6; i += 2 {
// 			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(i)*time.Second)
// 			defer cancel()
// 			req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, url, bytes.NewBuffer(compressedJSONBody))
// 			if err != nil {
// 				cancel()
// 				fmt.Println("Error creating request:", err)
// 				if i == 5 {
// 					fmt.Println("All retries exhausted, exiting...")
// 					break
// 				}
// 				fmt.Printf("Retry %v...", countRetry)
// 				countRetry++
// 				continue
// 			}
// 			req.Header.Set("Content-Encoding", "gzip")
// 			req.Header.Set("Content-Type", "application/json")
// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			resp.Body.Close()
// 			cancel()
// 			if err != nil {
// 				fmt.Println("Error sending request:", err)
//              continue
// 			}
// 			if resp.StatusCode != http.StatusOK {
// 				fmt.Println("Error status:", resp.StatusCode)
//              continue
// 			}  else {
// 			  	break
// 			}
// 		}
// 	}
// }//test

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
		for i := 0; i < 3; i++ {
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				time.Sleep(5 * time.Second) // Задержка перед следующей попыткой
				continue
			} else {
				i=3
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Error status:", resp.StatusCode)
				time.Sleep(5 * time.Second) // Задержка перед следующей попыткой
				continue
			} else {
				i=3
			}
			if i == 2 {
				fmt.Println("All retries exhausted, exiting...")
				break
			}
		}
	}
}