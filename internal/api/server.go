package api

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

const (
	EnvServerURL         = "ADDRESS"
	EnvStoreInterval     = "STORE_INTERVAL"
	EnvStoreFile         = "STORE_FILE"
	EnvRestore           = "RESTORE"
	DefaultServerURL     = "localhost:8080"
	DefaultStoreInterval = 300
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultRestore       = true
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type Config struct {
	host          string
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
}

type server struct {
	cfg    Config
	router *chi.Mux
}

func New() *server {
	r := chi.NewRouter()
	cfg := readConfig()

	return &server{
		router: r,
		cfg:    cfg,
	}
}

func (s *server) AddHandler(method, path string, handlerFn http.HandlerFunc) {
	s.router.MethodFunc(method, path, handlerFn)
}

func (s *server) Run() error {
	log.Printf("Start Server with config: %+v", s.cfg)
	return http.ListenAndServe(s.cfg.host, s.router)
}

func (s *server) GetConfig() Config {
	return Config{
		host:          s.cfg.host,
		StoreInterval: s.cfg.StoreInterval,
		StorePath:     s.cfg.StorePath,
		Restore:       s.cfg.Restore,
	}
}

type Server interface {
	Run() error
	AddHandler(method, path string, handlerFn http.HandlerFunc)
}

var _ Server = &server{}

func readConfig() Config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", DefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", DefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", DefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	flag.Parse()

	return Config{
		host:          getEnvString(EnvServerURL, *hostFlag),
		StoreInterval: getEnvDuration(EnvStoreInterval, *storeIntervalFlag),
		StorePath:     getEnvString(EnvStoreFile, *storePathFlag),
		Restore:       getEnvBool(EnvRestore, *restoreFlag),
	}
}

func getEnvString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Printf("empty env: %s", envName)
		return defaultValue
	}
	return value
}

func getEnvDuration(envName string, defaultValue time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(envName))
	if err != nil {
		log.Printf("error of env %s: %s", envName, err.Error())
		return defaultValue
	}
	return value
}

func getEnvBool(envName string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(envName))
	if err != nil {
		log.Printf("error of env %s: %s", envName, err.Error())
		return defaultValue
	}
	return value
}
