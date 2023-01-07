package runtimemetrics_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	mocks2 "github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver/mocks"
	runtimemetrics2 "github.com/IgorAleksandroff/musthave-devops/internal/runtimemetrics"
	"github.com/IgorAleksandroff/musthave-devops/internal/runtimemetrics/mocks"
)

func Test_usecase_SendMetrics(t *testing.T) {
	type fields struct {
		repository         runtimemetrics2.Repository
		devopsServerClient devopsserver.Client
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "success",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
					"name02",
				})
				metric01 := runtimemetrics2.Metrics{
					ID:    "name01",
					MType: "gauge",
					Value: func() *float64 { v := 0.1; return &v }(),
				}
				metric02 := runtimemetrics2.Metrics{
					ID:    "name02",
					MType: "counter",
					Delta: func() *int64 { v := int64(02); return &v }(),
				}

				repoMock.On("GetMetric", "name01").Return(metric01, nil).Once()
				repoMock.On("GetMetric", "name02").Return(metric02, nil).Once()

				clientMock := &mocks2.Client{}
				clientMock.On("DoPost", "/update/", metric01).Return(nil, nil).Once()
				clientMock.On("DoPost", "/update/", metric02).Return(nil, nil).Once()

				return fields{
					repository:         repoMock,
					devopsServerClient: clientMock,
				}
			},
		},
		{
			name: "error_get_metric",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
				})
				repoMock.On("GetMetric", "name01").Return(runtimemetrics2.Metrics{}, errors.New("err"))

				clientMock := &mocks2.Client{}

				return fields{
					repository:         repoMock,
					devopsServerClient: clientMock,
				}
			},
		},
		{
			name: "error_post",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
					"name02",
				})
				metric01 := runtimemetrics2.Metrics{
					ID:    "name01",
					MType: "gauge",
					Value: func() *float64 { v := 0.1; return &v }(),
				}
				metric02 := runtimemetrics2.Metrics{
					ID:    "name02",
					MType: "counter",
					Delta: func() *int64 { v := int64(02); return &v }(),
				}
				repoMock.On("GetMetric", "name01").Return(metric01, nil).Once()
				repoMock.On("GetMetric", "name02").Return(metric02, nil).Once()

				clientMock := &mocks2.Client{}
				clientMock.On("DoPost", "/update/", metric01).Return(nil, errors.New("err")).Once()
				clientMock.On("DoPost", "/update/", metric02).Return(nil, nil).Once()

				return fields{
					repository:         repoMock,
					devopsServerClient: clientMock,
				}
			},
		},
	}
	for _, tt := range tests {
		f := tt.fields()
		t.Run(tt.name, func(t *testing.T) {
			u := runtimemetrics2.NewRuntimeMetrics(
				f.repository,
				f.devopsServerClient,
			)
			u.SendMetrics()
		})
	}
}

func Test_usecase_SendMetricsBatch(t *testing.T) {
	type fields struct {
		repository         runtimemetrics2.Repository
		devopsServerClient devopsserver.Client
	}

	metric01 := runtimemetrics2.Metrics{
		ID:    "name01",
		MType: "gauge",
		Value: func() *float64 { v := 0.1; return &v }(),
	}
	metric02 := runtimemetrics2.Metrics{
		ID:    "name02",
		MType: "counter",
		Delta: func() *int64 { v := int64(02); return &v }(),
	}
	testMetrics := []runtimemetrics2.Metrics{
		metric01,
		metric02,
	}

	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "success",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return().Once()
				repoMock.On("GetMetrics").Return(testMetrics)

				clientMock := &mocks2.Client{}
				clientMock.On("DoPost", "/updates/", testMetrics).Return(nil, nil).Once()

				return fields{
					repository:         repoMock,
					devopsServerClient: clientMock,
				}
			},
		},
		{
			name: "error_post",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return().Once()
				repoMock.On("GetMetrics").Return(testMetrics)

				clientMock := &mocks2.Client{}
				clientMock.On("DoPost", "/updates/", testMetrics).Return(nil, errors.New("err")).Once()

				return fields{
					repository:         repoMock,
					devopsServerClient: clientMock,
				}
			},
		},
	}
	for _, tt := range tests {
		f := tt.fields()
		t.Run(tt.name, func(t *testing.T) {
			u := runtimemetrics2.NewRuntimeMetrics(
				f.repository,
				f.devopsServerClient,
			)
			u.SendMetricsBatch()
		})
	}
}

func Test_usecase_UpdateMetrics(t *testing.T) {
	type fields struct {
		repository         runtimemetrics2.Repository
		devopsServerClient devopsserver.Client
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "success",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return().Once()
				repoMock.On("GetMetric", "PollCount").Return(runtimemetrics2.Metrics{
					Delta: func() *int64 { v := int64(99); return &v }(),
				}, nil)
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(100)).Return().Once()
				repoMock.On("SaveMetric", mock.MatchedBy(func(name string) bool {
					return name != "PollCount"
				}), mock.Anything).Return()

				return fields{
					repository: repoMock,
				}
			},
		},
		{
			name: "error_get_poll_count",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return().Once()
				repoMock.On("GetMetric", "PollCount").Return(runtimemetrics2.Metrics{}, errors.New("err"))
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(1)).Return().Once()
				repoMock.On("SaveMetric", mock.MatchedBy(func(name string) bool {
					return name != "PollCount"
				}), mock.Anything).Return()

				return fields{
					repository: repoMock,
				}
			},
		},
	}

	clientMock := &mocks2.Client{}

	for _, tt := range tests {
		f := tt.fields()
		t.Run(tt.name, func(t *testing.T) {
			u := runtimemetrics2.NewRuntimeMetrics(
				f.repository,
				clientMock,
			)
			u.UpdateMetrics()
		})
	}
}

func Test_usecase_UpdateUtilMetrics(t *testing.T) {
	type fields struct {
		repository         runtimemetrics2.Repository
		devopsServerClient devopsserver.Client
	}
	tests := []struct {
		name   string
		fields func() fields
	}{
		{
			name: "success",
			fields: func() fields {
				repoMock := &mocks.Repository{}
				repoMock.On("SaveMetric", "PollCount", runtimemetrics2.Counter(0)).Return().Once()
				repoMock.On("SaveMetric", "TotalMemory", mock.AnythingOfType("Gauge")).Return().Once()
				repoMock.On("SaveMetric", "FreeMemory", mock.AnythingOfType("Gauge")).Return().Once()
				repoMock.On("SaveMetric", "CPUutilization1", mock.AnythingOfType("Gauge")).Return().Once()

				return fields{
					repository: repoMock,
				}
			},
		},
	}

	clientMock := &mocks2.Client{}

	for _, tt := range tests {
		f := tt.fields()
		t.Run(tt.name, func(t *testing.T) {
			u := runtimemetrics2.NewRuntimeMetrics(
				f.repository,
				clientMock,
			)
			u.UpdateUtilMetrics()
		})
	}
}

func BenchmarkUpdateMetrics(b *testing.B) {
	u := runtimemetrics2.NewRuntimeMetrics(
		runtimemetrics2.NewRepository("test"),
		&mocks2.Client{},
	)

	for i := 0; i < b.N; i++ {
		u.UpdateMetrics()
	}
}
