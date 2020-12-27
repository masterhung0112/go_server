// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import (
	"github.com/masterhung0112/hk_server/app/plugin_api_tests"
	"github.com/masterhung0112/hk_server/model"
	"github.com/masterhung0112/hk_server/plugin"
)

type MyPlugin struct {
	plugin.MattermostPlugin
	configuration plugin_api_tests.BasicConfig
}

func (p *MyPlugin) OnConfigurationChange() error {
	if err := p.API.LoadPluginConfiguration(&p.configuration); err != nil {
		return err
	}
	return nil
}

func (p *MyPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {

	channels, err := p.API.GetChannelsForTeamForUser(p.configuration.BasicTeamId, p.configuration.BasicUserId, false)
	if err != nil {
		return nil, err.Error()
	}
	if len(channels) != 3 {
		return nil, "Returned invalid number of channels"
	}
	channels, err = p.API.GetChannelsForTeamForUser("invalidid", p.configuration.BasicUserId, false)
	if err == nil {
		return nil, "Expected to get an error while retrieving channels for invalid id"
	}
	if len(channels) != 0 {
		return nil, "Returned invalid number of channels"
	}
	return nil, "OK"
}

func main() {
	plugin.ClientMain(&MyPlugin{})
}