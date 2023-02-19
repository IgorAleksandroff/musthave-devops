package runtimemetrics

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
)

const (
	GaugeTypeMetric   = "gauge"
	CounterTypeMetric = "counter"
)

type Gauge float64
type Counter int64

type Getter interface {
	GetType() string
}

type Metric struct {
	Value Getter
}

func (Gauge) GetType() string {
	return GaugeTypeMetric
}

func (Counter) GetType() string {
	return CounterTypeMetric
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m *Metrics) CalcHash(value, key string) {
	if key == enviroment.ClientDefaultString {
		m.Hash = ""
		return
	}
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	m.Hash = string(dst)
}
