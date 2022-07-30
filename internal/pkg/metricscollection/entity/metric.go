package entity

const GaugeTypeMetric = "gauge"
const CounterTypeMetric = "counter"

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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
