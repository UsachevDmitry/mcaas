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
	defaultKey            = ""
	//defaultKey            = "secretkey"
	defaultRateLimit      = 5
)

var (
	Addr           = flag.String("a", defaultAddr, "Адрес HTTP-сервера")
	PollInterval   = flag.Int("p", defaultPollInterval, "pollInterval")
	ReportInterval = flag.Int("r", defaultReportInterval, "reportInterval")
	Key            = flag.String("k", defaultKey, "Ключ шифрования")
	RateLimit      = flag.Int("l", defaultRateLimit, "reportInterval")
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
	rateLimitEnv := os.Getenv("RATE_LIMIT")
	if rateLimitEnv != "" {
		i, err := strconv.Atoi(rateLimitEnv)
		if err != nil {
			log.Fatal(err)
		}
		*RateLimit = i
	}
	fmt.Println("Адрес эндпоинта:", *Addr)
	fmt.Println("pollInterval:", *PollInterval)
	fmt.Println("reportInterval:", *ReportInterval)
}
