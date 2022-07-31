package entity

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
}
