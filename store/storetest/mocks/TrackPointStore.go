// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `npm run task store_mocks`.

package mocks

import (
	model "github.com/masterhung0112/hk_server/v5/model"
	mock "github.com/stretchr/testify/mock"
)

// TrackPointStore is an autogenerated mock type for the TrackPointStore type
type TrackPointStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: trackPointId
func (_m *TrackPointStore) Get(trackPointId string) (*model.TrackPoint, error) {
	ret := _m.Called(trackPointId)

	var r0 *model.TrackPoint
	if rf, ok := ret.Get(0).(func(string) *model.TrackPoint); ok {
		r0 = rf(trackPointId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackPoint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(trackPointId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTargetId provides a mock function with given fields: targetId
func (_m *TrackPointStore) GetByTargetId(targetId string) ([]*model.TrackPoint, error) {
	ret := _m.Called(targetId)

	var r0 []*model.TrackPoint
	if rf, ok := ret.Get(0).(func(string) []*model.TrackPoint); ok {
		r0 = rf(targetId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TrackPoint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(targetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: trackPoint
func (_m *TrackPointStore) Save(trackPoint *model.TrackPoint) (*model.TrackPoint, error) {
	ret := _m.Called(trackPoint)

	var r0 *model.TrackPoint
	if rf, ok := ret.Get(0).(func(*model.TrackPoint) *model.TrackPoint); ok {
		r0 = rf(trackPoint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackPoint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.TrackPoint) error); ok {
		r1 = rf(trackPoint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
