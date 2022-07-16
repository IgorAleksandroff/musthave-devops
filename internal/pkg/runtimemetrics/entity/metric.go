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
