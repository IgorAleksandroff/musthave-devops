package metricscollectionentity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

const GaugeTypeMetric = "gauge"
const CounterTypeMetric = "counter"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func (m *Metrics) Copy() Metrics {
	mOut := *m
	if m.Delta != nil {
		p := *m.Delta
		mOut.Delta = &p
	}
	if m.Value != nil {
		p := *m.Value
		mOut.Value = &p
	}
	return mOut
}

func (m *Metrics) CalcHash(value, key string) {
	if key == "" {
		m.Hash = ""
		return
	}
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	m.Hash = fmt.Sprintf("%x", dst)
}

func DownloadMetrics(path string) (map[string]Metrics, error) {
	metricDB := make(map[string]Metrics)
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return metricDB, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&metricDB); err != nil {
		return metricDB, err
	}

	return metricDB, nil
}
