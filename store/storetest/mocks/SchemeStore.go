// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/masterhung0112/go_server/model"
	mock "github.com/stretchr/testify/mock"
)

// SchemeStore is an autogenerated mock type for the SchemeStore type
type SchemeStore struct {
	mock.Mock
}

// Save provides a mock function with given fields: scheme
func (_m *SchemeStore) Save(scheme *model.Scheme) (*model.Scheme, error) {
	ret := _m.Called(scheme)

	var r0 *model.Scheme
	if rf, ok := ret.Get(0).(func(*model.Scheme) *model.Scheme); ok {
		r0 = rf(scheme)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Scheme)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Scheme) error); ok {
		r1 = rf(scheme)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
