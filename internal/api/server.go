package api

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

const EnvServerURL = "ADDRESS"
const EnvStoreInterval = "STORE_INTERVAL"
const EnvStoreFile = "STORE_FILE"
const EnvRestore = "RESTORE"
const DefaultServerURL = "localhost:8080"
const DefaultStoreInterval = 300
const DefaultStoreFile = "/tmp/devops-metrics-db.json"
const DefaultRestore = true

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type config struct {
	host          string
	storeInterval int
	storePath     string
	restore       bool
}

type server struct {
	cfg    config
	router *chi.Mux
}

func New() *server {
	r := chi.NewRouter()
	cfg := readConfig()

	return &server{router: r, cfg: cfg}
}

func (s *server) AddHandler(method, path string, handlerFn http.HandlerFunc) {
	s.router.MethodFunc(method, path, handlerFn)
}

func (s *server) Run() error {
	log.Printf("Start Server with config: %+v", s.cfg)
	return http.ListenAndServe(s.cfg.host, s.router)
}

type Server interface {
	Run() error
	AddHandler(method, path string, handlerFn http.HandlerFunc)
}

var _ Server = &server{}

func readConfig() config {
	return config{
		host:          getEnvString(EnvServerURL, DefaultServerURL),
		storeInterval: getEnvInt(EnvStoreInterval, DefaultStoreInterval),
		storePath:     getEnvString(EnvServerURL, DefaultStoreFile),
		restore:       getEnvBool(EnvRestore, DefaultRestore),
	}
}

func getEnvString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Println("empty env")
		return defaultValue
	}
	return value
}

func getEnvInt(envName string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(envName))
	if err != nil {
		log.Println(err.Error())
		return defaultValue
	}
	return value
}

func getEnvBool(envName string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(envName))
	if err != nil {
		log.Println(err.Error())
		return defaultValue
	}
	return value
}
