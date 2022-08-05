package serverconfig

import (
	"flag"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/cmd/environment"
)

const (
	EnvServerURL         = "ADDRESS"
	EnvStoreInterval     = "STORE_INTERVAL"
	EnvStoreFile         = "STORE_FILE"
	EnvRestore           = "RESTORE"
	DefaultServerURL     = "localhost:8080"
	DefaultStoreInterval = 300 * time.Second
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
)

type config struct {
	Host          string
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
}

func Read() config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", DefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", DefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", DefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	flag.Parse()

	cfg := config{
		Host:          environment.GetEnvString(EnvServerURL, *hostFlag),
		StoreInterval: environment.GetEnvDuration(EnvStoreInterval, *storeIntervalFlag),
		StorePath:     environment.GetEnvString(EnvStoreFile, *storePathFlag),
		Restore:       environment.GetEnvBool(EnvRestore, *restoreFlag),
	}

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}
