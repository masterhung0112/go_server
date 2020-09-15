// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	store "github.com/masterhung0112/go_server/store"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Channel provides a mock function with given fields:
func (_m *Store) Channel() store.ChannelStore {
	ret := _m.Called()

	var r0 store.ChannelStore
	if rf, ok := ret.Get(0).(func() store.ChannelStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.ChannelStore)
		}
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Store) Close() {
	_m.Called()
}

// DropAllTables provides a mock function with given fields:
func (_m *Store) DropAllTables() {
	_m.Called()
}

// MarkSystemRanUnitTests provides a mock function with given fields:
func (_m *Store) MarkSystemRanUnitTests() {
	_m.Called()
}

// Preference provides a mock function with given fields:
func (_m *Store) Preference() store.PreferenceStore {
	ret := _m.Called()

	var r0 store.PreferenceStore
	if rf, ok := ret.Get(0).(func() store.PreferenceStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.PreferenceStore)
		}
	}

	return r0
}

// Role provides a mock function with given fields:
func (_m *Store) Role() store.RoleStore {
	ret := _m.Called()

	var r0 store.RoleStore
	if rf, ok := ret.Get(0).(func() store.RoleStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.RoleStore)
		}
	}

	return r0
}

// Scheme provides a mock function with given fields:
func (_m *Store) Scheme() store.SchemeStore {
	ret := _m.Called()

	var r0 store.SchemeStore
	if rf, ok := ret.Get(0).(func() store.SchemeStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.SchemeStore)
		}
	}

	return r0
}

// Session provides a mock function with given fields:
func (_m *Store) Session() store.SessionStore {
	ret := _m.Called()

	var r0 store.SessionStore
	if rf, ok := ret.Get(0).(func() store.SessionStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.SessionStore)
		}
	}

	return r0
}

// System provides a mock function with given fields:
func (_m *Store) System() store.SystemStore {
	ret := _m.Called()

	var r0 store.SystemStore
	if rf, ok := ret.Get(0).(func() store.SystemStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.SystemStore)
		}
	}

	return r0
}

// Team provides a mock function with given fields:
func (_m *Store) Team() store.TeamStore {
	ret := _m.Called()

	var r0 store.TeamStore
	if rf, ok := ret.Get(0).(func() store.TeamStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.TeamStore)
		}
	}

	return r0
}

// Token provides a mock function with given fields:
func (_m *Store) Token() store.TokenStore {
	ret := _m.Called()

	var r0 store.TokenStore
	if rf, ok := ret.Get(0).(func() store.TokenStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.TokenStore)
		}
	}

	return r0
}

// User provides a mock function with given fields:
func (_m *Store) User() store.UserStore {
	ret := _m.Called()

	var r0 store.UserStore
	if rf, ok := ret.Get(0).(func() store.UserStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.UserStore)
		}
	}

	return r0
}

// UserAccessToken provides a mock function with given fields:
func (_m *Store) UserAccessToken() store.UserAccessTokenStore {
	ret := _m.Called()

	var r0 store.UserAccessTokenStore
	if rf, ok := ret.Get(0).(func() store.UserAccessTokenStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(store.UserAccessTokenStore)
		}
	}

	return r0
}
