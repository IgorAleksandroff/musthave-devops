package enviroment

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	ServerEnvServerURL        = "ADDRESS"
	ServerEnvStoreInterval    = "STORE_INTERVAL"
	ServerEnvStoreFile        = "STORE_FILE"
	ServerEnvRestore          = "RESTORE"
	ServerEnvHashKey          = "KEY"
	ServerEnvDB               = "DATABASE_DSN"
	ServerEnvPrivateCryptoKey = "CRYPTO_KEY"
	ServerEnvPublicCfgPath    = "CONFIG"

	ServerDefaultServerURL        = "localhost:8080"
	ServerDefaultStoreInterval    = 300 * time.Second
	ServerDefaultStoreFile        = "/tmp/devops-metrics-db.json"
	ServerDefaultRestore          = true
	ServerDefaultEnvHashKey       = ""
	ServerDefaultEnvDB            = ""
	ServerDefaultPrivateCryptoKey = ""
	ServerDefaultCfgPath          = ""
)

type ServerConfig struct {
	Host          string
	HashKey       string
	CryptoKeyPath string
}

type config struct {
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
	AddressDB     string
	ServerConfig
}

func NewServerConfig() config {
	log.Printf("os: %+v", os.Args)

	cfg := config{
		ServerConfig: ServerConfig{
			Host:          ServerDefaultServerURL,
			HashKey:       ServerDefaultEnvHashKey,
			CryptoKeyPath: ServerDefaultPrivateCryptoKey,
		},
		StoreInterval: ServerDefaultStoreInterval,
		StorePath:     ServerDefaultStoreFile,
		Restore:       ServerDefaultRestore,
		AddressDB:     ServerDefaultEnvDB,
	}

	hostFlag := flag.String("a", ServerDefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", ServerDefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", ServerDefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", ServerDefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	hashKeyFlag := flag.String("k", ServerDefaultEnvHashKey, "ключ подписи метрик")
	addressDBflag := flag.String("d", ServerDefaultEnvDB, "адрес подключения к БД")
	cryptoKeyFlag := flag.String("crypto-key", ServerDefaultPrivateCryptoKey, "путь до файла с приватным ключом")
	cfgPathFlag := flag.String("c", ServerDefaultCfgPath, "путь до json файла конфигурации сервера")

	flag.Parse()

	cfgJSONPath := GetEnvString(ServerEnvPublicCfgPath, *cfgPathFlag)
	if cfgJSONPath != ClientDefaultCfgPath {
		updateServerConfigByJSON(cfgJSONPath, &cfg)
	}

	// update Client config by flags
	if hostFlag != nil && isFlagPassed("a") {
		cfg.Host = *hostFlag
	}
	if restoreFlag != nil && isFlagPassed("r") {
		cfg.Restore = *restoreFlag
	}
	if storeIntervalFlag != nil && isFlagPassed("i") {
		cfg.StoreInterval = *storeIntervalFlag
	}
	if storePathFlag != nil && isFlagPassed("f") {
		cfg.StorePath = *storePathFlag
	}
	if addressDBflag != nil && isFlagPassed("d") {
		cfg.AddressDB = *addressDBflag
	}
	if cryptoKeyFlag != nil && isFlagPassed("crypto-key") {
		cfg.CryptoKeyPath = *cryptoKeyFlag
	}

	// update Client config by env, default is flag or json parameter
	cfg.Host = GetEnvString(ServerEnvServerURL, cfg.Host)
	cfg.HashKey = GetEnvString(ServerEnvHashKey, *hashKeyFlag)
	cfg.CryptoKeyPath = GetEnvString(ServerEnvPrivateCryptoKey, cfg.CryptoKeyPath)
	cfg.StoreInterval = GetEnvDuration(ServerEnvStoreInterval, cfg.StoreInterval)
	cfg.StorePath = GetEnvString(ServerEnvStoreFile, cfg.StorePath)
	cfg.Restore = GetEnvBool(ServerEnvRestore, cfg.Restore)
	cfg.AddressDB = GetEnvString(ServerEnvDB, cfg.AddressDB)

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}
