package app

import (
	"github.com/masterhung0112/go_server/mlog"
	"github.com/masterhung0112/go_server/model"
	"github.com/masterhung0112/go_server/store"
	"github.com/masterhung0112/go_server/utils"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

// CreateDefaultChannels creates channels in the given team for each channel returned by (*App).DefaultChannelNames.
//
func (a *App) CreateDefaultChannels(teamID string) ([]*model.Channel, *model.AppError) {
	displayNames := map[string]string{
		"town-square": utils.T("api.channel.create_default_channels.town_square"),
		"off-topic":   utils.T("api.channel.create_default_channels.off_topic"),
	}
	channels := []*model.Channel{}
	defaultChannelNames := a.DefaultChannelNames()
	for _, name := range defaultChannelNames {
		displayName := utils.TDefault(displayNames[name], name)
		channel := &model.Channel{DisplayName: displayName, Name: name, Type: model.CHANNEL_OPEN, TeamId: teamID}
		if _, err := a.CreateChannel(channel, false); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (a *App) GetChannel(channelId string) (*model.Channel, *model.AppError) {
	channel, err := a.Srv().Store.Channel().Get(channelId, true)
	if err != nil {
		var nfErr *store.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetChannel", "app.channel.get.existing.app_error", nil, nfErr.Error(), http.StatusNotFound)
		default:
			return nil, model.NewAppError("GetChannel", "app.channel.get.find.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	}
	return channel, nil
}

func (a *App) addUserToChannel(user *model.User, channel *model.Channel, teamMember *model.TeamMember) (*model.ChannelMember, *model.AppError) {
	if channel.Type != model.CHANNEL_OPEN && channel.Type != model.CHANNEL_PRIVATE {
		return nil, model.NewAppError("AddUserToChannel", "api.channel.add_user_to_channel.type.app_error", nil, "", http.StatusBadRequest)
	}

	channelMember, err := a.Srv().Store.Channel().GetMember(channel.Id, user.Id)
	if err != nil {
		if err.Id != store.MISSING_CHANNEL_MEMBER_ERROR {
			return nil, err
		}
	} else {
		return channelMember, nil
	}

	if channel.IsGroupConstrained() {
		nonMembers, err := a.FilterNonGroupChannelMembers([]string{user.Id}, channel)
		if err != nil {
			return nil, model.NewAppError("addUserToChannel", "api.channel.add_user_to_channel.type.app_error", nil, "", http.StatusInternalServerError)
		}
		if len(nonMembers) > 0 {
			return nil, model.NewAppError("addUserToChannel", "api.channel.add_members.user_denied", map[string]interface{}{"UserIDs": nonMembers}, "", http.StatusBadRequest)
		}
	}

	newMember := &model.ChannelMember{
		ChannelId:   channel.Id,
		UserId:      user.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
		SchemeGuest: user.IsGuest(),
		SchemeUser:  !user.IsGuest(),
	}

	//TODO: Open
	// if !user.IsGuest() {
	// 	var userShouldBeAdmin bool
	// 	userShouldBeAdmin, err = a.UserIsInAdminRoleGroup(user.Id, channel.Id, model.GroupSyncableTypeChannel)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	newMember.SchemeAdmin = userShouldBeAdmin
	// }

	newMember, err = a.Srv().Store.Channel().SaveMember(newMember)
	if err != nil {
		mlog.Error("Failed to add member", mlog.String("user_id", user.Id), mlog.String("channel_id", channel.Id), mlog.Err(err))
		return nil, model.NewAppError("AddUserToChannel", "api.channel.add_user.to.channel.failed.app_error", nil, "", http.StatusInternalServerError)
	}
	a.WaitForChannelMembership(channel.Id, user.Id)

	//TODO: Open
	// if nErr := a.Srv().Store.ChannelMemberHistory().LogJoinEvent(user.Id, channel.Id, model.GetMillis()); nErr != nil {
	// 	mlog.Error("Failed to update ChannelMemberHistory table", mlog.Err(nErr))
	// 	return nil, model.NewAppError("AddUserToChannel", "app.channel_member_history.log_join_event.internal_error", nil, nErr.Error(), http.StatusInternalServerError)
	// }

	//TODO: Open
	// a.InvalidateCacheForUser(user.Id)
	// a.invalidateCacheForChannelMembers(channel.Id)

	return newMember, nil
}

func (a *App) AddUserToChannel(user *model.User, channel *model.Channel) (*model.ChannelMember, *model.AppError) {
	teamMember, err := a.Srv().Store.Team().GetMember(channel.TeamId, user.Id)

	if err != nil {
		return nil, err
	}
	if teamMember.DeleteAt > 0 {
		return nil, model.NewAppError("AddUserToChannel", "api.channel.add_user.to.channel.failed.deleted.app_error", nil, "", http.StatusBadRequest)
	}

	newMember, err := a.addUserToChannel(user, channel, teamMember)
	if err != nil {
		return nil, err
	}

	//TODO: Open
	// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_ADDED, "", channel.Id, "", nil)
	// message.Add("user_id", user.Id)
	// message.Add("team_id", channel.TeamId)
	// a.Publish(message)

	return newMember, nil
}

func (a *App) WaitForChannelMembership(channelId string, userId string) {
	if len(a.Config().SqlSettings.DataSourceReplicas) == 0 {
		return
	}

	now := model.GetMillis()

	for model.GetMillis()-now < 12000 {

		time.Sleep(100 * time.Millisecond)

		_, err := a.Srv().Store.Channel().GetMember(channelId, userId)

		// If the membership was found then return
		if err == nil {
			return
		}

		// If we received an error, but it wasn't a missing channel member then return
		if err.Id != store.MISSING_CHANNEL_MEMBER_ERROR {
			return
		}
	}

	mlog.Error("WaitForChannelMembership giving up", mlog.String("channel_id", channelId), mlog.String("user_id", userId))
}

func (a *App) JoinDefaultChannels(teamId string, user *model.User, shouldBeAdmin bool, userRequestorId string) *model.AppError {
	//TODO: Open
	// var requestor *model.User
	// if userRequestorId != "" {
	// 	var err *model.AppError
	// 	requestor, err = a.Srv().Store.User().Get(userRequestorId)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	var err *model.AppError
	for _, channelName := range a.DefaultChannelNames() {
		channel, channelErr := a.Srv().Store.Channel().GetByName(teamId, channelName, true)
		if channelErr != nil {
			var nfErr *store.ErrNotFound
			switch {
			case errors.As(err, &nfErr):
				err = model.NewAppError("JoinDefaultChannels", "app.channel.get_by_name.missing.app_error", nil, nfErr.Error(), http.StatusNotFound)
			default:
				err = model.NewAppError("JoinDefaultChannels", "app.channel.get_by_name.existing.app_error", nil, channelErr.Error(), http.StatusInternalServerError)
			}
			continue
		}

		if channel.Type != model.CHANNEL_OPEN {
			continue
		}

		cm := &model.ChannelMember{
			ChannelId:   channel.Id,
			UserId:      user.Id,
			SchemeGuest: user.IsGuest(),
			SchemeUser:  !user.IsGuest(),
			SchemeAdmin: shouldBeAdmin,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		}

		_, err = a.Srv().Store.Channel().SaveMember(cm)
		//TODO: Open
		// if histErr := a.Srv().Store.ChannelMemberHistory().LogJoinEvent(user.Id, channel.Id, model.GetMillis()); histErr != nil {
		// 	mlog.Error("Failed to update ChannelMemberHistory table", mlog.Err(histErr))
		// 	return model.NewAppError("JoinDefaultChannels", "app.channel_member_history.log_join_event.internal_error", nil, histErr.Error(), http.StatusInternalServerError)
		// }

		//TODO: Open
		// if *a.Config().ServiceSettings.ExperimentalEnableDefaultChannelLeaveJoinMessages {
		// 	a.postJoinMessageForDefaultChannel(user, requestor, channel)
		// }

		//TODO: Open
		// a.invalidateCacheForChannelMembers(channel.Id)

		//TODO: Open
		// message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_USER_ADDED, "", channel.Id, "", nil)
		// message.Add("user_id", user.Id)
		// message.Add("team_id", channel.TeamId)
		// a.Publish(message)

	}

	return err
}

// DefaultChannelNames returns the list of system-wide default channel names.
//
// By default the list will be (not necessarily in this order):
//	['town-square', 'off-topic']
// However, if TeamSettings.ExperimentalDefaultChannels contains a list of channels then that list will replace
// 'off-topic' and be included in the return results in addition to 'town-square'. For example:
//	['town-square', 'game-of-thrones', 'wow']
//
func (a *App) DefaultChannelNames() []string {
	names := []string{"town-square"}

	if len(a.Config().TeamSettings.ExperimentalDefaultChannels) == 0 {
		names = append(names, "off-topic")
	} else {
		seenChannels := map[string]bool{"town-square": true}
		for _, channelName := range a.Config().TeamSettings.ExperimentalDefaultChannels {
			if !seenChannels[channelName] {
				names = append(names, channelName)
				seenChannels[channelName] = true
			}
		}
	}

	return names
}

func (a *App) CreateChannel(channel *model.Channel, addMember bool) (*model.Channel, *model.AppError) {
	channel.DisplayName = strings.TrimSpace(channel.DisplayName)
	sc, nErr := a.Srv().Store.Channel().Save(channel, *a.Config().TeamSettings.MaxChannelsPerTeam)
	if nErr != nil {
		var invErr *store.ErrInvalidInput
		var cErr *store.ErrConflict
		var ltErr *store.ErrLimitExceeded
		var appErr *model.AppError
		switch {
		case errors.As(nErr, &invErr):
			switch {
			case invErr.Entity == "Channel" && invErr.Field == "DeleteAt":
				return nil, model.NewAppError("CreateChannel", "store.sql_channel.save.archived_channel.app_error", nil, "", http.StatusBadRequest)
			case invErr.Entity == "Channel" && invErr.Field == "Type":
				return nil, model.NewAppError("CreateChannel", "store.sql_channel.save.direct_channel.app_error", nil, "", http.StatusBadRequest)
			case invErr.Entity == "Channel" && invErr.Field == "Id":
				return nil, model.NewAppError("CreateChannel", "store.sql_channel.save_channel.existing.app_error", nil, "id="+invErr.Value.(string), http.StatusBadRequest)
			}
		case errors.As(nErr, &cErr):
			return sc, model.NewAppError("CreateChannel", store.CHANNEL_EXISTS_ERROR, nil, cErr.Error(), http.StatusBadRequest)
		case errors.As(nErr, &ltErr):
			return nil, model.NewAppError("CreateChannel", "store.sql_channel.save_channel.limit.app_error", nil, ltErr.Error(), http.StatusBadRequest)
		case errors.As(nErr, &appErr): // in case we haven't converted to plain error.
			return nil, appErr
		default: // last fallback in case it doesn't map to an existing app error.
			return nil, model.NewAppError("CreateChannel", "app.channel.create_channel.internal_error", nil, nErr.Error(), http.StatusInternalServerError)
		}
	}

	if addMember {
		user, err := a.Srv().Store.User().Get(channel.CreatorId)
		if err != nil {
			return nil, err
		}

		cm := &model.ChannelMember{
			ChannelId:   sc.Id,
			UserId:      user.Id,
			SchemeGuest: user.IsGuest(),
			SchemeUser:  !user.IsGuest(),
			SchemeAdmin: true,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		}

		if _, err := a.Srv().Store.Channel().SaveMember(cm); err != nil {
			return nil, err
		}
		//TODO: Open
		// if err := a.Srv().Store.ChannelMemberHistory().LogJoinEvent(channel.CreatorId, sc.Id, model.GetMillis()); err != nil {
		// 	mlog.Error("Failed to update ChannelMemberHistory table", mlog.Err(err))
		// 	return nil, model.NewAppError("CreateChannel", "app.channel_member_history.log_join_event.internal_error", nil, err.Error(), http.StatusInternalServerError)
		// }

		//TODO: Open
		// a.InvalidateCacheForUser(channel.CreatorId)
	}

	//TODO: Open
	// if pluginsEnvironment := a.GetPluginsEnvironment(); pluginsEnvironment != nil {
	// 	a.Srv().Go(func() {
	// 		pluginContext := a.PluginContext()
	// 		pluginsEnvironment.RunMultiPluginHook(func(hooks plugin.Hooks) bool {
	// 			hooks.ChannelHasBeenCreated(pluginContext, sc)
	// 			return true
	// 		}, plugin.ChannelHasBeenCreatedId)
	// 	})
	// }

	return sc, nil
}
