package devopsserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/datacrypt"
	"github.com/IgorAleksandroff/musthave-devops/internal/generated/rpc"
)

const (
	EndpointUpdate    = "/update/"
	EndpointUpdates   = "/updates/"
	counterTypeMetric = "counter"
	gaugeTypeMetric   = "gauge"
)

//go:generate mockery --name Client

type Client interface {
	Update(url string, data []Metrics) error
}

type (
	clientHTTP struct {
		serverName, serverIPAddr string
		transport                *http.Client
		crypt                    datacrypt.Cypher
	}

	clientGRPC struct {
		client       rpc.MetricsCollectionClient
		serverIPAddr string
	}
)

func NewClient(cfg enviroment.ClientConfig) (Client, error) {
	if cfg.GRPSServerSocket != enviroment.ClientDefaultString {
		return NewClientGRPS(cfg.NetInterfaceAddr, cfg.GRPSServerSocket)
	}

	return NewClientHTTP(cfg.Host, cfg.NetInterfaceAddr, cfg.CryptoKeyPath)
}

func NewClientGRPS(netInterfaceAddr, socket string) (*clientGRPC, error) {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	return &clientGRPC{
		client:       rpc.NewMetricsCollectionClient(conn),
		serverIPAddr: netInterfaceAddr,
	}, nil
}

func (c clientGRPC) Update(_ string, data []Metrics) error {
	metrics := make([]*rpc.Metrics, 0, len(data))
	for _, m := range data {
		metric := rpc.Metrics{
			Id:   m.ID,
			Hash: m.Hash,
		}

		switch m.MType {
		case counterTypeMetric:
			metric.MType = rpc.Metrics_COUNTER
		case gaugeTypeMetric:
			metric.MType = rpc.Metrics_GAUGE
		default:
			return errors.New("unknown metric type")
		}

		if m.Delta != nil {
			metric.Delta = *m.Delta
		}

		if m.Value != nil {
			metric.Value = *m.Value
		}

		metrics = append(metrics, &metric)
	}

	md := metadata.New(map[string]string{"X-Real-IP": c.serverIPAddr})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.client.UpdateMetrics(ctx, &rpc.UpdateMetricsRequest{
		Metrics: metrics,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return errors.Wrap(err, "parameters validation failed ")
			case codes.NotFound:
				return errors.Wrap(err, "incorrect MType")
			case codes.Aborted:
				return errors.Wrap(err, "incorrect hash")
			}
		}
	}

	return err
}

func NewClientHTTP(serverName, netInterfaceAddr, cryptoKeyPathString string) (*clientHTTP, error) {
	var dc datacrypt.Cypher

	if cryptoKeyPathString != "" {
		var err error

		dc, err = datacrypt.New(
			datacrypt.WithPublicKey(cryptoKeyPathString),
			datacrypt.WithLabel("metrics"),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &clientHTTP{
		serverName:   serverName,
		serverIPAddr: netInterfaceAddr,
		transport:    &http.Client{},
		crypt:        dc,
	}, nil
}

func (c clientHTTP) Update(url string, data []Metrics) error {
	var err error
	switch url {
	case EndpointUpdate:
		_, err = c.doPost(url, data[0])
	case EndpointUpdates:
		_, err = c.doPost(url, data)
	default:
		err = errors.New("unknown url")

	}

	return err
}

func (c clientHTTP) do(req *http.Request) (body []byte, err error) {
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

func (c clientHTTP) doGet(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.serverName+path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.do(req)

	return body, err
}

func (c clientHTTP) doPost(path string, data interface{}) ([]byte, error) {
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

var _ Client = &clientHTTP{}
