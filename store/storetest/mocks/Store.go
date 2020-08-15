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
