// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	model "github.com/masterhung0112/go_server/model"
	mock "github.com/stretchr/testify/mock"
)

// ChannelStore is an autogenerated mock type for the ChannelStore type
type ChannelStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: id, allowFromCache
func (_m *ChannelStore) Get(id string, allowFromCache bool) (*model.Channel, error) {
	ret := _m.Called(id, allowFromCache)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string, bool) *model.Channel); ok {
		r0 = rf(id, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(id, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllChannelMembersForUser provides a mock function with given fields: userId, allowFromCache, includeDeleted
func (_m *ChannelStore) GetAllChannelMembersForUser(userId string, allowFromCache bool, includeDeleted bool) (map[string]string, *model.AppError) {
	ret := _m.Called(userId, allowFromCache, includeDeleted)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string, bool, bool) map[string]string); ok {
		r0 = rf(userId, allowFromCache, includeDeleted)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, bool, bool) *model.AppError); ok {
		r1 = rf(userId, allowFromCache, includeDeleted)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetForPost provides a mock function with given fields: postId
func (_m *ChannelStore) GetForPost(postId string) (*model.Channel, error) {
	ret := _m.Called(postId)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(string) *model.Channel); ok {
		r0 = rf(postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberForPost provides a mock function with given fields: postId, userId
func (_m *ChannelStore) GetMemberForPost(postId string, userId string) (*model.ChannelMember, *model.AppError) {
	ret := _m.Called(postId, userId)

	var r0 *model.ChannelMember
	if rf, ok := ret.Get(0).(func(string, string) *model.ChannelMember); ok {
		r0 = rf(postId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ChannelMember)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(postId, userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetSidebarCategory provides a mock function with given fields: categoryId
func (_m *ChannelStore) GetSidebarCategory(categoryId string) (*model.SidebarCategoryWithChannels, *model.AppError) {
	ret := _m.Called(categoryId)

	var r0 *model.SidebarCategoryWithChannels
	if rf, ok := ret.Get(0).(func(string) *model.SidebarCategoryWithChannels); ok {
		r0 = rf(categoryId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SidebarCategoryWithChannels)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(categoryId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Save provides a mock function with given fields: channel, maxChannelsPerTeam
func (_m *ChannelStore) Save(channel *model.Channel, maxChannelsPerTeam int64) (*model.Channel, error) {
	ret := _m.Called(channel, maxChannelsPerTeam)

	var r0 *model.Channel
	if rf, ok := ret.Get(0).(func(*model.Channel, int64) *model.Channel); ok {
		r0 = rf(channel, maxChannelsPerTeam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Channel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Channel, int64) error); ok {
		r1 = rf(channel, maxChannelsPerTeam)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
