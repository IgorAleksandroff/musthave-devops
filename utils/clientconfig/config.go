package clientconfig

import (
	"flag"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/utils"
)

const (
	EnvServerURL      = "ADDRESS"
	EnvPollInterval   = "POLL_INTERVAL"
	EnvReportInterval = "REPORT_INTERVAL"
	EnvHashKey        = "KEY"

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
	hashKey := flag.String("k", DefaultEnvHashKey, "ключ подписи метрик")
	flag.Parse()

	cfg := config{
		Host:           "http://" + utils.GetEnvString(EnvServerURL, *hostFlag),
		PollInterval:   utils.GetEnvDuration(EnvPollInterval, *pollIntervalFlag),
		ReportInterval: utils.GetEnvDuration(EnvReportInterval, *reportIntervalFlag),
		HashKey:        utils.GetEnvString(EnvHashKey, *hashKey),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
