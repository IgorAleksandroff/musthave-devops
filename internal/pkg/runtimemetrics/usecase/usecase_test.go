package usecase

import (
	"errors"
	"testing"

	"github.com/IgorAleksandroff/yp-musthave-devops/internal/api/services/devopsserver"
	mocks2 "github.com/IgorAleksandroff/yp-musthave-devops/internal/api/services/devopsserver/mocks"
	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/runtimemetrics"
	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/runtimemetrics/entity"
	"github.com/IgorAleksandroff/yp-musthave-devops/internal/pkg/runtimemetrics/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_usecase_SendMetrics(t *testing.T) {
	type fields struct {
		repository         runtimemetrics.Repository
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
				repoMock.On("SaveMetric", "PollCount", entity.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
					"name02",
				})
				repoMock.On("GetMetric", "name01").Return(entity.Metric{
					Value: entity.Gauge(0.1),
				}, nil).Once()
				repoMock.On("GetMetric", "name02").Return(entity.Metric{
					Value: entity.Counter(02),
				}, nil).Once()

				clientMock := &mocks2.Client{}
				endpoint01 := "/update/gauge/name01/0.1/"
				clientMock.On("DoPost", endpoint01, nil).Return(nil, nil).Once()
				endpoint02 := "/update/counter/name02/2/"
				clientMock.On("DoPost", endpoint02, nil).Return(nil, nil).Once()

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
				repoMock.On("SaveMetric", "PollCount", entity.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
				})
				repoMock.On("GetMetric", "name01").Return(entity.Metric{}, errors.New("err"))

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
				repoMock.On("SaveMetric", "PollCount", entity.Counter(0)).Return()
				repoMock.On("GetMetricsName").Return([]string{
					"name01",
					"name02",
				})
				repoMock.On("GetMetric", "name01").Return(entity.Metric{
					Value: entity.Gauge(0.1),
				}, nil).Once()
				repoMock.On("GetMetric", "name02").Return(entity.Metric{
					Value: entity.Counter(02),
				}, nil).Once()

				clientMock := &mocks2.Client{}
				endpoint01 := "/update/gauge/name01/0.1/"
				clientMock.On("DoPost", endpoint01, nil).Return(nil, errors.New("err")).Once()
				endpoint02 := "/update/counter/name02/2/"
				clientMock.On("DoPost", endpoint02, nil).Return(nil, nil).Once()

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
			u := New(
				f.repository,
				f.devopsServerClient,
			)
			u.SendMetrics()
		})
	}
}

func Test_usecase_UpdateMetrics(t *testing.T) {
	type fields struct {
		repository         runtimemetrics.Repository
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
				repoMock.On("SaveMetric", "PollCount", entity.Counter(0)).Return().Once()
				repoMock.On("GetMetric", "PollCount").Return(entity.Metric{
					Value: entity.Counter(99),
				}, nil)
				repoMock.On("SaveMetric", "PollCount", entity.Counter(100)).Return().Once()
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
				repoMock.On("SaveMetric", "PollCount", entity.Counter(0)).Return().Once()
				repoMock.On("GetMetric", "PollCount").Return(entity.Metric{}, errors.New("err"))
				repoMock.On("SaveMetric", "PollCount", entity.Counter(1)).Return().Once()
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
			u := New(
				f.repository,
				clientMock,
			)
			u.UpdateMetrics()
		})
	}
}
