// Code generated by mockery v1.0.0. DO NOT EDIT.

// 'Regenerate

package mocks

import (
	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// UploadSessionStore is an autogenerated mock type for the UploadSessionStore type
type UploadSessionStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *UploadSessionStore) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *UploadSessionStore) Get(id string) (*model.UploadSession, error) {
	ret := _m.Called(id)

	var r0 *model.UploadSession
	if rf, ok := ret.Get(0).(func(string) *model.UploadSession); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UploadSession)
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

// GetForUser provides a mock function with given fields: userId
func (_m *UploadSessionStore) GetForUser(userId string) ([]*model.UploadSession, error) {
	ret := _m.Called(userId)

	var r0 []*model.UploadSession
	if rf, ok := ret.Get(0).(func(string) []*model.UploadSession); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.UploadSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: session
func (_m *UploadSessionStore) Save(session *model.UploadSession) (*model.UploadSession, error) {
	ret := _m.Called(session)

	var r0 *model.UploadSession
	if rf, ok := ret.Get(0).(func(*model.UploadSession) *model.UploadSession); ok {
		r0 = rf(session)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UploadSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.UploadSession) error); ok {
		r1 = rf(session)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: session
func (_m *UploadSessionStore) Update(session *model.UploadSession) error {
	ret := _m.Called(session)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.UploadSession) error); ok {
		r0 = rf(session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
