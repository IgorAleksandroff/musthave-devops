package enviroment

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	ClientEnvServerURL      = "ADDRESS"
	ClientEnvPollInterval   = "POLL_INTERVAL"
	ClientEnvReportInterval = "REPORT_INTERVAL"
	ClientEnvHashKey        = "KEY"

	ClientDefaultServerURL      = "localhost:8080"
	ClientDefaultPollInterval   = 2 * time.Second
	ClientDefaultReportInterval = 10 * time.Second
	ClientDefaultEnvHashKey     = ""
)

type ClientConfig struct {
	Host           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	HashKey        string
}

func NewClientConfig() ClientConfig {
	log.Println("os:", os.Args)
	hostFlag := flag.String("a", ClientDefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", ClientDefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", ClientDefaultReportInterval, "частота отправки метрик в секундах")
	hashKey := flag.String("k", ClientDefaultEnvHashKey, "ключ подписи метрик")
	flag.Parse()

	cfg := ClientConfig{
		Host:           "http://" + GetEnvString(ClientEnvServerURL, *hostFlag),
		PollInterval:   GetEnvDuration(ClientEnvPollInterval, *pollIntervalFlag),
		ReportInterval: GetEnvDuration(ClientEnvReportInterval, *reportIntervalFlag),
		HashKey:        GetEnvString(ClientEnvHashKey, *hashKey),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
