package enviroment

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	ClientEnvServerURL       = "ADDRESS"
	ClientEnvPollInterval    = "POLL_INTERVAL"
	ClientEnvReportInterval  = "REPORT_INTERVAL"
	ClientEnvHashKey         = "KEY"
	ClientEnvPublicCryptoKey = "CRYPTO_KEY"
	ClientEnvPublicCfgPath   = "CONFIG"

	ClientDefaultServerURL      = "localhost:8080"
	ClientDefaultPollInterval   = 2 * time.Second
	ClientDefaultReportInterval = 10 * time.Second
)

type ClientConfig struct {
	Host             string
	NetInterfaceAddr string
	PollInterval     time.Duration
	ReportInterval   time.Duration
	HashKey          string
	CryptoKeyPath    string
	GRPSServerSocket string
}

func NewClientConfig() ClientConfig {
	log.Printf("os: %+v", os.Args)

	cfg := ClientConfig{
		Host:           ClientDefaultServerURL,
		PollInterval:   ClientDefaultPollInterval,
		ReportInterval: ClientDefaultReportInterval,
		HashKey:        "",
		CryptoKeyPath:  "",
	}

	hostFlag := flag.String("a", ClientDefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Duration("p", ClientDefaultPollInterval, "частота обновления метрик в секундах")
	reportIntervalFlag := flag.Duration("r", ClientDefaultReportInterval, "частота отправки метрик в секундах")
	hashKey := flag.String("k", "", "ключ подписи метрик")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с публичным ключом")
	cfgPathFlag := flag.String("c", "", "путь до json файла конфигурации сервера")
	socketFlag := flag.String("s", "", "если не указан gRPC сервер:порт, то используется HTTP клиент")

	flag.Parse()

	cfgJSONPath := GetEnvString(ClientEnvPublicCfgPath, *cfgPathFlag)
	if cfgJSONPath != "" {
		updateClientConfigByJSON(cfgJSONPath, &cfg)
	}

	// update only parameters of Client from json config by flags
	if hostFlag != nil && isFlagPassed("a") {
		cfg.Host = *hostFlag
	}
	if pollIntervalFlag != nil && isFlagPassed("p") {
		cfg.PollInterval = *pollIntervalFlag
	}
	if reportIntervalFlag != nil && isFlagPassed("r") {
		cfg.ReportInterval = *reportIntervalFlag
	}
	if cryptoKey != nil && isFlagPassed("crypto-key") {
		cfg.CryptoKeyPath = *cryptoKey
	}
	if socketFlag != nil && isFlagPassed("s") {
		cfg.GRPSServerSocket = *socketFlag
	}

	// update Client config by env, default is flag or json parameter
	cfg.Host = GetEnvString(ClientEnvServerURL, cfg.Host)
	cfg.PollInterval = GetEnvDuration(ClientEnvPollInterval, cfg.PollInterval)
	cfg.ReportInterval = GetEnvDuration(ClientEnvReportInterval, cfg.ReportInterval)
	cfg.HashKey = GetEnvString(ClientEnvHashKey, *hashKey)
	cfg.CryptoKeyPath = GetEnvString(ClientEnvPublicCryptoKey, cfg.CryptoKeyPath)

	cfg.Host = "http://" + cfg.Host

	ip, err := getInterfaceIP()
	if err != nil {
		log.Println("failed to get net interfaces:", err)
	}
	cfg.NetInterfaceAddr = ip

	log.Printf("Parsed Client config: %+v", cfg)

	return cfg
}

func getInterfaceIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}
	return "", errors.New("you aren't connected to the network")
}
