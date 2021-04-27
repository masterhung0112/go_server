// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
package jobs

import "github.com/masterhung0112/hk_server/v5/model"

// ResendInvitationEmailJobInterface defines the interface for the job to resend invitation emails
type ResendInvitationEmailJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
