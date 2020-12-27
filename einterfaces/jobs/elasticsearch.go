// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package jobs

import (
	"github.com/masterhung0112/hk_server/model"
)

type ElasticsearchIndexerInterface interface {
	MakeWorker() model.Worker
}

type ElasticsearchAggregatorInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
