package repository

import (
	"log"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

type config struct {
	storePath string
	syncMemo  bool
}

type rep struct {
	metricDB map[string]entity.Metrics
	cfg      config
}

func New(path string, syncMemo, restore bool) *rep {
	metricDB := make(map[string]entity.Metrics)
	var err error

	if restore && path != "" {
		if metricDB, err = entity.DownloadMetrics(path); err != nil {
			log.Printf("error to restore metrics from %s: %s.\n", path, err.Error())
		}
	}

	return &rep{
		metricDB: metricDB,
		cfg: config{
			storePath: path,
			syncMemo:  syncMemo,
		},
	}
}
