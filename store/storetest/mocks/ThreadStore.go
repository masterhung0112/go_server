// Code generated by mockery v1.0.0. DO NOT EDIT.

// 'Regenerate

package mocks

import (
	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// ThreadStore is an autogenerated mock type for the ThreadStore type
type ThreadStore struct {
	mock.Mock
}

// CollectThreadsWithNewerReplies provides a mock function with given fields: userId, channelIds, timestamp
func (_m *ThreadStore) CollectThreadsWithNewerReplies(userId string, channelIds []string, timestamp int64) ([]string, error) {
	ret := _m.Called(userId, channelIds, timestamp)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, []string, int64) []string); ok {
		r0 = rf(userId, channelIds, timestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string, int64) error); ok {
		r1 = rf(userId, channelIds, timestamp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *ThreadStore) Get(id string) (*model.Thread, error) {
	ret := _m.Called(id)

	var r0 *model.Thread
	if rf, ok := ret.Get(0).(func(string) *model.Thread); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Thread)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: thread
func (_m *ThreadStore) Update(thread *model.Thread) (*model.Thread, error) {
	ret := _m.Called(thread)

	var r0 *model.Thread
	if rf, ok := ret.Get(0).(func(*model.Thread) *model.Thread); ok {
		r0 = rf(thread)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Thread)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Thread) error); ok {
		r1 = rf(thread)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUnreadsByChannel provides a mock function with given fields: userId, changedThreads, timestamp, updateViewedTimestamp
func (_m *ThreadStore) UpdateUnreadsByChannel(userId string, changedThreads []string, timestamp int64, updateViewedTimestamp bool) error {
	ret := _m.Called(userId, changedThreads, timestamp, updateViewedTimestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string, int64, bool) error); ok {
		r0 = rf(userId, changedThreads, timestamp, updateViewedTimestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
