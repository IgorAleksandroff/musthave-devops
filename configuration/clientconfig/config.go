package clientconfig

import (
	"flag"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/configuration"
)

const (
	EnvServerURL          = "ADDRESS"
	EnvPollInterval       = "POLL_INTERVAL"
	EnvReportInterval     = "REPORT_INTERVAL"
	EnvHashKey            = "KEY"
	DefaultServerURL      = "localhost:8080"
	DefaultPollInterval   = 2 * time.Second
	DefaultReportInterval = 10 * time.Second
	DefaultEnvHashKey     = ""
)

type config struct {
	Host           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	HashKey        string
}

func Read() config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", DefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", DefaultReportInterval, "частота отправки метрик в секундах")
	hashKey := flag.String("k", DefaultEnvHashKey, "адрес и порт сервера")
	flag.Parse()

	cfg := config{
		Host:           "http://" + configuration.GetEnvString(EnvServerURL, *hostFlag),
		PollInterval:   configuration.GetEnvDuration(EnvPollInterval, *pollIntervalFlag),
		ReportInterval: configuration.GetEnvDuration(EnvReportInterval, *reportIntervalFlag),
		HashKey:        configuration.GetEnvString(EnvHashKey, *hashKey),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
