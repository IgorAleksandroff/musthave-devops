package grpchandler

import (
	"context"
	"fmt"
	"log"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/generated/rpc"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	rpc.UnimplementedMetricsCollectionServer
	metricsUC metricscollection.MetricsCollection
	hashKey   string
}

func New(metricsUC metricscollection.MetricsCollection, k string) *handler {
	return &handler{
		metricsUC: metricsUC,
		hashKey:   k,
	}
}

func (h handler) UpdateMetrics(ctx context.Context, r *rpc.UpdateMetricsRequest) (*rpc.UpdateMetricsResponse, error) {
	metrics := make([]metricscollectionentity.Metrics, 0)
	for _, m := range r.GetMetrics() {
		metricToSave := metricscollectionentity.Metrics{}

		switch m.GetMType() {
		case rpc.Metrics_COUNTER:
			metricToSave.MType = metricscollectionentity.CounterTypeMetric
			metricToSave.CalcHash(fmt.Sprintf("%s:counter:%d", m.GetId(), m.GetDelta()), h.hashKey)
		case rpc.Metrics_GAUGE:
			metricToSave.MType = metricscollectionentity.CounterTypeMetric
			metricToSave.CalcHash(fmt.Sprintf("%s:gauge:%f", m.GetId(), m.GetValue()), h.hashKey)
		default:
			return nil, status.Errorf(codes.NotFound, "unknown type metric: %s", m)
		}

		if h.hashKey != enviroment.ServerDefaultString && metricToSave.Hash != m.Hash {
			log.Println("hash isn't valid:", metricToSave.Hash, m)
			return nil, status.Errorf(codes.Aborted, "hash isn't valid: %s", m)
		}

		metricToSave.ID = m.GetId()
		metricToSave.Delta = &m.Delta
		metricToSave.Value = &m.Value

		metrics = append(metrics, metricToSave)
	}

	h.metricsUC.SaveMetrics(metrics)

	return &rpc.UpdateMetricsResponse{}, nil
}
