package enviroment

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const errorParseJSONConfig = "failed to parse client JSON config from path - %s: %s"

type clientJSONConfig struct {
	Host           string `json:"address,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
	ReportInterval string `json:"report_interval,omitempty"`
	CryptoKeyPath  string `json:"crypto_key,omitempty"`
}

func updateClientConfigByJSON(path string, cfg *clientConfig) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf(errorParseJSONConfig, path, err)

		return
	}
	defer f.Close()

	var cfgJSON clientJSONConfig
	err = json.NewDecoder(f).Decode(&cfgJSON)
	if err != nil {
		log.Printf(errorParseJSONConfig, path, err)
		return
	}

	cfg.Host = cfgJSON.Host
	cfg.CryptoKeyPath = cfgJSON.CryptoKeyPath

	if v, err := time.ParseDuration(cfgJSON.PollInterval); err != nil {
		log.Printf(errorParseJSONConfig, path, err)
	} else {
		cfg.PollInterval = v
	}

	if v, err := time.ParseDuration(cfgJSON.ReportInterval); err != nil {
		log.Printf(errorParseJSONConfig, path, err)
	} else {
		cfg.ReportInterval = v
	}
}