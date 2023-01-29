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

	ServerDefaultServerURL        = "localhost:8080"
	ServerDefaultStoreInterval    = 300 * time.Second
	ServerDefaultStoreFile        = "/tmp/devops-metrics-db.json"
	ServerDefaultRestore          = true
	ServerDefaultEnvHashKey       = ""
	ServerDefaultEnvDB            = ""
	ServerDefaultPrivateCryptoKey = ""
)

type ServerConfig struct {
	Host          string
	HashKey       string
	CryptoKeyPath string
}

type Config struct {
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
	AddressDB     string
	ServerConfig
}

func NewServerConfig() Config {
	log.Println("os:", os.Args)

	hostFlag := flag.String("a", ServerDefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", ServerDefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", ServerDefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", ServerDefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	hashKeyFlag := flag.String("k", ServerDefaultEnvHashKey, "ключ подписи метрик")
	addressDBflag := flag.String("d", ServerDefaultEnvDB, "адрес подключения к БД")
	cryptoKeyFlag := flag.String("crypto-key", ServerDefaultPrivateCryptoKey, "путь до файла с приватным ключом")

	flag.Parse()

	storePathEnv := GetEnvString(ServerEnvStoreFile, *storePathFlag)
	restoreEnv := GetEnvBool(ServerEnvRestore, *restoreFlag)
	addressDBenv := GetEnvString(ServerEnvDB, *addressDBflag)

	cfg := Config{
		ServerConfig: ServerConfig{
			Host:          GetEnvString(ServerEnvServerURL, *hostFlag),
			HashKey:       GetEnvString(ServerEnvHashKey, *hashKeyFlag),
			CryptoKeyPath: GetEnvString(ServerEnvPrivateCryptoKey, *cryptoKeyFlag),
		},
		StoreInterval: GetEnvDuration(ServerEnvStoreInterval, *storeIntervalFlag),
		StorePath:     storePathEnv,
		Restore:       restoreEnv,
		AddressDB:     addressDBenv,
	}

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}
