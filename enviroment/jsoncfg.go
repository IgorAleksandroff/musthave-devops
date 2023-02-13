package enviroment

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const errorParseClientJSONConfig = "failed to parse client JSON config from path - %s: %s"
const errorParseServerJSONConfig = "failed to parse server JSON config from path - %s: %s"

type clientJSONConfig struct {
	Host           string `json:"address,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
	ReportInterval string `json:"report_interval,omitempty"`
	CryptoKeyPath  string `json:"crypto_key,omitempty"`
}

type serverJSONConfig struct {
	Host          string `json:"address,omitempty"`
	Restore       string `json:"restore,omitempty"`
	StoreInterval string `json:"store_interval,omitempty"`
	StorePath     string `json:"store_file,omitempty"`
	AddressDB     string `json:"database_dsn,omitempty"`
	CryptoKeyPath string `json:"crypto_key,omitempty"`
	TrustedSubnet string `json:"trusted_subnet,omitempty"`
}

func updateClientConfigByJSON(path string, cfg *clientConfig) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf(errorParseClientJSONConfig, path, err)

		return
	}
	defer f.Close()

	var cfgJSON clientJSONConfig
	err = json.NewDecoder(f).Decode(&cfgJSON)
	if err != nil {
		log.Printf(errorParseClientJSONConfig, path, err)
		return
	}

	cfg.Host = cfgJSON.Host
	cfg.CryptoKeyPath = cfgJSON.CryptoKeyPath

	if v, err := time.ParseDuration(cfgJSON.PollInterval); err != nil {
		log.Printf(errorParseClientJSONConfig, path, err)
	} else {
		cfg.PollInterval = v
	}

	if v, err := time.ParseDuration(cfgJSON.ReportInterval); err != nil {
		log.Printf(errorParseClientJSONConfig, path, err)
	} else {
		cfg.ReportInterval = v
	}
}

func updateServerConfigByJSON(path string, cfg *config) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf(errorParseServerJSONConfig, path, err)

		return
	}
	defer f.Close()

	var cfgJSON serverJSONConfig
	err = json.NewDecoder(f).Decode(&cfgJSON)
	if err != nil {
		log.Printf(errorParseServerJSONConfig, path, err)
		return
	}

	cfg.Host = cfgJSON.Host
	cfg.StorePath = cfgJSON.StorePath
	cfg.AddressDB = cfgJSON.AddressDB
	cfg.CryptoKeyPath = cfgJSON.CryptoKeyPath

	if v, err := strconv.ParseBool(cfgJSON.Restore); err != nil {
		log.Printf(errorParseServerJSONConfig, path, err)
	} else {
		cfg.Restore = v
	}

	if v, err := time.ParseDuration(cfgJSON.StoreInterval); err != nil {
		log.Printf(errorParseServerJSONConfig, path, err)
	} else {
		cfg.StoreInterval = v
	}

	if _, v, err := net.ParseCIDR(cfgJSON.TrustedSubnet); err != nil {
		log.Printf(errorParseServerJSONConfig, path, err)
	} else {
		cfg.subnet = v
	}
}
