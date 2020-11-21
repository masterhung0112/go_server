// Code generated by mockery v1.0.0. DO NOT EDIT.

// 'Regenerate

package mocks

import (
	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// SystemStore is an autogenerated mock type for the SystemStore type
type SystemStore struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *SystemStore) Get() (model.StringMap, error) {
	ret := _m.Called()

	var r0 model.StringMap
	if rf, ok := ret.Get(0).(func() model.StringMap); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.StringMap)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *SystemStore) GetByName(name string) (*model.System, error) {
	ret := _m.Called(name)

	var r0 *model.System
	if rf, ok := ret.Get(0).(func(string) *model.System); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.System)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertIfExists provides a mock function with given fields: system
func (_m *SystemStore) InsertIfExists(system *model.System) (*model.System, error) {
	ret := _m.Called(system)

	var r0 *model.System
	if rf, ok := ret.Get(0).(func(*model.System) *model.System); ok {
		r0 = rf(system)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.System)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.System) error); ok {
		r1 = rf(system)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteByName provides a mock function with given fields: name
func (_m *SystemStore) PermanentDeleteByName(name string) (*model.System, error) {
	ret := _m.Called(name)

	var r0 *model.System
	if rf, ok := ret.Get(0).(func(string) *model.System); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.System)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: system
func (_m *SystemStore) Save(system *model.System) error {
	ret := _m.Called(system)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.System) error); ok {
		r0 = rf(system)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveOrUpdate provides a mock function with given fields: system
func (_m *SystemStore) SaveOrUpdate(system *model.System) error {
	ret := _m.Called(system)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.System) error); ok {
		r0 = rf(system)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: system
func (_m *SystemStore) Update(system *model.System) error {
	ret := _m.Called(system)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.System) error); ok {
		r0 = rf(system)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
