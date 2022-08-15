package repositorypg

import (
	"context"
	"errors"
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Pinger interface {
	Ping() error
}

type rep struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewRepository(ctx context.Context, db *pgxpool.Pool) *rep {
	return &rep{ctx: ctx, db: db}
}

func NewPinger(ctx context.Context, db *pgxpool.Pool) *rep {
	return &rep{ctx: ctx, db: db}
}

func (r *rep) SaveMetric(value metricscollection.Metrics) {
	//r.metricDB[value.ID] = value
	//if r.cfg.StoreInterval == 0 && r.cfg.StorePath != "" {
	//	if err := r.flushMemo(); err != nil {
	//		fmt.Printf("error to save metric in file %s: %s.\n", r.cfg.StorePath, err.Error())
	//	}
	//}
}

func (r *rep) GetMetric(name string) (*metricscollection.Metrics, error) {
	//if metric, ok := r.metricDB[name]; ok {
	//	return &metric, nil
	//}
	//
	return nil, fmt.Errorf("can not found a metric: %s", name)
}

func (r *rep) GetMetrics() map[string]metricscollection.Metrics {
	result := make(map[string]metricscollection.Metrics)
	//
	//for name, metric := range r.metricDB {
	//	result[name] = metricscollection.CopyMetric(metric)
	//}
	//
	return result
}

func (r *rep) Ping() error {
	if r.db == nil {
		return errors.New("DB isn't configured")
	}
	return r.db.Ping(r.ctx)
}
