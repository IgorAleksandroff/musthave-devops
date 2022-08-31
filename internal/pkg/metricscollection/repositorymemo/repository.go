package repositorymemo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

type Config struct {
	StorePath     string
	StoreInterval time.Duration
	Restore       bool
}

type rep struct {
	ctx      context.Context
	metricDB map[string]entity.Metrics
	cfg      Config
}

func NewRepository(ctx context.Context, cfg Config) *rep {
	metricDB := make(map[string]entity.Metrics)
	var err error

	if cfg.Restore && cfg.StorePath != "" {
		if metricDB, err = entity.DownloadMetrics(cfg.StorePath); err != nil {
			log.Printf("error to restore metrics from %s: %s.\n", cfg.StorePath, err.Error())
		}
	}

	return &rep{ctx: ctx, metricDB: metricDB, cfg: cfg}
}

func (r *rep) SaveMetric(value entity.Metrics) {
	r.metricDB[value.ID] = value
	if r.cfg.StoreInterval == 0 && r.cfg.StorePath != "" {
		if err := r.flushMemo(); err != nil {
			fmt.Printf("error to save metric in file %s: %s.\n", r.cfg.StorePath, err.Error())
		}
	}
}

func (r *rep) GetMetric(name string) (*entity.Metrics, error) {
	if metric, ok := r.metricDB[name]; ok {
		return &metric, nil
	}

	return nil, fmt.Errorf("can not found a metric: %s", name)
}

func (r *rep) GetMetrics() map[string]entity.Metrics {
	result := make(map[string]entity.Metrics, len(r.metricDB))

	for name, metric := range r.metricDB {
		result[name] = metric.Copy()
	}

	return result
}

func (r *rep) MemSync() {
	go func() {
		ticker := time.NewTicker(r.cfg.StoreInterval)
		if r.cfg.StoreInterval == 0 {
			ticker.Stop()
		}
		defer ticker.Stop()
		for {
			select {
			case <-r.ctx.Done():
				err := r.flushMemo()
				if err != nil {
					log.Printf("can't save metrics, %s", err.Error())
				}
				return
			case <-ticker.C:
				err := r.flushMemo()
				if err != nil {
					log.Printf("can't save metrics, %s", err.Error())
				}
			}
		}
	}()
}

func (r *rep) Ping() error {
	return nil
}

func (r *rep) Close() {}

func (r *rep) flushMemo() error {
	file, err := os.OpenFile(r.cfg.StorePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(r.metricDB); err != nil {
		return err
	}

	return nil
}
