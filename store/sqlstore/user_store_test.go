package sqlstore

import (
	"github.com/masterhung0112/hk_server/model"
	"github.com/masterhung0112/hk_server/store"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"

	"github.com/masterhung0112/hk_server/store/storetest"
)

const (
	DAY_MILLISECONDS   = 24 * 60 * 60 * 1000
	MONTH_MILLISECONDS = 31 * DAY_MILLISECONDS
)

// type UserStoreTestSuite struct {
// 	suite.Suite
// 	StoreTestSuite
// }

func sanitized(user *model.User) *model.User {
	clonedUser := user.DeepCopy()
	clonedUser.Sanitize(map[string]bool{})

	return clonedUser
}

type UserStoreTS struct {
	suite.Suite
	StoreTestSuite
}

func TestUserStoreTS(t *testing.T) {
	StoreTestSuiteWithSqlSupplier(t, &UserStoreTS{}, func(t *testing.T, testSuite StoreTestBaseSuite) {
		suite.Run(t, testSuite)
	})
}

func (s *UserStoreTS) cleanupStatusStore() {
	_, execerr := s.SqlStore().GetMaster().ExecNoTimeout(` DELETE FROM Status `)
	s.Require().Nil(execerr)
}

func (s *UserStoreTS) TestCount() {
	teamId := model.NewId()
	channelId := model.NewId()
	regularUser := &model.User{}
	regularUser.Email = MakeEmail()
	regularUser.Roles = model.SYSTEM_USER_ROLE_ID
	_, err := s.Store().User().Save(regularUser)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(regularUser.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: regularUser.Id, SchemeAdmin: false, SchemeUser: true}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{UserId: regularUser.Id, ChannelId: channelId, SchemeAdmin: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	s.Require().Nil(nErr)

	guestUser := &model.User{}
	guestUser.Email = MakeEmail()
	guestUser.Roles = model.SYSTEM_GUEST_ROLE_ID
	_, err = s.Store().User().Save(guestUser)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(guestUser.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: guestUser.Id, SchemeAdmin: false, SchemeUser: false, SchemeGuest: true}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{UserId: guestUser.Id, ChannelId: channelId, SchemeAdmin: false, SchemeUser: false, SchemeGuest: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	s.Require().Nil(nErr)

	teamAdmin := &model.User{}
	teamAdmin.Email = MakeEmail()
	teamAdmin.Roles = model.SYSTEM_USER_ROLE_ID
	_, err = s.Store().User().Save(teamAdmin)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(teamAdmin.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: teamAdmin.Id, SchemeAdmin: true, SchemeUser: true}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{UserId: teamAdmin.Id, ChannelId: channelId, SchemeAdmin: true, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	s.Require().Nil(nErr)

	sysAdmin := &model.User{}
	sysAdmin.Email = MakeEmail()
	sysAdmin.Roles = model.SYSTEM_ADMIN_ROLE_ID + " " + model.SYSTEM_USER_ROLE_ID
	_, err = s.Store().User().Save(sysAdmin)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(sysAdmin.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: sysAdmin.Id, SchemeAdmin: false, SchemeUser: true}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{UserId: sysAdmin.Id, ChannelId: channelId, SchemeAdmin: true, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
	s.Require().Nil(nErr)

	// Deleted
	deletedUser := &model.User{}
	deletedUser.Email = MakeEmail()
	deletedUser.DeleteAt = model.GetMillis()
	_, err = s.Store().User().Save(deletedUser)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(deletedUser.Id)) }()

	// Bot
	botUser, err := s.Store().User().Save(&model.User{
		Email: MakeEmail(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(botUser.Id)) }()
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   botUser.Id,
		Username: botUser.Username,
		OwnerId:  regularUser.Id,
	})
	s.Require().Nil(nErr)
	botUser.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(botUser.Id)) }()

	testCases := []struct {
		Description string
		Options     model.UserCountOptions
		Expected    int64
	}{
		{
			"No bot accounts no deleted accounts and no team id",
			model.UserCountOptions{
				IncludeBotAccounts: false,
				IncludeDeleted:     false,
				TeamId:             "",
			},
			4,
		},
		{
			"Include bot accounts no deleted accounts and no team id",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     false,
				TeamId:             "",
			},
			5,
		},
		{
			"Include delete accounts no bots and no team id",
			model.UserCountOptions{
				IncludeBotAccounts: false,
				IncludeDeleted:     true,
				TeamId:             "",
			},
			5,
		},
		{
			"Include bot accounts and deleted accounts and no team id",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     true,
				TeamId:             "",
			},
			6,
		},
		{
			"Include bot accounts, deleted accounts, exclude regular users with no team id",
			model.UserCountOptions{
				IncludeBotAccounts:  true,
				IncludeDeleted:      true,
				ExcludeRegularUsers: true,
				TeamId:              "",
			},
			1,
		},
		{
			"Include bot accounts and deleted accounts with existing team id",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     true,
				TeamId:             teamId,
			},
			4,
		},
		{
			"Include bot accounts and deleted accounts with fake team id",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     true,
				TeamId:             model.NewId(),
			},
			0,
		},
		{
			"Include bot accounts and deleted accounts with existing team id and view restrictions allowing team",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     true,
				TeamId:             teamId,
				ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{teamId}},
			},
			4,
		},
		{
			"Include bot accounts and deleted accounts with existing team id and view restrictions not allowing current team",
			model.UserCountOptions{
				IncludeBotAccounts: true,
				IncludeDeleted:     true,
				TeamId:             teamId,
				ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{model.NewId()}},
			},
			0,
		},
		{
			"Filter by system admins only",
			model.UserCountOptions{
				TeamId: teamId,
				Roles:  []string{model.SYSTEM_ADMIN_ROLE_ID},
			},
			1,
		},
		{
			"Filter by system users only",
			model.UserCountOptions{
				TeamId: teamId,
				Roles:  []string{model.SYSTEM_USER_ROLE_ID},
			},
			2,
		},
		{
			"Filter by system guests only",
			model.UserCountOptions{
				TeamId: teamId,
				Roles:  []string{model.SYSTEM_GUEST_ROLE_ID},
			},
			1,
		},
		{
			"Filter by system admins and system users",
			model.UserCountOptions{
				TeamId: teamId,
				Roles:  []string{model.SYSTEM_ADMIN_ROLE_ID, model.SYSTEM_USER_ROLE_ID},
			},
			3,
		},
		{
			"Filter by system admins, system user and system guests",
			model.UserCountOptions{
				TeamId: teamId,
				Roles:  []string{model.SYSTEM_ADMIN_ROLE_ID, model.SYSTEM_USER_ROLE_ID, model.SYSTEM_GUEST_ROLE_ID},
			},
			4,
		},
		{
			"Filter by team admins",
			model.UserCountOptions{
				TeamId:    teamId,
				TeamRoles: []string{model.TEAM_ADMIN_ROLE_ID},
			},
			1,
		},
		{
			"Filter by team members",
			model.UserCountOptions{
				TeamId:    teamId,
				TeamRoles: []string{model.TEAM_USER_ROLE_ID},
			},
			1,
		},
		{
			"Filter by team guests",
			model.UserCountOptions{
				TeamId:    teamId,
				TeamRoles: []string{model.TEAM_GUEST_ROLE_ID},
			},
			1,
		},
		{
			"Filter by team guests and any system role",
			model.UserCountOptions{
				TeamId:    teamId,
				TeamRoles: []string{model.TEAM_GUEST_ROLE_ID},
				Roles:     []string{model.SYSTEM_ADMIN_ROLE_ID},
			},
			2,
		},
		{
			"Filter by channel members",
			model.UserCountOptions{
				ChannelId:    channelId,
				ChannelRoles: []string{model.CHANNEL_USER_ROLE_ID},
			},
			1,
		},
		{
			"Filter by channel members and system admins",
			model.UserCountOptions{
				ChannelId:    channelId,
				Roles:        []string{model.SYSTEM_ADMIN_ROLE_ID},
				ChannelRoles: []string{model.CHANNEL_USER_ROLE_ID},
			},
			2,
		},
		{
			"Filter by channel members and system admins and channel admins",
			model.UserCountOptions{
				ChannelId:    channelId,
				Roles:        []string{model.SYSTEM_ADMIN_ROLE_ID},
				ChannelRoles: []string{model.CHANNEL_USER_ROLE_ID, model.CHANNEL_ADMIN_ROLE_ID},
			},
			3,
		},
		{
			"Filter by channel guests",
			model.UserCountOptions{
				ChannelId:    channelId,
				ChannelRoles: []string{model.CHANNEL_GUEST_ROLE_ID},
			},
			1,
		},
		{
			"Filter by channel guests and any system role",
			model.UserCountOptions{
				ChannelId:    channelId,
				ChannelRoles: []string{model.CHANNEL_GUEST_ROLE_ID},
				Roles:        []string{model.SYSTEM_ADMIN_ROLE_ID},
			},
			2,
		},
	}
	for _, testCase := range testCases {
		s.T().Run(testCase.Description, func(t *testing.T) {
			count, err := s.Store().User().Count(testCase.Options)
			s.Require().Nil(err)
			s.Require().Equal(testCase.Expected, count)
		})
	}
}

func (s *UserStoreTS) TestSave() {
	teamId := model.NewId()
	maxUsersPerTeam := 50

	u1 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}

	_, err := s.Store().User().Save(&u1)
	s.Require().Nil(err, "couldn't save user")

	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam)
	s.Require().Nil(nErr)

	_, err = s.Store().User().Save(&u1)
	s.Require().NotNil(err, "shouldn't be able to update user from save")

	u2 := model.User{
		Email:    u1.Email,
		Username: model.NewId(),
	}
	_, err = s.Store().User().Save(&u2)
	s.Require().NotNil(err, "should be unique email")

	u2.Email = MakeEmail()
	u2.Username = u1.Username
	_, err = s.Store().User().Save(&u2)
	s.Require().NotNil(err, "should be unique username")

	u2.Username = ""
	_, err = s.Store().User().Save(&u2)
	s.Require().NotNil(err, "should be unique username")

	for i := 0; i < 49; i++ {
		u := model.User{
			Email:    MakeEmail(),
			Username: model.NewId(),
		}
		_, err = s.Store().User().Save(&u)
		s.Require().Nil(err, "couldn't save item")

		defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u.Id)) }()

		_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u.Id}, maxUsersPerTeam)
		s.Require().Nil(nErr)
	}

	u2.Id = ""
	u2.Email = MakeEmail()
	u2.Username = model.NewId()
	_, err = s.Store().User().Save(&u2)
	s.Require().Nil(err, "couldn't save item")

	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam)
	s.Require().NotNil(nErr, "should be the limit")
}

func (s *UserStoreTS) TestUpdate() {
	u1 := &model.User{
		Email: MakeEmail(),
	}
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2 := &model.User{
		Email:       MakeEmail(),
		AuthService: "ldap",
	}
	_, err = s.Store().User().Save(u2)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	_, err = s.Store().User().Update(u1, false)
	s.Require().Nil(err)

	missing := &model.User{}
	_, err = s.Store().User().Update(missing, false)
	s.Require().NotNil(err, "Update should have failed because of missing key")

	newId := &model.User{
		Id: model.NewId(),
	}
	_, err = s.Store().User().Update(newId, false)
	s.Require().NotNil(err, "Update should have failed because id change")

	u2.Email = MakeEmail()
	_, err = s.Store().User().Update(u2, false)
	s.Require().NotNil(err, "Update should have failed because you can't modify AD/LDAP fields")

	u3 := &model.User{
		Email:       MakeEmail(),
		AuthService: "gitlab",
	}
	oldEmail := u3.Email
	_, err = s.Store().User().Save(u3)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u3.Id}, -1)
	s.Require().Nil(nErr)

	u3.Email = MakeEmail()
	userUpdate, err := s.Store().User().Update(u3, false)
	s.Require().Nil(err, "Update should not have failed")
	s.Assert().Equal(oldEmail, userUpdate.New.Email, "Email should not have been updated as the update is not trusted")

	u3.Email = MakeEmail()
	userUpdate, err = s.Store().User().Update(u3, true)
	s.Require().Nil(err, "Update should not have failed")
	s.Assert().NotEqual(oldEmail, userUpdate.New.Email, "Email should have been updated as the update is trusted")

	err = s.Store().User().UpdateLastPictureUpdate(u1.Id)
	s.Require().Nil(err, "Update should not have failed")
}

func (s *UserStoreTS) TestResetLastPictureUpdate() {

	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	err = s.Store().User().UpdateLastPictureUpdate(u1.Id)
	s.Require().Nil(err)

	user, err := s.Store().User().Get(u1.Id)
	s.Require().Nil(err)

	s.Assert().NotZero(user.LastPictureUpdate)
	s.Assert().NotZero(user.UpdateAt)

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond)

	err = s.Store().User().ResetLastPictureUpdate(u1.Id)
	s.Require().Nil(err)

	s.Store().User().InvalidateProfileCacheForUser(u1.Id)

	user2, err := s.Store().User().Get(u1.Id)
	s.Require().Nil(err)

	s.Assert().True(user2.UpdateAt > user.UpdateAt)
	s.Assert().Zero(user2.LastPictureUpdate)
}

func (s *UserStoreTS) TestUpdatePassword() {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	hashedPassword := model.HashPassword("newpwd")

	err = s.Store().User().UpdatePassword(u1.Id, hashedPassword)
	s.Require().Nil(err)

	user, err := s.Store().User().GetByEmail(u1.Email)
	s.Require().Nil(err)
	s.Require().Equal(user.Password, hashedPassword, "Password was not updated correctly")
}

func (s *UserStoreTS) TestUpdateUpdateAt() {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	_, err = s.Store().User().UpdateUpdateAt(u1.Id)
	s.Require().Nil(err)

	user, err := s.Store().User().Get(u1.Id)
	s.Require().Nil(err)
	s.Require().Less(u1.UpdateAt, user.UpdateAt, "UpdateAt not updated correctly")
}

func (s *UserStoreTS) TestUpdateAuthData() {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	service := "someservice"
	authData := model.NewId()

	_, err = s.Store().User().UpdateAuthData(u1.Id, service, &authData, "", true)
	s.Require().Nil(err)

	user, err := s.Store().User().GetByEmail(u1.Email)
	s.Require().Nil(err)
	s.Require().Equal(service, user.AuthService, "AuthService was not updated correctly")
	s.Require().Equal(authData, *user.AuthData, "AuthData was not updated correctly")
	s.Require().Equal("", user.Password, "Password was not cleared properly")
}

func (s *UserStoreTS) TestUpdateMfaSecret() {
	u1 := model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(&u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	err = s.Store().User().UpdateMfaSecret(u1.Id, "12345")
	s.Require().Nil(err)

	// should pass, no update will occur though
	err = s.Store().User().UpdateMfaSecret("junk", "12345")
	s.Require().Nil(err)
}

func (s *UserStoreTS) TestUpdateMfaActive() {
	u1 := model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(&u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	time.Sleep(time.Millisecond)

	err = s.Store().User().UpdateMfaActive(u1.Id, true)
	s.Require().Nil(err)

	err = s.Store().User().UpdateMfaActive(u1.Id, false)
	s.Require().Nil(err)

	// should pass, no update will occur though
	err = s.Store().User().UpdateMfaActive("junk", true)
	s.Require().Nil(err)
}

func (s *UserStoreTS) TestGetProfilesInChannel() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u4.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1)
	s.Require().Nil(nErr)

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, nErr := s.Store().Channel().Save(ch1, -1)
	s.Require().Nil(nErr)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, nErr := s.Store().Channel().Save(ch2, -1)
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u4.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	u4.DeleteAt = 1
	_, err = s.Store().User().Update(u4, true)
	s.Require().Nil(err)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	s.T().Run("get all users in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u1), sanitized(u2), sanitized(u3), sanitized(u4)}, users)
	})

	s.T().Run("get only active users in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
			Active:      true,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, users)
	})

	s.T().Run("get inactive users in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
			Inactive:    true,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u4)}, users)
	})

	s.T().Run("get in channel 1, offset 1, limit 2", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        1,
			PerPage:     1,
		})
		s.Require().Nil(err)
		users_p2, err2 := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        2,
			PerPage:     1,
		})
		s.Require().Nil(err2)
		users = append(users, users_p2...)
		s.Assert().Equal([]*model.User{sanitized(u2), sanitized(u3)}, users)
	})

	s.T().Run("get in channel 2, offset 0, limit 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c2.Id,
			Page:        0,
			PerPage:     1,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u1)}, users)
	})
}

func (s *UserStoreTS) TestProfilesInChannelByStatus() {

	s.cleanupStatusStore()

	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u4.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1)
	s.Require().Nil(nErr)

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, nErr := s.Store().Channel().Save(ch1, -1)
	s.Require().Nil(nErr)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, nErr := s.Store().Channel().Save(ch2, -1)
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u4.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)

	u4.DeleteAt = 1
	_, err = s.Store().User().Update(u4, true)
	s.Require().Nil(err)

	_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	s.Require().Nil(nErr)
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{
		UserId: u1.Id,
		Status: model.STATUS_DND,
	}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{
		UserId: u2.Id,
		Status: model.STATUS_AWAY,
	}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{
		UserId: u3.Id,
		Status: model.STATUS_ONLINE,
	}))

	s.T().Run("get all users in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u1), sanitized(u2), sanitized(u3), sanitized(u4)}, users)
	})

	s.T().Run("get active in channel 1 by status, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannelByStatus(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
			Active:      true,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u3), sanitized(u2), sanitized(u1)}, users)
	})

	s.T().Run("get inactive users in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannel(&model.UserGetOptions{
			InChannelId: c1.Id,
			Page:        0,
			PerPage:     100,
			Inactive:    true,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u4)}, users)
	})

	s.T().Run("get in channel 2 by status, offset 0, limit 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesInChannelByStatus(&model.UserGetOptions{
			InChannelId: c2.Id,
			Page:        0,
			PerPage:     1,
		})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u1)}, users)
	})
}

func (s *UserStoreTS) TestGetProfilesWithoutTeam() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
		DeleteAt: 1,
		Roles:    "system_admin",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get, page 0, per_page 100", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u2), sanitized(u3)}, users)
	})

	s.T().Run("get, page 1, per_page 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 1, PerPage: 1})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u3)}, users)
	})

	s.T().Run("get, page 2, per_page 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 2, PerPage: 1})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{}, users)
	})

	s.T().Run("get, page 0, per_page 100, inactive", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100, Inactive: true})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u3)}, users)
	})

	s.T().Run("get, page 0, per_page 100, role", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100, Role: "system_admin"})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{sanitized(u3)}, users)
	})
}

func (s *UserStoreTS) TestGetProfilesByUsernames() {
	teamId := model.NewId()
	team2Id := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: team2Id, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get by u1 and u2 usernames, team id 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesByUsernames([]string{u1.Username, u2.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u1, u2}, users)
	})

	s.T().Run("get by u1 username, team id 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesByUsernames([]string{u1.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u1}, users)
	})

	s.T().Run("get by u1 and u3 usernames, no team id", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, nil)
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u1, u3}, users)
	})

	s.T().Run("get by u1 and u3 usernames, team id 1", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u1}, users)
	})

	s.T().Run("get by u1 and u3 usernames, team id 2", func(t *testing.T) {
		users, err := s.Store().User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, &model.ViewUsersRestrictions{Teams: []string{team2Id}})
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u3}, users)
	})
}

type UserStoreGetAllProfilesTS struct {
	suite.Suite
	StoreTestSuite

	u1 *model.User
	u2 *model.User
	u3 *model.User
	u4 *model.User
	u5 *model.User
	u6 *model.User
	u7 *model.User
}

func TestUserStoreGetAllProfilesTS(t *testing.T) {
	StoreTestSuiteWithSqlSupplier(t, &UserStoreGetAllProfilesTS{}, func(t *testing.T, testSuite StoreTestBaseSuite) {
		suite.Run(t, testSuite)
	})
}

func (s *UserStoreGetAllProfilesTS) SetupSuite() {
	// Clean all user first
	users, err := s.Store().User().GetAll()
	s.Require().Nil(err, "failed cleaning up test users")

	for _, u := range users {
		err := s.Store().User().PermanentDelete(u.Id)
		s.Require().Nil(err, "failed cleaning up test user %s", u.Username)
	}

	u1, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	s.u1 = u1

	u2, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	s.u2 = u2

	u3, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	s.u3 = u3

	_, nErr := s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true

	u4, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u4" + model.NewId(),
		Roles:    "system_user some-other-role",
	})
	s.Require().Nil(err)
	s.u4 = u4

	u5, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u5" + model.NewId(),
		Roles:    "system_admin",
	})
	s.Require().Nil(err)
	s.u5 = u5

	u6, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u6" + model.NewId(),
		DeleteAt: model.GetMillis(),
		Roles:    "system_admin",
	})
	s.Require().Nil(err)
	s.u6 = u6

	u7, err := s.Store().User().Save(&model.User{
		Email:    storetest.MakeEmail(),
		Username: "u7" + model.NewId(),
		DeleteAt: model.GetMillis(),
	})
	s.Require().Nil(err)
	s.u7 = u7
}

func (s *UserStoreGetAllProfilesTS) TearDownSuite() {
	s.Require().Nil(s.Store().User().PermanentDelete(s.u1.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u2.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u3.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u4.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u5.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u6.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u7.Id))
}

func (s *UserStoreGetAllProfilesTS) TestGetOffset0Limit100() {
	options := &model.UserGetOptions{Page: 0, PerPage: 100}
	actual, err := s.Store().User().GetAllProfiles(options)
	s.Require().Nil(err)

	s.Require().Equal([]*model.User{
		sanitized(s.u1),
		sanitized(s.u2),
		sanitized(s.u3),
		sanitized(s.u4),
		sanitized(s.u5),
		sanitized(s.u6),
		sanitized(s.u7),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestGetOffset0Limit1() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:    0,
		PerPage: 1,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u1),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestGetAll() {
	actual, err := s.Store().User().GetAll()
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		s.u1,
		s.u2,
		s.u3,
		s.u4,
		s.u5,
		s.u6,
		s.u7,
	}, actual)
}

// etag changes for all after user creation
func (s *UserStoreGetAllProfilesTS) TestEtagChangesForAllAfterUserCreation() {
	etag := s.Store().User().GetEtagForAllProfiles()

	uNew := &model.User{}
	uNew.Email = storetest.MakeEmail()
	_, err := s.Store().User().Save(uNew)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(uNew.Id)) }()

	updatedEtag := s.Store().User().GetEtagForAllProfiles()
	s.Require().NotEqual(etag, updatedEtag)
}

func (s *UserStoreGetAllProfilesTS) TestFilterToSystemAdminRole() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:    0,
		PerPage: 10,
		Role:    "system_admin",
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u5),
		sanitized(s.u6),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestFilterToSystemAdminRoleInactive() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:     0,
		PerPage:  10,
		Role:     "system_admin",
		Inactive: true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u6),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestFilterToInactive() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:     0,
		PerPage:  10,
		Inactive: true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u6),
		sanitized(s.u7),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestFilterToActive() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:    0,
		PerPage: 10,
		Active:  true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u1),
		sanitized(s.u2),
		sanitized(s.u3),
		sanitized(s.u4),
		sanitized(s.u5),
	}, actual)
}

func (s *UserStoreGetAllProfilesTS) TestTryToFilterToActiveAndInactive() {
	actual, err := s.Store().User().GetAllProfiles(&model.UserGetOptions{
		Page:     0,
		PerPage:  10,
		Inactive: true,
		Active:   true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u6),
		sanitized(s.u7),
	}, actual)
}

type UserStoreGetProfilesTS struct {
	suite.Suite
	StoreTestSuite

	u1     *model.User
	u2     *model.User
	u3     *model.User
	u4     *model.User
	u5     *model.User
	u6     *model.User
	u7     *model.User
	teamId string
}

func TestUserStoreGetProfilesTS(t *testing.T) {
	StoreTestSuiteWithSqlSupplier(t, &UserStoreGetProfilesTS{}, func(t *testing.T, testSuite StoreTestBaseSuite) {
		suite.Run(t, testSuite)
	})
}

func (s *UserStoreGetProfilesTS) SetupSuite() {
	teamId := model.NewId()
	s.teamId = teamId

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)
	s.u1 = u1

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)
	s.u2 = u2

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	s.u3 = u3
	s.u3.IsBot = true
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)

	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
		Roles:    "system_admin",
	})
	s.Require().Nil(err)
	defer func() {}()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1)
	s.Require().Nil(nErr)
	s.u4 = u4

	u5, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u5" + model.NewId(),
		DeleteAt: model.GetMillis(),
	})
	s.Require().Nil(err)
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u5.Id}, -1)
	s.Require().Nil(nErr)
	s.u5 = u5
}

func (s *UserStoreGetProfilesTS) TearDownSuite() {
	s.Require().Nil(s.Store().User().PermanentDelete(s.u1.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u2.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u3.Id))
	s.Require().Nil(s.Store().Bot().PermanentDelete(s.u3.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u4.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u5.Id))
}

// get page 0, perPage 100
func (s *UserStoreGetProfilesTS) TestGetPage0PerPage100() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  100,
	})
	s.Require().Nil(err)

	s.Require().Equal([]*model.User{
		sanitized(s.u1),
		sanitized(s.u2),
		sanitized(s.u3),
		sanitized(s.u4),
		sanitized(s.u5),
	}, actual)
}

// get page 0, perPage 1
func (s *UserStoreGetProfilesTS) TestGetPage0PerPage1() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  1,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{sanitized(s.u1)}, actual)
}

// get unknown team id
func (s *UserStoreGetProfilesTS) TestGetUnknownTeamId() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: "123",
		Page:     0,
		PerPage:  100,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{}, actual)
}

// etag changes for all after user creation
func (s *UserStoreGetProfilesTS) TestEtagChangesForAllAfterUserCreation() {
	etag := s.Store().User().GetEtagForProfiles(s.teamId)
	uNew := &model.User{}
	uNew.Email = MakeEmail()
	_, err := s.Store().User().Save(uNew)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(uNew.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: s.teamId, UserId: uNew.Id}, -1)
	s.Require().Nil(nErr)
	updatedEtag := s.Store().User().GetEtagForProfiles(s.teamId)
	s.Require().NotEqual(etag, updatedEtag)
}

// filter to system_admin role
func (s *UserStoreGetProfilesTS) TestFilterToSystemAdminRole() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  10,
		Role:     "system_admin",
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u4),
	}, actual)
}

// filter to inactive
func (s *UserStoreGetProfilesTS) TestFilterToInActive() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  10,
		Inactive: true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u5),
	}, actual)
}

// filter to active
func (s *UserStoreGetProfilesTS) TestFilterToActive() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  10,
		Active:   true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u1),
		sanitized(s.u2),
		sanitized(s.u3),
		sanitized(s.u4),
	}, actual)
}

// try to filter to active and inactive
func (s *UserStoreGetProfilesTS) TestFilterToActiveAndInactive() {
	actual, err := s.Store().User().GetProfiles(&model.UserGetOptions{
		InTeamId: s.teamId,
		Page:     0,
		PerPage:  10,
		Inactive: true,
		Active:   true,
	})
	s.Require().Nil(err)
	s.Require().Equal([]*model.User{
		sanitized(s.u5),
	}, actual)
}

type UserStoreGetProfilesByIdsTS struct {
	suite.Suite
	StoreTestSuite

	u1     *model.User
	u2     *model.User
	u3     *model.User
	u4     *model.User
	u5     *model.User
	u6     *model.User
	u7     *model.User
	teamId string
}

func TestUserStoreGetProfilesByIdsTS(t *testing.T) {
	StoreTestSuiteWithSqlSupplier(t, &UserStoreGetProfilesByIdsTS{}, func(t *testing.T, testSuite StoreTestBaseSuite) {
		suite.Run(t, testSuite)
	})
}

func (s *UserStoreGetProfilesByIdsTS) SetupSuite() {
	teamId := model.NewId()
	s.teamId = teamId

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)
	s.u1 = u1

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)
	s.u2 = u2

	time.Sleep(time.Millisecond)
	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	s.u3 = u3

	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	s.u4 = u4
}

func (s *UserStoreGetProfilesByIdsTS) TearDownSuite() {
	s.Require().Nil(s.Store().User().PermanentDelete(s.u1.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u2.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u3.Id))
	s.Require().Nil(s.Store().User().PermanentDelete(s.u4.Id))
}

// get u1 by id, no caching
func (s *UserStoreGetProfilesByIdsTS) TestGetU1ByIdNoCaching() {
	users, err := s.Store().User().GetProfileByIds([]string{s.u1.Id}, nil, false)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{s.u1}, users)
}

func (s *UserStoreGetProfilesByIdsTS) TestGetU1ByIdCaching() {
	users, err := s.Store().User().GetProfileByIds([]string{s.u1.Id}, nil, true)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{s.u1}, users)
}

// get u1, u2, u3 by id, no caching
func (s *UserStoreGetProfilesByIdsTS) TestGetU1U2U3ByIdNoCaching() {
	users, err := s.Store().User().GetProfileByIds([]string{s.u1.Id, s.u2.Id, s.u3.Id}, nil, false)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{s.u1, s.u2, s.u3}, users)
}

// get u1, u2, u3 by id, caching
func (s *UserStoreGetProfilesByIdsTS) TestGetU1U2U3ByIdCaching() {
	users, err := s.Store().User().GetProfileByIds([]string{s.u1.Id, s.u2.Id, s.u3.Id}, nil, true)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{s.u1, s.u2, s.u3}, users)
}

// get unknown id, caching
func (s *UserStoreGetProfilesByIdsTS) TestGetUnknownByIdCaching() {
	users, err := s.Store().User().GetProfileByIds([]string{"123"}, nil, true)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{}, users)
}

// should only return users with UpdateAt greater than the since time
func (s *UserStoreGetProfilesByIdsTS) TestReturnUsersUpdateAtGreater() {
	users, err := s.Store().User().GetProfileByIds([]string{s.u1.Id, s.u2.Id, s.u3.Id, s.u4.Id}, &store.UserGetByIdsOpts{
		Since: s.u2.CreateAt,
	}, true)
	s.Require().Nil(err)
	s.Assert().Equal([]*model.User{s.u3, s.u4}, users)
}

func (s *UserStoreTS) TestGetProfileByGroupChannelIdsForUser() {
	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()

	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u4.Id)) }()

	gc1, nErr := s.Store().Channel().Save(&model.Channel{
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_GROUP,
	}, -1)
	s.Require().Nil(nErr)

	for _, uId := range []string{u1.Id, u2.Id, u3.Id} {
		_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc1.Id,
			UserId:      uId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		s.Require().Nil(nErr)
	}

	gc2, nErr := s.Store().Channel().Save(&model.Channel{
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_GROUP,
	}, -1)
	s.Require().Nil(nErr)

	for _, uId := range []string{u1.Id, u3.Id, u4.Id} {
		_, nErr = s.Store().Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc2.Id,
			UserId:      uId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		s.Require().Nil(nErr)
	}

	testCases := []struct {
		Name                       string
		UserId                     string
		ChannelIds                 []string
		ExpectedUserIdsByChannel   map[string][]string
		EnsureChannelsNotInResults []string
	}{
		{
			Name:       "Get group 1 as user 1",
			UserId:     u1.Id,
			ChannelIds: []string{gc1.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u2.Id, u3.Id},
			},
			EnsureChannelsNotInResults: []string{},
		},
		{
			Name:       "Get groups 1 and 2 as user 1",
			UserId:     u1.Id,
			ChannelIds: []string{gc1.Id, gc2.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u2.Id, u3.Id},
				gc2.Id: {u3.Id, u4.Id},
			},
			EnsureChannelsNotInResults: []string{},
		},
		{
			Name:       "Get groups 1 and 2 as user 2",
			UserId:     u2.Id,
			ChannelIds: []string{gc1.Id, gc2.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u1.Id, u3.Id},
			},
			EnsureChannelsNotInResults: []string{gc2.Id},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.Name, func(t *testing.T) {
			res, err := s.Store().User().GetProfileByGroupChannelIdsForUser(tc.UserId, tc.ChannelIds)
			s.Require().Nil(err)

			for channelId, expectedUsers := range tc.ExpectedUserIdsByChannel {
				users, ok := res[channelId]
				s.Require().True(ok)

				var userIds []string
				for _, user := range users {
					userIds = append(userIds, user.Id)
				}
				s.Require().ElementsMatch(expectedUsers, userIds)
			}

			for _, channelId := range tc.EnsureChannelsNotInResults {
				_, ok := res[channelId]
				s.Require().False(ok)
			}
		})
	}
}

func (s *UserStoreTS) TestGetByEmail() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get u1 by email", func(t *testing.T) {
		u, err := s.Store().User().GetByEmail(u1.Email)
		s.Require().Nil(err)
		s.Assert().Equal(u1, u)
	})

	s.T().Run("get u2 by email", func(t *testing.T) {
		u, err := s.Store().User().GetByEmail(u2.Email)
		s.Require().Nil(err)
		s.Assert().Equal(u2, u)
	})

	s.T().Run("get u3 by email", func(t *testing.T) {
		u, err := s.Store().User().GetByEmail(u3.Email)
		s.Require().Nil(err)
		s.Assert().Equal(u3, u)
	})

	s.T().Run("get by empty email", func(t *testing.T) {
		_, err := s.Store().User().GetByEmail("")
		s.Require().NotNil(err)
	})

	s.T().Run("get by unknown", func(t *testing.T) {
		_, err := s.Store().User().GetByEmail("unknown")
		s.Require().NotNil(err)
	})
}

func (s *UserStoreTS) TestGetByAuthData() {
	teamId := model.NewId()
	auth1 := model.NewId()
	auth3 := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthData:    &auth1,
		AuthService: "service",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthData:    &auth3,
		AuthService: "service2",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get by u1 auth", func(t *testing.T) {
		u, err := s.Store().User().GetByAuth(u1.AuthData, u1.AuthService)
		s.Require().Nil(err)
		s.Assert().Equal(u1, u)
	})

	s.T().Run("get by u3 auth", func(t *testing.T) {
		u, err := s.Store().User().GetByAuth(u3.AuthData, u3.AuthService)
		s.Require().Nil(err)
		s.Assert().Equal(u3, u)
	})

	s.T().Run("get by u1 auth, unknown service", func(t *testing.T) {
		_, err := s.Store().User().GetByAuth(u1.AuthData, "unknown")
		s.Require().NotNil(err)
		var nfErr *store.ErrNotFound
		s.Require().True(errors.As(err, &nfErr))
	})

	s.T().Run("get by unknown auth, u1 service", func(t *testing.T) {
		unknownAuth := ""
		_, err := s.Store().User().GetByAuth(&unknownAuth, u1.AuthService)
		s.Require().NotNil(err)
		var invErr *store.ErrInvalidInput
		s.Require().True(errors.As(err, &invErr))
	})

	s.T().Run("get by unknown auth, unknown service", func(t *testing.T) {
		unknownAuth := ""
		_, err := s.Store().User().GetByAuth(&unknownAuth, "unknown")
		s.Require().NotNil(err)
		var invErr *store.ErrInvalidInput
		s.Require().True(errors.As(err, &invErr))
	})
}

func (s *UserStoreTS) TestGetByUsername() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get u1 by username", func(t *testing.T) {
		result, err := s.Store().User().GetByUsername(u1.Username)
		s.Require().Nil(err)
		s.Assert().Equal(u1, result)
	})

	s.T().Run("get u2 by username", func(t *testing.T) {
		result, err := s.Store().User().GetByUsername(u2.Username)
		s.Require().Nil(err)
		s.Assert().Equal(u2, result)
	})

	s.T().Run("get u3 by username", func(t *testing.T) {
		result, err := s.Store().User().GetByUsername(u3.Username)
		s.Require().Nil(err)
		s.Assert().Equal(u3, result)
	})

	s.T().Run("get by empty username", func(t *testing.T) {
		_, err := s.Store().User().GetByUsername("")
		s.Require().NotNil(err)
		var nfErr *store.ErrNotFound
		s.Require().True(errors.As(err, &nfErr))
	})

	s.T().Run("get by unknown", func(t *testing.T) {
		_, err := s.Store().User().GetByUsername("unknown")
		s.Require().NotNil(err)
		var nfErr *store.ErrNotFound
		s.Require().True(errors.As(err, &nfErr))
	})
}

func (s *UserStoreTS) TestGetForLogin() {
	teamId := model.NewId()
	auth := model.NewId()
	auth2 := model.NewId()
	auth3 := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_GITLAB,
		AuthData:    &auth,
	})

	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u2" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_LDAP,
		AuthData:    &auth2,
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_LDAP,
		AuthData:    &auth3,
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("get u1 by username, allow both", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(u1.Username, true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u1, user)
	})

	s.T().Run("get u1 by username, check for case issues", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(strings.ToUpper(u1.Username), true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u1, user)
	})

	s.T().Run("get u1 by username, allow only email", func(t *testing.T) {
		_, err := s.Store().User().GetForLogin(u1.Username, false, true)
		s.Require().NotNil(err)
		s.Require().Equal("user not found", err.Error())
	})

	s.T().Run("get u1 by email, allow both", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(u1.Email, true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u1, user)
	})

	s.T().Run("get u1 by email, check for case issues", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(strings.ToUpper(u1.Email), true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u1, user)
	})

	s.T().Run("get u1 by email, allow only username", func(t *testing.T) {
		_, err := s.Store().User().GetForLogin(u1.Email, true, false)
		s.Require().NotNil(err)
		s.Require().Equal("user not found", err.Error())
	})

	s.T().Run("get u2 by username, allow both", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(u2.Username, true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u2, user)
	})

	s.T().Run("get u2 by email, allow both", func(t *testing.T) {
		user, err := s.Store().User().GetForLogin(u2.Email, true, true)
		s.Require().Nil(err)
		s.Assert().Equal(u2, user)
	})

	s.T().Run("get u2 by username, allow neither", func(t *testing.T) {
		_, err := s.Store().User().GetForLogin(u2.Username, false, false)
		s.Require().NotNil(err)
		s.Require().Equal("sign in with username and email are disabled", err.Error())
	})
}

func (s *UserStoreTS) TestGetAllUsingAuthService() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthService: "service",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u2" + model.NewId(),
		AuthService: "service",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthService: "service2",
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()

	s.T().Run("get by unknown auth service", func(t *testing.T) {
		users, err := s.Store().User().GetAllUsingAuthService("unknown")
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{}, users)
	})

	s.T().Run("get by auth service", func(t *testing.T) {
		users, err := s.Store().User().GetAllUsingAuthService("service")
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u1, u2}, users)
	})

	s.T().Run("get by other auth service", func(t *testing.T) {
		users, err := s.Store().User().GetAllUsingAuthService("service2")
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{u3}, users)
	})
}

func (s *UserStoreTS) TestUpdateFailedPasswordAttempts() {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	err = s.Store().User().UpdateFailedPasswordAttempts(u1.Id, 3)
	s.Require().Nil(err)

	user, err := s.Store().User().Get(u1.Id)
	s.Require().Nil(err)
	s.Require().Equal(3, user.FailedAttempts, "FailedAttempts not updated correctly")
}

func (s *UserStoreTS) TestUserStoreGetSystemAdminProfiles() {
	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	s.T().Run("all system admin profiles", func(t *testing.T) {
		result, userError := s.Store().User().GetSystemAdminProfiles()
		s.Require().Nil(userError)
		s.Assert().Equal(map[string]*model.User{
			u1.Id: sanitized(u1),
			u3.Id: sanitized(u3),
		}, result)
	})
}

func (s *UserStoreTS) TestAnalyticsActiveCount() {
	s.cleanupStatusStore()

	// Create 5 users statuses u0, u1, u2, u3, u4.
	// u4 is also a bot
	u0, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u0" + model.NewId(),
	})
	s.Require().Nil(err)
	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() {
		s.Require().Nil(s.Store().User().PermanentDelete(u0.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u1.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u2.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u3.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u4.Id))
	}()

	_, nErr := s.Store().Bot().Save(&model.Bot{
		UserId:   u4.Id,
		Username: u4.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)

	millis := model.GetMillis()
	millisTwoDaysAgo := model.GetMillis() - (2 * DAY_MILLISECONDS)
	millisTwoMonthsAgo := model.GetMillis() - (2 * MONTH_MILLISECONDS)

	// u0 last activity status is two months ago.
	// u1 last activity status is two days ago.
	// u2, u3, u4 last activity is within last day
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u0.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoMonthsAgo}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoDaysAgo}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u2.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u3.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u4.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))

	// Daily counts (without bots)
	count, err := s.Store().User().AnalyticsActiveCount(DAY_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: true})
	s.Require().Nil(err)
	s.Assert().Equal(int64(2), count)

	// Daily counts (with bots)
	count, err = s.Store().User().AnalyticsActiveCount(DAY_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: true})
	s.Require().Nil(err)
	s.Assert().Equal(int64(3), count)

	// Monthly counts (without bots)
	count, err = s.Store().User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: true})
	s.Require().Nil(err)
	s.Assert().Equal(int64(3), count)

	// Monthly counts - (with bots)
	count, err = s.Store().User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: true})
	s.Require().Nil(err)
	s.Assert().Equal(int64(4), count)

	// Monthly counts - (with bots, excluding deleted)
	count, err = s.Store().User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: false})
	s.Require().Nil(err)
	s.Assert().Equal(int64(4), count)
}

func (s *UserStoreTS) TestAnalyticsActiveCountForPeriod() {

	s.cleanupStatusStore()

	// Create 5 users statuses u0, u1, u2, u3, u4.
	// u4 is also a bot
	u0, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u0" + model.NewId(),
	})
	s.Require().Nil(err)
	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	u4, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() {
		s.Require().Nil(s.Store().User().PermanentDelete(u0.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u1.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u2.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u3.Id))
		s.Require().Nil(s.Store().User().PermanentDelete(u4.Id))
	}()

	_, nErr := s.Store().Bot().Save(&model.Bot{
		UserId:   u4.Id,
		Username: u4.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)

	millis := model.GetMillis()
	millisTwoDaysAgo := model.GetMillis() - (2 * DAY_MILLISECONDS)
	millisTwoMonthsAgo := model.GetMillis() - (2 * MONTH_MILLISECONDS)

	// u0 last activity status is two months ago.
	// u1 last activity status is one month ago
	// u2 last activiy is two days ago
	// u2 last activity is one day ago
	// u3 last activity is within last day
	// u4 last activity is within last day
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u0.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoMonthsAgo}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoMonthsAgo + MONTH_MILLISECONDS}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u2.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoDaysAgo}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u3.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoDaysAgo + DAY_MILLISECONDS}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u4.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))

	// Two months to two days (without bots)
	count, nerr := s.Store().User().AnalyticsActiveCountForPeriod(millisTwoMonthsAgo, millisTwoDaysAgo, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: false})
	s.Require().Nil(nerr)
	s.Assert().Equal(int64(2), count)

	// Two months to two days (without bots)
	count, nerr = s.Store().User().AnalyticsActiveCountForPeriod(millisTwoMonthsAgo, millisTwoDaysAgo, model.UserCountOptions{IncludeBotAccounts: false, IncludeDeleted: true})
	s.Require().Nil(nerr)
	s.Assert().Equal(int64(2), count)

	// Two days to present - (with bots)
	count, nerr = s.Store().User().AnalyticsActiveCountForPeriod(millisTwoDaysAgo, millis, model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: false})
	s.Require().Nil(nerr)
	s.Assert().Equal(int64(2), count)

	// Two days to present - (with bots, excluding deleted)
	count, nerr = s.Store().User().AnalyticsActiveCountForPeriod(millisTwoDaysAgo, millis, model.UserCountOptions{IncludeBotAccounts: true, IncludeDeleted: true})
	s.Require().Nil(nerr)
	s.Assert().Equal(int64(2), count)
}

func (s *UserStoreTS) TestAnalyticsGetInactiveUsersCount() {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	count, err := s.Store().User().AnalyticsGetInactiveUsersCount()
	s.Require().Nil(err)

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.DeleteAt = model.GetMillis()
	_, err = s.Store().User().Save(u2)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	newCount, err := s.Store().User().AnalyticsGetInactiveUsersCount()
	s.Require().Nil(err)
	s.Require().Equal(count, newCount-1, "Expected 1 more inactive users but found otherwise.")
}

func (s *UserStoreTS) TestAnalyticsGetSystemAdminCount() {
	countBefore, err := s.Store().User().AnalyticsGetSystemAdminCount()
	s.Require().Nil(err)

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	u1.Roles = "system_user system_admin"

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Username = model.NewId()

	_, nErr := s.Store().User().Save(&u1)
	s.Require().Nil(nErr, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	_, nErr = s.Store().User().Save(&u2)
	s.Require().Nil(nErr, "couldn't save user")

	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	result, err := s.Store().User().AnalyticsGetSystemAdminCount()
	s.Require().Nil(err)
	s.Require().Equal(countBefore+1, result, "Did not get the expected number of system admins.")

}

func (s *UserStoreTS) TestAnalyticsGetGuestCount() {
	countBefore, err := s.Store().User().AnalyticsGetGuestCount()
	s.Require().Nil(err)

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	u1.Roles = "system_user system_admin"

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Username = model.NewId()
	u2.Roles = "system_user"

	u3 := model.User{}
	u3.Email = MakeEmail()
	u3.Username = model.NewId()
	u3.Roles = "system_guest"

	_, nErr := s.Store().User().Save(&u1)
	s.Require().Nil(nErr, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	_, nErr = s.Store().User().Save(&u2)
	s.Require().Nil(nErr, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	_, nErr = s.Store().User().Save(&u3)
	s.Require().Nil(nErr, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()

	result, err := s.Store().User().AnalyticsGetGuestCount()
	s.Require().Nil(err)
	s.Require().Equal(countBefore+1, result, "Did not get the expected number of guests.")
}

func (s *UserStoreTS) TestAnalyticsGetExternalUsers() {
	localHostDomain := "mattermost.com"
	result, err := s.Store().User().AnalyticsGetExternalUsers(localHostDomain)
	s.Require().Nil(err)
	s.Assert().False(result)

	u1 := model.User{}
	u1.Email = "a@mattermost.com"
	u1.Username = model.NewId()
	u1.Roles = "system_user system_admin"

	u2 := model.User{}
	u2.Email = "b@example.com"
	u2.Username = model.NewId()
	u2.Roles = "system_user"

	u3 := model.User{}
	u3.Email = "c@test.com"
	u3.Username = model.NewId()
	u3.Roles = "system_guest"

	_, err = s.Store().User().Save(&u1)
	s.Require().Nil(err, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()

	_, err = s.Store().User().Save(&u2)
	s.Require().Nil(err, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()

	_, err = s.Store().User().Save(&u3)
	s.Require().Nil(err, "couldn't save user")
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()

	result, err = s.Store().User().AnalyticsGetExternalUsers(localHostDomain)
	s.Require().Nil(err)
	s.Assert().True(result)
}

func (s *UserStoreTS) TestUserUnreadCount() {
	teamId := model.NewId()

	c1 := model.Channel{}
	c1.TeamId = teamId
	c1.DisplayName = "Unread Messages"
	c1.Name = "unread-messages-" + model.NewId()
	c1.Type = model.CHANNEL_OPEN

	c2 := model.Channel{}
	c2.TeamId = teamId
	c2.DisplayName = "Unread Direct"
	c2.Name = "unread-direct-" + model.NewId()
	c2.Type = model.CHANNEL_DIRECT

	u1 := &model.User{}
	u1.Username = "user1" + model.NewId()
	u1.Email = MakeEmail()
	_, err := s.Store().User().Save(u1)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Username = "user2" + model.NewId()
	_, err = s.Store().User().Save(u2)
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	_, nErr = s.Store().Channel().Save(&c1, -1)
	s.Require().Nil(nErr, "couldn't save item")

	m1 := model.ChannelMember{}
	m1.ChannelId = c1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = c1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	_, nErr = s.Store().Channel().SaveMember(&m2)
	s.Require().Nil(nErr)

	m1.ChannelId = c2.Id
	m2.ChannelId = c2.Id

	_, nErr = s.Store().Channel().SaveDirectChannel(&c2, &m1, &m2)
	s.Require().Nil(nErr, "couldn't save direct channel")

	p1 := model.Post{}
	p1.ChannelId = c1.Id
	p1.UserId = u1.Id
	p1.Message = "this is a message for @" + u2.Username

	// Post one message with mention to open channel
	_, nErr = s.Store().Post().Save(&p1)
	s.Require().Nil(nErr)
	nErr = s.Store().Channel().IncrementMentionCount(c1.Id, u2.Id, false)
	s.Require().Nil(nErr)

	// Post 2 messages without mention to direct channel
	p2 := model.Post{}
	p2.ChannelId = c2.Id
	p2.UserId = u1.Id
	p2.Message = "first message"

	_, nErr = s.Store().Post().Save(&p2)
	s.Require().Nil(nErr)
	nErr = s.Store().Channel().IncrementMentionCount(c2.Id, u2.Id, false)
	s.Require().Nil(nErr)

	p3 := model.Post{}
	p3.ChannelId = c2.Id
	p3.UserId = u1.Id
	p3.Message = "second message"
	_, nErr = s.Store().Post().Save(&p3)
	s.Require().Nil(nErr)

	nErr = s.Store().Channel().IncrementMentionCount(c2.Id, u2.Id, false)
	s.Require().Nil(nErr)

	badge, unreadCountErr := s.Store().User().GetUnreadCount(u2.Id)
	s.Require().Nil(unreadCountErr)
	s.Require().Equal(int64(3), badge, "should have 3 unread messages")

	badge, unreadCountErr = s.Store().User().GetUnreadCountForChannel(u2.Id, c1.Id)
	s.Require().Nil(unreadCountErr)
	s.Require().Equal(int64(1), badge, "should have 1 unread messages for that channel")

	badge, unreadCountErr = s.Store().User().GetUnreadCountForChannel(u2.Id, c2.Id)
	s.Require().Nil(unreadCountErr)
	s.Require().Equal(int64(2), badge, "should have 2 unread messages for that channel")
}

func (s *UserStoreTS) TestGetRecentlyActiveUsersForTeam() {

	s.cleanupStatusStore()

	teamId := model.NewId()

	u1, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u1.Id)) }()
	_, nErr := s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	s.Require().Nil(nErr)

	u2, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u2.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	s.Require().Nil(nErr)

	u3, err := s.Store().User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	s.Require().Nil(err)
	defer func() { s.Require().Nil(s.Store().User().PermanentDelete(u3.Id)) }()
	_, nErr = s.Store().Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	s.Require().Nil(nErr)
	_, nErr = s.Store().Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	s.Require().Nil(nErr)
	u3.IsBot = true
	defer func() { s.Require().Nil(s.Store().Bot().PermanentDelete(u3.Id)) }()

	millis := model.GetMillis()
	u3.LastActivityAt = millis
	u2.LastActivityAt = millis - 1
	u1.LastActivityAt = millis - 1

	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u1.LastActivityAt, ActiveChannel: ""}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u2.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u2.LastActivityAt, ActiveChannel: ""}))
	s.Require().Nil(s.Store().Status().SaveOrUpdate(&model.Status{UserId: u3.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u3.LastActivityAt, ActiveChannel: ""}))

	s.T().Run("get team 1, offset 0, limit 100", func(t *testing.T) {
		users, err := s.Store().User().GetRecentlyActiveUsersForTeam(teamId, 0, 100, nil)
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{
			sanitized(u3),
			sanitized(u1),
			sanitized(u2),
		}, users)
	})

	s.T().Run("get team 1, offset 0, limit 1", func(t *testing.T) {
		users, err := s.Store().User().GetRecentlyActiveUsersForTeam(teamId, 0, 1, nil)
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{
			sanitized(u3),
		}, users)
	})

	s.T().Run("get team 1, offset 2, limit 1", func(t *testing.T) {
		users, err := s.Store().User().GetRecentlyActiveUsersForTeam(teamId, 2, 1, nil)
		s.Require().Nil(err)
		s.Assert().Equal([]*model.User{
			sanitized(u2),
		}, users)
	})
}