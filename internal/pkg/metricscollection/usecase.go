package metricscollection

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorymemo"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorypg"
)

//go:generate mockery --name Usecase

type Usecase interface {
	SaveMetric(value entity.Metrics)
	SaveCounterMetric(value entity.Metrics)
	GetMetric(name string) (*entity.Metrics, error)
	GetMetricsValue() map[string]string
	SaveMetrics(metrics []entity.Metrics)
	Ping() error
}

type (
	usecase struct {
		repository Repository
	}

	Config struct {
		StorePath     string
		StoreInterval time.Duration
		Restore       bool
		AddressDB     string
	}
)

func NewUsecase(ctx context.Context, cfg Config) (*usecase, error) {
	repository := repositorymemo.NewRepository(ctx, repositorymemo.Config{
		StorePath:     cfg.StorePath,
		StoreInterval: cfg.StoreInterval,
		Restore:       cfg.Restore,
	})

	if cfg.AddressDB != "" {
		repository, err := repositorypg.NewRepository(ctx, cfg.AddressDB)
		if err != nil {
			return nil, err
		}
		defer repository.Close()
	} else {
		repository.MemSync()
	}

	return &usecase{repository: repository}, nil
}

func (u usecase) SaveMetric(value entity.Metrics) {
	u.repository.SaveMetric(value)
}

func (u usecase) SaveCounterMetric(value entity.Metrics) {
	if metric, err := u.repository.GetMetric(value.ID); err == nil {
		delta := *metric.Delta + *value.Delta
		value.Delta = &delta
	}

	u.repository.SaveMetric(value)
}

func (u usecase) GetMetric(name string) (*entity.Metrics, error) {
	return u.repository.GetMetric(name)
}

func (u usecase) GetMetricsValue() map[string]string {
	metrics := u.repository.GetMetrics()

	result := make(map[string]string, len(metrics))
	for name, m := range metrics {
		switch m.MType {
		case entity.CounterTypeMetric:
			if m.Delta != nil {
				result[name] = fmt.Sprintf("%v", *m.Delta)
			}
		case entity.GaugeTypeMetric:
			if m.Value != nil {
				result[name] = fmt.Sprintf("%v", *m.Value)
			}
		default:
			log.Printf("wrong type: %v of metric: %v", m.MType, m.ID)
		}
	}

	return result
}

func (u usecase) SaveMetrics(metrics []entity.Metrics) {
	for _, metric := range metrics {
		if m, err := u.repository.GetMetric(metric.ID); err == nil && m.Delta != nil && metric.Delta != nil {
			delta := *m.Delta + *metric.Delta
			metric.Delta = &delta
		}

		u.repository.SaveMetric(metric)
	}
}

func (u usecase) Ping() error {
	return u.repository.Ping()
}

func (u usecase) Close() {
	u.repository.Close()
}
