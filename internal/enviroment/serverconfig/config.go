package serverconfig

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/enviroment"
)

const (
	EnvServerURL     = "ADDRESS"
	EnvStoreInterval = "STORE_INTERVAL"
	EnvStoreFile     = "STORE_FILE"
	EnvRestore       = "RESTORE"
	EnvHashKey       = "KEY"
	EnvDB            = "DATABASE_DSN"

	DefaultServerURL     = "localhost:8080"
	DefaultStoreInterval = 300 * time.Second
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
	DefaultEnvHashKey    = ""
	DefaultEnvDB         = ""
)

type config struct {
	Host          string
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
	HashKey       string
	AddressDB     string
}

func NewConfig() config {
	log.Println("os:", os.Args)
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", DefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", DefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", DefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	hashKeyFlag := flag.String("k", DefaultEnvHashKey, "ключ подписи метрик")
	addressDBflag := flag.String("d", DefaultEnvDB, "адрес подключения к БД")
	flag.Parse()

	storePathEnv := enviroment.GetEnvString(EnvStoreFile, *storePathFlag)
	restoreEnv := enviroment.GetEnvBool(EnvRestore, *restoreFlag)
	addressDBenv := enviroment.GetEnvString(EnvDB, *addressDBflag)
	//if addressDBenv != DefaultEnvDB {
	//	storePathEnv = ""
	//	restoreEnv = false
	//}

	cfg := config{
		Host:          enviroment.GetEnvString(EnvServerURL, *hostFlag),
		StoreInterval: enviroment.GetEnvDuration(EnvStoreInterval, *storeIntervalFlag),
		StorePath:     storePathEnv,
		Restore:       restoreEnv,
		HashKey:       enviroment.GetEnvString(EnvHashKey, *hashKeyFlag),
		AddressDB:     addressDBenv,
	}

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}
