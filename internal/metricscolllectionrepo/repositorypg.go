package metricscolllectionrepo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	queryCreateTable = `CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(64) UNIQUE,
		m_type VARCHAR(8) not null,
		delta BIGINT DEFAULT NULL,
		value double precision DEFAULT NULL,
		hash VARCHAR(64) DEFAULT NULL
	  )`
	querySave = `INSERT INTO metrics (id, m_type, delta, value, hash) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE 
		SET delta = $3, value = $4, hash = $5`
	queryGet        = `SELECT id, m_type, delta, value, hash FROM metrics WHERE id = $1`
	queryGetMetrics = `SELECT id, m_type, delta, value, hash FROM metrics`
)

type PGRep struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewPGRepository(ctx context.Context, addressDB string) (*PGRep, error) {
	conn, err := pgxpool.Connect(ctx, addressDB)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	log.Printf("connect to database: %v", conn.Config())

	repositoryPG := PGRep{ctx: ctx, db: conn}
	if err = repositoryPG.Init(); err != nil {
		return nil, fmt.Errorf("init database error: %v", err)
	}

	return &repositoryPG, nil
}

func (r *PGRep) Ping() error {
	if r.db == nil {
		return errors.New("DB isn't configured")
	}
	return r.db.Ping(r.ctx)
}

func (r *PGRep) Init() error {
	_, err := r.db.Exec(r.ctx, queryCreateTable)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRep) SaveMetric(value metricscollectionentity.Metrics) {
	_, err := r.db.Exec(r.ctx, querySave,
		value.ID,
		value.MType,
		value.Delta,
		value.Value,
		value.Hash,
	)
	if err != nil {
		log.Println(err)
	}
}

func (r *PGRep) GetMetric(name string) (*metricscollectionentity.Metrics, error) {
	var m metricscollectionentity.Metrics

	row := r.db.QueryRow(r.ctx, queryGet, name)
	if err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
		log.Printf("%v: can not found a metric in DB: %s\n", err, name)
		return nil, err
	}

	return &m, nil
}

func (r *PGRep) GetMetrics() map[string]metricscollectionentity.Metrics {
	result := make(map[string]metricscollectionentity.Metrics)
	rows, err := r.db.Query(r.ctx, queryGetMetrics)
	if err != nil {
		log.Printf("can not get all metrics: %v\n", err)
		return result
	}
	for rows.Next() {
		var m metricscollectionentity.Metrics
		if err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash); err != nil {
			log.Printf("can not scan a metric: %v\n", err)
			continue
		}

		result[m.ID] = m
	}
	if rows.Err() != nil {
		log.Printf("can not get all metrics: %v\n", err)
	}

	return result
}

func (r *PGRep) Close() {
	r.db.Close()
}
