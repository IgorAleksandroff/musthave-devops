// Code generated by mockery v2.12.3. DO NOT EDIT.

package mocks

import (
	entity "github.com/IgorAleksandroff/musthave-devops/internal/pkg/runtimemetrics/entity"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetMetric provides a mock function with given fields: name
func (_m *Repository) GetMetric(name string) (entity.Metric, error) {
	ret := _m.Called(name)

	var r0 entity.Metric
	if rf, ok := ret.Get(0).(func(string) entity.Metric); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(entity.Metric)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetricsName provides a mock function with given fields:
func (_m *Repository) GetMetricsName() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// SaveMetric provides a mock function with given fields: name, value
func (_m *Repository) SaveMetric(name string, value entity.Getter) {
	_m.Called(name, value)
}

type NewRepositoryT interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t NewRepositoryT) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
