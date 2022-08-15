package metricscollection

import (
	"encoding/json"
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

func CopyMetric(mIn Metrics) Metrics {
	mOut := mIn
	if mIn.Delta != nil {
		p := *mIn.Delta
		mOut.Delta = &p
	}
	if mIn.Value != nil {
		p := *mIn.Value
		mOut.Value = &p
	}
	return mOut
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
