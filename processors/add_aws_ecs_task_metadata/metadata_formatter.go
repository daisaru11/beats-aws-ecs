package add_aws_ecs_task_metadata

import (
	"github.com/elastic/beats/libbeat/common"

	"github.com/daisaru11/beats-aws-ecs/ecs_task_metadata"
)

type MetadataFormatter struct {
}

func NewMetadataFormatter() *MetadataFormatter {
	return &MetadataFormatter{}
}

func (f *MetadataFormatter) FormatContainerMetadata(task *ecs_task_metadata.TaskMetadata, name string) common.MapStr {
	meta := f.FormatTaskMetadata(task)

	for _, container := range task.Containers {
		if container.Name == name {
			meta["container"] = common.MapStr{
				"name":        container.Name,
				"docker_name": container.DockerName,
				"image":       container.Image,
				"image_id":    container.ImageID,
			}
		}
	}

	return meta
}

func (f *MetadataFormatter) FormatTaskMetadata(task *ecs_task_metadata.TaskMetadata) common.MapStr {
	meta := common.MapStr{}

	meta["task"] = common.MapStr{
		"cluster":  task.Cluster,
		"arn":      task.TaskARN,
		"family":   task.Family,
		"revision": task.Revision,
	}

	return meta
}
