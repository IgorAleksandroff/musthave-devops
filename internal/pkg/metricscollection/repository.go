package metricscollection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockery --name Repository

type Repository interface {
	SaveMetric(value Metrics)
	GetMetric(name string) (*Metrics, error)
	GetMetrics() map[string]Metrics
	PingDB() error
}

type Config struct {
	StorePath     string
	AddressDB     string
	StoreInterval time.Duration
	Restore       bool
}

type rep struct {
	ctx      context.Context
	metricDB map[string]Metrics
	db       *pgxpool.Pool
	cfg      Config
}

func NewRepository(ctx context.Context, cfg Config, db *pgxpool.Pool) *rep {
	metricDB := make(map[string]Metrics)
	var err error

	if cfg.Restore && cfg.StorePath != "" {
		if metricDB, err = DownloadMetrics(cfg.StorePath); err != nil {
			log.Printf("error to restore metrics from %s: %s.\n", cfg.StorePath, err.Error())
		}
	}

	return &rep{ctx: ctx, metricDB: metricDB, cfg: cfg, db: db}
}

func (r *rep) SaveMetric(value Metrics) {
	r.metricDB[value.ID] = value
	if r.cfg.StoreInterval == 0 && r.cfg.StorePath != "" {
		if err := r.flushMemo(); err != nil {
			fmt.Printf("error to save metric in file %s: %s.\n", r.cfg.StorePath, err.Error())
		}
	}
}

func (r *rep) GetMetric(name string) (*Metrics, error) {
	if metric, ok := r.metricDB[name]; ok {
		return &metric, nil
	}

	return nil, fmt.Errorf("can not found a metric: %s", name)
}

func (r *rep) GetMetrics() map[string]Metrics {
	result := make(map[string]Metrics, len(r.metricDB))

	for name, metric := range r.metricDB {
		result[name] = CopyMetric(metric)
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

func (r *rep) PingDB() error {
	if r.cfg.AddressDB == "" {
		return errors.New("DB isn't configured")
	}
	return r.db.Ping(r.ctx)
}
