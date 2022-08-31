package metricscollection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
	metricscolllectionrepo2 "github.com/IgorAleksandroff/musthave-devops/internal/metricscolllectionrepo"
)

//go:generate mockery --name MetricsCollection

type MetricsCollection interface {
	SaveMetric(value metricscollectionentity.Metrics)
	SaveCounterMetric(value metricscollectionentity.Metrics)
	GetMetric(name string) (*metricscollectionentity.Metrics, error)
	GetMetricsValue() map[string]string
	SaveMetrics(metrics []metricscollectionentity.Metrics)
	Ping() error
}

type (
	metricsCollection struct {
		repository metricscolllectionrepo2.Repository
	}

	Config struct {
		StorePath     string
		StoreInterval time.Duration
		Restore       bool
		AddressDB     string
	}
)

func NewMetricsCollection(ctx context.Context, cfg Config) (*metricsCollection, error) {
	repository := metricscolllectionrepo2.NewMemoRepository(ctx, metricscolllectionrepo2.MemoConfig{
		StorePath:     cfg.StorePath,
		StoreInterval: cfg.StoreInterval,
		Restore:       cfg.Restore,
	})

	if cfg.AddressDB != "" {
		repository, err := metricscolllectionrepo2.NewPGRepository(ctx, cfg.AddressDB)
		if err != nil {
			return nil, err
		}
		defer repository.Close()
	} else {
		repository.MemSync()
	}

	return &metricsCollection{repository: repository}, nil
}

func (u metricsCollection) SaveMetric(value metricscollectionentity.Metrics) {
	u.repository.SaveMetric(value)
}

func (u metricsCollection) SaveCounterMetric(value metricscollectionentity.Metrics) {
	if metric, err := u.repository.GetMetric(value.ID); err == nil {
		delta := *metric.Delta + *value.Delta
		value.Delta = &delta
	}

	u.repository.SaveMetric(value)
}

func (u metricsCollection) GetMetric(name string) (*metricscollectionentity.Metrics, error) {
	return u.repository.GetMetric(name)
}

func (u metricsCollection) GetMetricsValue() map[string]string {
	metrics := u.repository.GetMetrics()

	result := make(map[string]string, len(metrics))
	for name, m := range metrics {
		switch m.MType {
		case metricscollectionentity.CounterTypeMetric:
			if m.Delta != nil {
				result[name] = fmt.Sprintf("%v", *m.Delta)
			}
		case metricscollectionentity.GaugeTypeMetric:
			if m.Value != nil {
				result[name] = fmt.Sprintf("%v", *m.Value)
			}
		default:
			log.Printf("wrong type: %v of metric: %v", m.MType, m.ID)
		}
	}

	return result
}

func (u metricsCollection) SaveMetrics(metrics []metricscollectionentity.Metrics) {
	for _, metric := range metrics {
		if m, err := u.repository.GetMetric(metric.ID); err == nil && m.Delta != nil && metric.Delta != nil {
			delta := *m.Delta + *metric.Delta
			metric.Delta = &delta
		}

		u.repository.SaveMetric(metric)
	}
}

func (u metricsCollection) Ping() error {
	return u.repository.Ping()
}

func (u metricsCollection) Close() {
	u.repository.Close()
}

func (u metricsCollection) MemSync() {
	u.repository.MemSync()
}
