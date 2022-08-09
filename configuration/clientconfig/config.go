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
	DefaultServerURL      = "localhost:8080"
	DefaultPollInterval   = 2 * time.Second
	DefaultReportInterval = 10 * time.Second
)

type config struct {
	Host           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func Read() config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", DefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", DefaultReportInterval, "частота отправки метрик в секундах")
	flag.Parse()

	cfg := config{
		Host:           "http://" + configuration.GetEnvString(EnvServerURL, *hostFlag),
		PollInterval:   configuration.GetEnvDuration(EnvPollInterval, *pollIntervalFlag),
		ReportInterval: configuration.GetEnvDuration(EnvReportInterval, *reportIntervalFlag),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}