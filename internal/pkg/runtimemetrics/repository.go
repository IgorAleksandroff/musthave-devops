package runtimemetrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/configuration/clientconfig"
)

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(name string, value Getter)
	GetMetric(name string) (m Metrics, err error)
	GetMetricsName() []string
}

type rep struct {
	storage map[string]Metrics
	hashKey string
}

func NewRepository(key string) *rep {
	return &rep{storage: make(map[string]Metrics), hashKey: key}
}

func (r *rep) SaveMetric(name string, value Getter) {
	switch value := value.(type) {
	case Counter:
		valueInt64 := int64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Delta: &valueInt64,
			Hash:  getHash(fmt.Sprintf("%s:counter:%d", name, valueInt64), r.hashKey),
		}
	case Gauge:
		valueFloat64 := float64(value)
		r.storage[name] = Metrics{
			ID:    name,
			MType: value.GetType(),
			Value: &valueFloat64,
			Hash:  getHash(fmt.Sprintf("%s:gauge:%f", name, valueFloat64), r.hashKey),
		}
	}
}

func (r *rep) GetMetric(name string) (m Metrics, err error) {
	if metric, ok := r.storage[name]; ok {
		return metric, nil
	}

	return m, fmt.Errorf("can not get a metric: %s", name)
}

func (r *rep) GetMetricsName() []string {
	metricsName := make([]string, 0, len(r.storage))

	for metricName := range r.storage {
		metricsName = append(metricsName, metricName)
	}

	return metricsName
}

func getHash(value, key string) string {
	if key == clientconfig.DefaultEnvHashKey {
		return ""
	}
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	fmt.Printf("%x", dst)
	return fmt.Sprintf("%x", dst)
}
