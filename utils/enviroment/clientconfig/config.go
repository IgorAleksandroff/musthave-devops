package clientconfig

import (
	"flag"
	"log"
	"os"
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
	log.Println("os:", os.Args)
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", DefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", DefaultReportInterval, "частота отправки метрик в секундах")
	hashKey := flag.String("k", DefaultEnvHashKey, "ключ подписи метрик")
	flag.Parse()

	cfg := config{
		Host:           "http://" + enviroment.GetEnvString(EnvServerURL, *hostFlag),
		PollInterval:   enviroment.GetEnvDuration(EnvPollInterval, *pollIntervalFlag),
		ReportInterval: enviroment.GetEnvDuration(EnvReportInterval, *reportIntervalFlag),
		HashKey:        enviroment.GetEnvString(EnvHashKey, *hashKey),
	}
	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
