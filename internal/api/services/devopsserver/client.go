package devopsserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//go:generate mockery --name "Client"

const EnvServerURL = "ADDRESS"
const DefaultServerURL = "127.0.0.1:8080"

type (
	client struct {
		host      string
		transport *http.Client
	}

	Client interface {
		Do(req *http.Request) (body []byte, err error)
		DoGet(url string) ([]byte, error)
		DoPost(url string, data interface{}) ([]byte, error)
	}
)

func NewClient() Client {
	return &client{
		host:      getEnvString(EnvServerURL, DefaultServerURL),
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
	req, err := http.NewRequest(http.MethodGet, c.host+path, nil)
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

	req, err := http.NewRequest(http.MethodPost, c.host+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set(`Content-Type`, `text/plain`)

	body, err := c.Do(req)

	return body, err
}

func getEnvString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Println("empty env")
		return defaultValue
	}
	return value
}
