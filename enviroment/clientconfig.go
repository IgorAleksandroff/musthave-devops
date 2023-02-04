package enviroment

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	ClientEnvServerURL       = "ADDRESS"
	ClientEnvPollInterval    = "POLL_INTERVAL"
	ClientEnvReportInterval  = "REPORT_INTERVAL"
	ClientEnvHashKey         = "KEY"
	ClientEnvPublicCryptoKey = "CRYPTO_KEY"
	ClientEnvPublicCfgPath   = "CONFIG"

	ClientDefaultServerURL          = "localhost:8080"
	ClientDefaultPollInterval       = 2 * time.Second
	ClientDefaultReportInterval     = 10 * time.Second
	ClientDefaultEnvHashKey         = ""
	ClientDefaultEnvPublicCryptoKey = ""
	ClientDefaultCfgPath            = ""
)

type clientConfig struct {
	Host           string
	PollInterval   time.Duration
	ReportInterval time.Duration
	HashKey        string
	CryptoKeyPath  string
}

func NewClientConfig() clientConfig {
	log.Printf("os: %+v", os.Args)

	cfg := clientConfig{
		Host:           ClientDefaultServerURL,
		PollInterval:   ClientDefaultPollInterval,
		ReportInterval: ClientDefaultReportInterval,
		HashKey:        ClientDefaultEnvHashKey,
		CryptoKeyPath:  ClientDefaultEnvPublicCryptoKey,
	}

	hostFlag := flag.String("a", ClientDefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", ClientDefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", ClientDefaultReportInterval, "частота отправки метрик в секундах")
	hashKey := flag.String("k", ClientDefaultEnvHashKey, "ключ подписи метрик")
	cryptoKey := flag.String("crypto-key", ClientDefaultEnvPublicCryptoKey, "путь до файла с публичным ключом")
	cfgPathFlag := flag.String("c", ClientDefaultCfgPath, "адрес и порт сервера")

	flag.Parse()

	cfgJSONPath := GetEnvString(ClientEnvPublicCfgPath, *cfgPathFlag)

	if cfgJSONPath != ClientDefaultCfgPath {
		updateClientConfigByJSON(cfgJSONPath, &cfg)
	}

	// update Client Config by flags
	if hostFlag != nil && isFlagPassed("a") {
		cfg.Host = *hostFlag
	}
	if pollIntervalFlag != nil && isFlagPassed("p") {
		cfg.PollInterval = *pollIntervalFlag
	}
	if reportIntervalFlag != nil && isFlagPassed("r") {
		cfg.ReportInterval = *reportIntervalFlag
	}
	if hashKey != nil && isFlagPassed("k") {
		cfg.HashKey = *hashKey
	}
	if cryptoKey != nil && isFlagPassed("crypto-key") {
		cfg.CryptoKeyPath = *cryptoKey
	}

	// update Client Config by env
	cfg.Host = GetEnvString(ClientEnvServerURL, cfg.Host)
	cfg.PollInterval = GetEnvDuration(ClientEnvPollInterval, cfg.PollInterval)
	cfg.ReportInterval = GetEnvDuration(ClientEnvReportInterval, cfg.ReportInterval)
	cfg.HashKey = GetEnvString(ClientEnvHashKey, cfg.HashKey)
	cfg.CryptoKeyPath = GetEnvString(ClientEnvPublicCryptoKey, cfg.CryptoKeyPath)

	cfg.Host = "http://" + cfg.Host

	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}
