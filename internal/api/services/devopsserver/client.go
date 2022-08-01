package devopsserver

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//go:generate mockery --name "Client"

const (
	EnvServerURL          = "ADDRESS"
	EnvPollInterval       = "POLL_INTERVAL"
	EnvReportInterval     = "REPORT_INTERVAL"
	DefaultServerURL      = "localhost:8080"
	DefaultPollInterval   = 2
	DefaultReportInterval = 10
)

type (
	Config struct {
		host           string
		PollInterval   time.Duration
		ReportInterval time.Duration
	}
	client struct {
		cfg       Config
		transport *http.Client
	}

	Client interface {
		Do(req *http.Request) (body []byte, err error)
		DoGet(url string) ([]byte, error)
		DoPost(url string, data interface{}) ([]byte, error)
		GetConfig() Config
	}
)

func NewClient() Client {
	cfg := readConfig()

	return &client{
		cfg:       cfg,
		transport: &http.Client{},
	}
}

func (c client) Do(req *http.Request) (body []byte, err error) {
	r, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c client) DoGet(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.cfg.host+path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.Do(req)

	return body, err
}

func (c client) DoPost(path string, data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("payload marshal error")

		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.cfg.host+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set(`Content-Type`, `application/json`)

	body, err := c.Do(req)

	return body, err
}

func (c client) GetConfig() Config {
	return Config{
		host:           c.cfg.host,
		PollInterval:   c.cfg.PollInterval,
		ReportInterval: c.cfg.ReportInterval,
	}
}

var _ Client = &client{}

func readConfig() Config {
	hostFlag := flag.String("a", DefaultServerURL, "адрес и порт сервера")
	pollIntervalFlag := flag.Int("p", DefaultPollInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	reportIntervalFlag := flag.Int("r", DefaultReportInterval, "строка, имя файла, где хранятся значения")
	flag.Parse()

	return Config{
		host:           "http://" + getEnvString(EnvServerURL, *hostFlag),
		PollInterval:   time.Duration(getEnvInt(EnvPollInterval, *pollIntervalFlag)) * time.Second,
		ReportInterval: time.Duration(getEnvInt(EnvReportInterval, *reportIntervalFlag)) * time.Second,
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
	value, err := strconv.Atoi(strings.TrimRight(os.Getenv(envName), "s"))
	if err != nil {
		log.Printf("error of env %s: %s", envName, err.Error())
		return defaultValue
	}
	return value
}
