// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/masterhung0112/go_server/model"
	mock "github.com/stretchr/testify/mock"
)

// PreferenceStore is an autogenerated mock type for the PreferenceStore type
type PreferenceStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: userId, category, name
func (_m *PreferenceStore) Get(userId string, category string, name string) (*model.Preference, error) {
	ret := _m.Called(userId, category, name)

	var r0 *model.Preference
	if rf, ok := ret.Get(0).(func(string, string, string) *model.Preference); ok {
		r0 = rf(userId, category, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Preference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(userId, category, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: preferences
func (_m *PreferenceStore) Save(preferences *model.Preferences) error {
	ret := _m.Called(preferences)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Preferences) error); ok {
		r0 = rf(preferences)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
