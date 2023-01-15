// Package metricscollection saves metrics and provides read access to them.
package metricscollection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscolllectionrepo"
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
		repository metricscolllectionrepo.Repository
	}

	Config struct {
		StorePath     string
		StoreInterval time.Duration
		Restore       bool
		AddressDB     string
	}
)

func NewMetricsCollection(ctx context.Context, cfg Config) (*metricsCollection, error) {
	var repository metricscolllectionrepo.Repository

	var err error
	if cfg.AddressDB != "" {
		repository, err = metricscolllectionrepo.NewPGRepository(ctx, cfg.AddressDB)
		if err != nil {
			return nil, err
		}
	} else {
		repository = metricscolllectionrepo.NewMemoRepository(ctx, metricscolllectionrepo.MemoConfig{
			StorePath:     cfg.StorePath,
			StoreInterval: cfg.StoreInterval,
			Restore:       cfg.Restore,
		})
	}

	return &metricsCollection{repository: repository}, nil
}

// SaveMetric saves any metrics in repository.
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

// GetMetric provides read access to metric by name.
func (u metricsCollection) GetMetric(name string) (*metricscollectionentity.Metrics, error) {
	return u.repository.GetMetric(name)
}

// GetMetricsValue provides read access to metrics value by names.
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

// SaveMetrics saves batch metrics in repository
func (u metricsCollection) SaveMetrics(metrics []metricscollectionentity.Metrics) {
	for _, metric := range metrics {
		if m, err := u.repository.GetMetric(metric.ID); err == nil && m.Delta != nil && metric.Delta != nil {
			delta := *m.Delta + *metric.Delta
			metric.Delta = &delta
		}

		u.repository.SaveMetric(metric)
	}
}

// Ping checks connection to repository, return nil if ok.
func (u metricsCollection) Ping() error {
	return u.repository.Ping()
}

// Close connection to repository.
func (u metricsCollection) Close() {
	u.repository.Close()
}

// MemSync periodical copies metrics from inmemory to file.
func (u metricsCollection) MemSync() {
	u.repository.MemSync()
}
