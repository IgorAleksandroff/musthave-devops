package entity

const GaugeTypeMetric = "gauge"
const CounterTypeMetric = "counter"

type MetricGauge struct {
	TypeMetric string
	Value      float64
}
type MetricCounter struct {
	TypeMetric string
	Value      int64
}
