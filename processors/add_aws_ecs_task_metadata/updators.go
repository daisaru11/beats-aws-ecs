package add_aws_ecs_task_metadata

import (
	"github.com/daisaru11/beats-aws-ecs/ecs_task_metadata"
	"github.com/elastic/beats/libbeat/logp"
)

type Updater interface {
	Run()
}

type OnceAtStartupUpdater struct {
	client        *ecs_task_metadata.TaskMetadataClient
	updateHandler func(task *ecs_task_metadata.TaskMetadata)
}

func NewOnceAtStartupUpdater(
	client *ecs_task_metadata.TaskMetadataClient,
	updateHandler func(task *ecs_task_metadata.TaskMetadata),
) *OnceAtStartupUpdater {
	return &OnceAtStartupUpdater{
		client:        client,
		updateHandler: updateHandler,
	}
}

func (u *OnceAtStartupUpdater) Run() {
	logp.Debug("add_aws_ecs_task_metadata", "Starting to fetch ecs task metadata")
	task, err := u.client.GetTaskMetadata()

	if err != nil {
		logp.Err("failed to fetch ecs task metadata. %s", err)
		return
	}

	logp.Debug("add_aws_ecs_task_metadata", "Update task metadata: %v", task)
	u.updateHandler(task)
}
