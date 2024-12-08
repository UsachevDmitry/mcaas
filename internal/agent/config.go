package internal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	defaultAddr           = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultKey            = "secretkey"
)

var (
	Addr           = flag.String("a", defaultAddr, "Адрес HTTP-сервера")
	PollInterval   = flag.Int("p", defaultPollInterval, "pollInterval")
	ReportInterval = flag.Int("r", defaultReportInterval, "reportInterval")
	Key            = flag.String("k", defaultKey, "Ключ шифрования")
)

func GetConfig() {
	flag.Parse()
	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		*Addr = addrEnv
	}
	pollEnv := os.Getenv("POLL_INTERVAL")
	if pollEnv != "" {
		i, err := strconv.Atoi(pollEnv)
		if err != nil {
			log.Fatal(err)
		}
		*PollInterval = i
	}
	reportEnv := os.Getenv("REPORT_INTERVAL")
	if reportEnv != "" {
		i, err := strconv.Atoi(reportEnv)
		if err != nil {
			log.Fatal(err)
		}
		*ReportInterval = i
	}
	keyEnv := os.Getenv("KEY")
	if keyEnv != "" {
		*Key = keyEnv
	}
	fmt.Println("Адрес эндпоинта:", *Addr)
	fmt.Println("pollInterval:", *PollInterval)
	fmt.Println("reportInterval:", *ReportInterval)
	fmt.Println("key:", *Key)
}
