package runtimemetrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/internal/enviroment/clientconfig"
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
	return "gauge"
}

func (Counter) GetType() string {
	return "counter"
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m *Metrics) CalcHash(value, key string) {
	if key == clientconfig.DefaultEnvHashKey {
		m.Hash = ""
		return
	}
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	m.Hash = fmt.Sprintf("%x", dst)
}
