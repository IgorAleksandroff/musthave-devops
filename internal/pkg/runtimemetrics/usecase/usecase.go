package usecase

import (
	"github.com/IgorAleksandroff/musthave-devops/internal/api/services/devopsserver"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"
)

type usecase struct {
	repository         runtimemetrics.Repository
	devopsServerClient devopsserver.Client
}

func New(
	r runtimemetrics.Repository,
	client devopsserver.Client,
) *usecase {
	r.SaveMetric("PollCount", entity.Counter(0))

	return &usecase{
		repository:         r,
		devopsServerClient: client,
	}
}
