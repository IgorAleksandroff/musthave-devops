package repository

import (
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (r rep) SaveMetric(value entity.Metrics) {
	r.metricDB[value.ID] = value
	if r.cfg.StoreInterval == 0 && r.cfg.StorePath != "" {
		if err := r.FlushMemo(); err != nil {
			fmt.Printf("error to save metric in file %s: %s.\n", r.cfg.StorePath, err.Error())
		}
	}
}
