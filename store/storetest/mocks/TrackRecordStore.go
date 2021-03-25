// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `npm run task store_mocks`.

package mocks

import (
	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// TrackRecordStore is an autogenerated mock type for the TrackRecordStore type
type TrackRecordStore struct {
	mock.Mock
}

// End provides a mock function with given fields: trackRecordId
func (_m *TrackRecordStore) End(trackRecordId string) (*model.TrackRecord, error) {
	ret := _m.Called(trackRecordId)

	var r0 *model.TrackRecord
	if rf, ok := ret.Get(0).(func(string) *model.TrackRecord); ok {
		r0 = rf(trackRecordId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(trackRecordId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: trackRecordId
func (_m *TrackRecordStore) Get(trackRecordId string) (*model.TrackRecord, error) {
	ret := _m.Called(trackRecordId)

	var r0 *model.TrackRecord
	if rf, ok := ret.Get(0).(func(string) *model.TrackRecord); ok {
		r0 = rf(trackRecordId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(trackRecordId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: trackRecord
func (_m *TrackRecordStore) Save(trackRecord *model.TrackRecord) (*model.TrackRecord, error) {
	ret := _m.Called(trackRecord)

	var r0 *model.TrackRecord
	if rf, ok := ret.Get(0).(func(*model.TrackRecord) *model.TrackRecord); ok {
		r0 = rf(trackRecord)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.TrackRecord) error); ok {
		r1 = rf(trackRecord)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: trackRecordId
func (_m *TrackRecordStore) Start(trackRecordId string) (*model.TrackRecord, error) {
	ret := _m.Called(trackRecordId)

	var r0 *model.TrackRecord
	if rf, ok := ret.Get(0).(func(string) *model.TrackRecord); ok {
		r0 = rf(trackRecordId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(trackRecordId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: trackRecord
func (_m *TrackRecordStore) Update(trackRecord *model.TrackRecord) (*model.TrackRecord, error) {
	ret := _m.Called(trackRecord)

	var r0 *model.TrackRecord
	if rf, ok := ret.Get(0).(func(*model.TrackRecord) *model.TrackRecord); ok {
		r0 = rf(trackRecord)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TrackRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.TrackRecord) error); ok {
		r1 = rf(trackRecord)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
