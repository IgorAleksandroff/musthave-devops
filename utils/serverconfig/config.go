package serverconfig

import (
	"flag"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/utils"
)

const (
	EnvServerURL     = "ADDRESS"
	EnvStoreInterval = "STORE_INTERVAL"
	EnvStoreFile     = "STORE_FILE"
	EnvRestore       = "RESTORE"
	EnvHashKey       = "KEY"

	DefaultServerURL     = "localhost:8080"
	DefaultStoreInterval = 300 * time.Second
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
	DefaultEnvHashKey    = ""
)

type config struct {
	Host          string
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
	HashKey       string
}

func Read() config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", DefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", DefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", DefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	hashKey := flag.String("k", DefaultEnvHashKey, "ключ подписи метрик")
	flag.Parse()

	cfg := config{
		Host:          utils.GetEnvString(EnvServerURL, *hostFlag),
		StoreInterval: utils.GetEnvDuration(EnvStoreInterval, *storeIntervalFlag),
		StorePath:     utils.GetEnvString(EnvStoreFile, *storePathFlag),
		Restore:       utils.GetEnvBool(EnvRestore, *restoreFlag),
		HashKey:       utils.GetEnvString(EnvHashKey, *hashKey),
	}

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}
