package clientconfig

import (
	"flag"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/utils/enviroment"
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
	//hashKey := flag.String("k", DefaultEnvHashKey, "ключ подписи метрик")
	flag.Parse()

	cfg := config{
		Host:           "http://" + enviroment.GetEnvString(EnvServerURL, *hostFlag),
		PollInterval:   enviroment.GetEnvDuration(EnvPollInterval, *pollIntervalFlag),
		ReportInterval: enviroment.GetEnvDuration(EnvReportInterval, *reportIntervalFlag),
		HashKey:        "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21",
		//enviroment.GetEnvString(EnvHashKey, *hashKey),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
