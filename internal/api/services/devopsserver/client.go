package devopsserver

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/IgorAleksandroff/musthave-devops/internal/datacrypt"
	"github.com/pkg/errors"
)

//go:generate mockery --name "Client"

type (
	client struct {
		serverName, serverIPAddr string
		transport                *http.Client
		crypt                    datacrypt.Cypher
	}

	Client interface {
		DoGet(url string) ([]byte, error)
		DoPost(url string, data interface{}) ([]byte, error)
	}
)

func NewClient(serverName, netInterfaceAddr, cryptoKeyPathstring string) (Client, error) {
	var dc datacrypt.Cypher

	if cryptoKeyPathstring != "" {
		var err error

		dc, err = datacrypt.New(
			datacrypt.WithPublicKey(cryptoKeyPathstring),
			datacrypt.WithLabel("metrics"),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &client{
		serverName:   serverName,
		serverIPAddr: netInterfaceAddr,
		transport:    &http.Client{},
		crypt:        dc,
	}, nil
}

func (c client) do(req *http.Request) (body []byte, err error) {
	req.Header.Add("X-Real-IP", c.serverIPAddr)

	r, err := c.transport.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c client) DoGet(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.serverName+path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.do(req)

	return body, err
}

func (c client) DoPost(path string, data interface{}) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("payload marshal error")

		return nil, err
	}

	if c.crypt != nil {
		payload, err = c.crypt.Encrypt(payload)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, c.serverName+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set(`Content-Type`, `application/json`)

	body, err := c.do(req)

	return body, err
}

var _ Client = &client{}
