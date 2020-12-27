// Code generated by mockery v1.0.0. DO NOT EDIT.

// 'Regenerate

package mocks

import (
	context "context"

	model "github.com/masterhung0112/hk_server/model"
	mock "github.com/stretchr/testify/mock"
)

// TeamStore is an autogenerated mock type for the TeamStore type
type TeamStore struct {
	mock.Mock
}

// AnalyticsGetTeamCountForScheme provides a mock function with given fields: schemeId
func (_m *TeamStore) AnalyticsGetTeamCountForScheme(schemeId string) (int64, error) {
	ret := _m.Called(schemeId)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(schemeId)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(schemeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsPrivateTeamCount provides a mock function with given fields:
func (_m *TeamStore) AnalyticsPrivateTeamCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsPublicTeamCount provides a mock function with given fields:
func (_m *TeamStore) AnalyticsPublicTeamCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AnalyticsTeamCount provides a mock function with given fields: includeDeleted
func (_m *TeamStore) AnalyticsTeamCount(includeDeleted bool) (int64, error) {
	ret := _m.Called(includeDeleted)

	var r0 int64
	if rf, ok := ret.Get(0).(func(bool) int64); ok {
		r0 = rf(includeDeleted)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(includeDeleted)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClearAllCustomRoleAssignments provides a mock function with given fields:
func (_m *TeamStore) ClearAllCustomRoleAssignments() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearCaches provides a mock function with given fields:
func (_m *TeamStore) ClearCaches() {
	_m.Called()
}

// Get provides a mock function with given fields: id
func (_m *TeamStore) Get(id string) (*model.Team, error) {
	ret := _m.Called(id)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
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

// GetActiveMemberCount provides a mock function with given fields: teamId, restrictions
func (_m *TeamStore) GetActiveMemberCount(teamId string, restrictions *model.ViewUsersRestrictions) (int64, error) {
	ret := _m.Called(teamId, restrictions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, *model.ViewUsersRestrictions) int64); ok {
		r0 = rf(teamId, restrictions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.ViewUsersRestrictions) error); ok {
		r1 = rf(teamId, restrictions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *TeamStore) GetAll() ([]*model.Team, error) {
	ret := _m.Called()

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func() []*model.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
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

// GetAllForExportAfter provides a mock function with given fields: limit, afterId
func (_m *TeamStore) GetAllForExportAfter(limit int, afterId string) ([]*model.TeamForExport, error) {
	ret := _m.Called(limit, afterId)

	var r0 []*model.TeamForExport
	if rf, ok := ret.Get(0).(func(int, string) []*model.TeamForExport); ok {
		r0 = rf(limit, afterId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamForExport)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, string) error); ok {
		r1 = rf(limit, afterId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPage provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllPage(offset int, limit int) ([]*model.Team, error) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPrivateTeamListing provides a mock function with given fields:
func (_m *TeamStore) GetAllPrivateTeamListing() ([]*model.Team, error) {
	ret := _m.Called()

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func() []*model.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
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

// GetAllPrivateTeamPageListing provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllPrivateTeamPageListing(offset int, limit int) ([]*model.Team, error) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPublicTeamPageListing provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllPublicTeamPageListing(offset int, limit int) ([]*model.Team, error) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllTeamListing provides a mock function with given fields:
func (_m *TeamStore) GetAllTeamListing() ([]*model.Team, error) {
	ret := _m.Called()

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func() []*model.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
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

// GetAllTeamPageListing provides a mock function with given fields: offset, limit
func (_m *TeamStore) GetAllTeamPageListing(offset int, limit int) ([]*model.Team, error) {
	ret := _m.Called(offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(int, int) []*model.Team); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByInviteId provides a mock function with given fields: inviteId
func (_m *TeamStore) GetByInviteId(inviteId string) (*model.Team, error) {
	ret := _m.Called(inviteId)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(inviteId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(inviteId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *TeamStore) GetByName(name string) (*model.Team, error) {
	ret := _m.Called(name)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(string) *model.Team); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
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

// GetByNames provides a mock function with given fields: name
func (_m *TeamStore) GetByNames(name []string) ([]*model.Team, error) {
	ret := _m.Called(name)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func([]string) []*model.Team); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChannelUnreadsForAllTeams provides a mock function with given fields: excludeTeamId, userId
func (_m *TeamStore) GetChannelUnreadsForAllTeams(excludeTeamId string, userId string) ([]*model.ChannelUnread, error) {
	ret := _m.Called(excludeTeamId, userId)

	var r0 []*model.ChannelUnread
	if rf, ok := ret.Get(0).(func(string, string) []*model.ChannelUnread); ok {
		r0 = rf(excludeTeamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ChannelUnread)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(excludeTeamId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChannelUnreadsForTeam provides a mock function with given fields: teamId, userId
func (_m *TeamStore) GetChannelUnreadsForTeam(teamId string, userId string) ([]*model.ChannelUnread, error) {
	ret := _m.Called(teamId, userId)

	var r0 []*model.ChannelUnread
	if rf, ok := ret.Get(0).(func(string, string) []*model.ChannelUnread); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ChannelUnread)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(teamId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMember provides a mock function with given fields: teamId, userId
func (_m *TeamStore) GetMember(teamId string, userId string) (*model.TeamMember, error) {
	ret := _m.Called(teamId, userId)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(string, string) *model.TeamMember); ok {
		r0 = rf(teamId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(teamId, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMembers provides a mock function with given fields: teamId, offset, limit, teamMembersGetOptions
func (_m *TeamStore) GetMembers(teamId string, offset int, limit int, teamMembersGetOptions *model.TeamMembersGetOptions) ([]*model.TeamMember, error) {
	ret := _m.Called(teamId, offset, limit, teamMembersGetOptions)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, int, int, *model.TeamMembersGetOptions) []*model.TeamMember); ok {
		r0 = rf(teamId, offset, limit, teamMembersGetOptions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int, *model.TeamMembersGetOptions) error); ok {
		r1 = rf(teamId, offset, limit, teamMembersGetOptions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMembersByIds provides a mock function with given fields: teamId, userIds, restrictions
func (_m *TeamStore) GetMembersByIds(teamId string, userIds []string, restrictions *model.ViewUsersRestrictions) ([]*model.TeamMember, error) {
	ret := _m.Called(teamId, userIds, restrictions)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, []string, *model.ViewUsersRestrictions) []*model.TeamMember); ok {
		r0 = rf(teamId, userIds, restrictions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string, *model.ViewUsersRestrictions) error); ok {
		r1 = rf(teamId, userIds, restrictions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeamMembersForExport provides a mock function with given fields: userId
func (_m *TeamStore) GetTeamMembersForExport(userId string) ([]*model.TeamMemberForExport, error) {
	ret := _m.Called(userId)

	var r0 []*model.TeamMemberForExport
	if rf, ok := ret.Get(0).(func(string) []*model.TeamMemberForExport); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMemberForExport)
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

// GetTeamsByScheme provides a mock function with given fields: schemeId, offset, limit
func (_m *TeamStore) GetTeamsByScheme(schemeId string, offset int, limit int) ([]*model.Team, error) {
	ret := _m.Called(schemeId, offset, limit)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.Team); ok {
		r0 = rf(schemeId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(schemeId, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeamsByUserId provides a mock function with given fields: userId
func (_m *TeamStore) GetTeamsByUserId(userId string) ([]*model.Team, error) {
	ret := _m.Called(userId)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
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

// GetTeamsForUser provides a mock function with given fields: ctx, userId
func (_m *TeamStore) GetTeamsForUser(ctx context.Context, userId string) ([]*model.TeamMember, error) {
	ret := _m.Called(ctx, userId)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.TeamMember); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTeamsForUserWithPagination provides a mock function with given fields: userId, page, perPage
func (_m *TeamStore) GetTeamsForUserWithPagination(userId string, page int, perPage int) ([]*model.TeamMember, error) {
	ret := _m.Called(userId, page, perPage)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.TeamMember); ok {
		r0 = rf(userId, page, perPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(userId, page, perPage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalMemberCount provides a mock function with given fields: teamId, restrictions
func (_m *TeamStore) GetTotalMemberCount(teamId string, restrictions *model.ViewUsersRestrictions) (int64, error) {
	ret := _m.Called(teamId, restrictions)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, *model.ViewUsersRestrictions) int64); ok {
		r0 = rf(teamId, restrictions)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.ViewUsersRestrictions) error); ok {
		r1 = rf(teamId, restrictions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserTeamIds provides a mock function with given fields: userId, allowFromCache
func (_m *TeamStore) GetUserTeamIds(userId string, allowFromCache bool) ([]string, error) {
	ret := _m.Called(userId, allowFromCache)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, bool) []string); ok {
		r0 = rf(userId, allowFromCache)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(userId, allowFromCache)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GroupSyncedTeamCount provides a mock function with given fields:
func (_m *TeamStore) GroupSyncedTeamCount() (int64, error) {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvalidateAllTeamIdsForUser provides a mock function with given fields: userId
func (_m *TeamStore) InvalidateAllTeamIdsForUser(userId string) {
	_m.Called(userId)
}

// MigrateTeamMembers provides a mock function with given fields: fromTeamId, fromUserId
func (_m *TeamStore) MigrateTeamMembers(fromTeamId string, fromUserId string) (map[string]string, error) {
	ret := _m.Called(fromTeamId, fromUserId)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(string, string) map[string]string); ok {
		r0 = rf(fromTeamId, fromUserId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(fromTeamId, fromUserId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDelete provides a mock function with given fields: teamId
func (_m *TeamStore) PermanentDelete(teamId string) error {
	ret := _m.Called(teamId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(teamId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAllMembersByTeam provides a mock function with given fields: teamId
func (_m *TeamStore) RemoveAllMembersByTeam(teamId string) error {
	ret := _m.Called(teamId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(teamId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAllMembersByUser provides a mock function with given fields: userId
func (_m *TeamStore) RemoveAllMembersByUser(userId string) error {
	ret := _m.Called(userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveMember provides a mock function with given fields: teamId, userId
func (_m *TeamStore) RemoveMember(teamId string, userId string) error {
	ret := _m.Called(teamId, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(teamId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveMembers provides a mock function with given fields: teamId, userIds
func (_m *TeamStore) RemoveMembers(teamId string, userIds []string) error {
	ret := _m.Called(teamId, userIds)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(teamId, userIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResetAllTeamSchemes provides a mock function with given fields:
func (_m *TeamStore) ResetAllTeamSchemes() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: team
func (_m *TeamStore) Save(team *model.Team) (*model.Team, error) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Team) error); ok {
		r1 = rf(team)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMember provides a mock function with given fields: member, maxUsersPerTeam
func (_m *TeamStore) SaveMember(member *model.TeamMember, maxUsersPerTeam int) (*model.TeamMember, error) {
	ret := _m.Called(member, maxUsersPerTeam)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(*model.TeamMember, int) *model.TeamMember); ok {
		r0 = rf(member, maxUsersPerTeam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.TeamMember, int) error); ok {
		r1 = rf(member, maxUsersPerTeam)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMultipleMembers provides a mock function with given fields: members, maxUsersPerTeam
func (_m *TeamStore) SaveMultipleMembers(members []*model.TeamMember, maxUsersPerTeam int) ([]*model.TeamMember, error) {
	ret := _m.Called(members, maxUsersPerTeam)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func([]*model.TeamMember, int) []*model.TeamMember); ok {
		r0 = rf(members, maxUsersPerTeam)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.TeamMember, int) error); ok {
		r1 = rf(members, maxUsersPerTeam)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchAll provides a mock function with given fields: term, opts
func (_m *TeamStore) SearchAll(term string, opts *model.TeamSearch) ([]*model.Team, error) {
	ret := _m.Called(term, opts)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string, *model.TeamSearch) []*model.Team); ok {
		r0 = rf(term, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *model.TeamSearch) error); ok {
		r1 = rf(term, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchAllPaged provides a mock function with given fields: term, opts
func (_m *TeamStore) SearchAllPaged(term string, opts *model.TeamSearch) ([]*model.Team, int64, error) {
	ret := _m.Called(term, opts)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string, *model.TeamSearch) []*model.Team); ok {
		r0 = rf(term, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(string, *model.TeamSearch) int64); ok {
		r1 = rf(term, opts)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, *model.TeamSearch) error); ok {
		r2 = rf(term, opts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SearchOpen provides a mock function with given fields: term
func (_m *TeamStore) SearchOpen(term string) ([]*model.Team, error) {
	ret := _m.Called(term)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(term)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(term)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchPrivate provides a mock function with given fields: term
func (_m *TeamStore) SearchPrivate(term string) ([]*model.Team, error) {
	ret := _m.Called(term)

	var r0 []*model.Team
	if rf, ok := ret.Get(0).(func(string) []*model.Team); ok {
		r0 = rf(term)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(term)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: team
func (_m *TeamStore) Update(team *model.Team) (*model.Team, error) {
	ret := _m.Called(team)

	var r0 *model.Team
	if rf, ok := ret.Get(0).(func(*model.Team) *model.Team); ok {
		r0 = rf(team)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Team)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Team) error); ok {
		r1 = rf(team)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateLastTeamIconUpdate provides a mock function with given fields: teamId, curTime
func (_m *TeamStore) UpdateLastTeamIconUpdate(teamId string, curTime int64) error {
	ret := _m.Called(teamId, curTime)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(teamId, curTime)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMember provides a mock function with given fields: member
func (_m *TeamStore) UpdateMember(member *model.TeamMember) (*model.TeamMember, error) {
	ret := _m.Called(member)

	var r0 *model.TeamMember
	if rf, ok := ret.Get(0).(func(*model.TeamMember) *model.TeamMember); ok {
		r0 = rf(member)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.TeamMember) error); ok {
		r1 = rf(member)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMembersRole provides a mock function with given fields: teamID, userIDs
func (_m *TeamStore) UpdateMembersRole(teamID string, userIDs []string) error {
	ret := _m.Called(teamID, userIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []string) error); ok {
		r0 = rf(teamID, userIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMultipleMembers provides a mock function with given fields: members
func (_m *TeamStore) UpdateMultipleMembers(members []*model.TeamMember) ([]*model.TeamMember, error) {
	ret := _m.Called(members)

	var r0 []*model.TeamMember
	if rf, ok := ret.Get(0).(func([]*model.TeamMember) []*model.TeamMember); ok {
		r0 = rf(members)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TeamMember)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.TeamMember) error); ok {
		r1 = rf(members)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserBelongsToTeams provides a mock function with given fields: userId, teamIds
func (_m *TeamStore) UserBelongsToTeams(userId string, teamIds []string) (bool, error) {
	ret := _m.Called(userId, teamIds)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, []string) bool); ok {
		r0 = rf(userId, teamIds)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []string) error); ok {
		r1 = rf(userId, teamIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
