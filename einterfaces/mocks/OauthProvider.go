// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `npm run task einterfaces_mocks`.

package mocks

import (
	io "io"

	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// OauthProvider is an autogenerated mock type for the OauthProvider type
type OauthProvider struct {
	mock.Mock
}

// GetSSOSettings provides a mock function with given fields: config, service
func (_m *OauthProvider) GetSSOSettings(config *model.Config, service string) (*model.SSOSettings, error) {
	ret := _m.Called(config, service)

	var r0 *model.SSOSettings
	if rf, ok := ret.Get(0).(func(*model.Config, string) *model.SSOSettings); ok {
		r0 = rf(config, service)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SSOSettings)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Config, string) error); ok {
		r1 = rf(config, service)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserFromIdToken provides a mock function with given fields: idToken
func (_m *OauthProvider) GetUserFromIdToken(idToken string) (*model.User, error) {
	ret := _m.Called(idToken)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(idToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(idToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserFromJson provides a mock function with given fields: data, tokenUser
func (_m *OauthProvider) GetUserFromJson(data io.Reader, tokenUser *model.User) (*model.User, error) {
	ret := _m.Called(data, tokenUser)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(io.Reader, *model.User) *model.User); ok {
		r0 = rf(data, tokenUser)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader, *model.User) error); ok {
		r1 = rf(data, tokenUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
