package repository

import (
	"context"
	"log"
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

func New(ctx context.Context, cfg Config) *rep {
	metricDB := make(map[string]entity.Metrics)
	var err error

	if cfg.Restore && cfg.StorePath != "" {
		if metricDB, err = entity.DownloadMetrics(cfg.StorePath); err != nil {
			log.Printf("error to restore metrics from %s: %s.\n", cfg.StorePath, err.Error())
		}
	}

	return &rep{ctx: ctx, metricDB: metricDB, cfg: cfg}
}
