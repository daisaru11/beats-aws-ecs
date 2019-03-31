package add_aws_ecs_task_metadata

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/daisaru11/beats-aws-ecs/ecs_task_metadata"
)

type Indexer interface {
	GetIndices(task *ecs_task_metadata.TaskMetadata) []MetadataIndex
}

type MetadataIndex struct {
	Index string
	Data  common.MapStr
}

type Indexers struct {
	indexers []Indexer
}

func NewIndexers(config *addAwsEcsTaskMetadataConfig) *Indexers {
	indexers := []Indexer{}

	for _, indexerConfigs := range config.Indexers {
		for name, indexerConfig := range indexerConfigs {
			switch name {
			case "container_name":
				i, err := NewContainerNameIndexer(&indexerConfig)
				if err != nil {
					logp.Err("Unable to configure the indexer. %s, %v", name, indexerConfig)
					continue
				}
				indexers = append(indexers, i)
			default:
				logp.Warn("Unable to find indexer plugin %s", name)
			}
		}
	}

	return &Indexers{
		indexers: indexers,
	}
}

func (i *Indexers) GetIndices(task *ecs_task_metadata.TaskMetadata) []MetadataIndex {
	var metadata []MetadataIndex

	for _, indexer := range i.indexers {
		for _, m := range indexer.GetIndices(task) {
			metadata = append(metadata, m)
		}
	}
	return metadata
}

type ContainerNameIndexer struct {
	formatter *MetadataFormatter
}

func NewContainerNameIndexer(cfg *common.Config) (*ContainerNameIndexer, error) {
	formatter := NewMetadataFormatter()
	return &ContainerNameIndexer{
		formatter: formatter,
	}, nil
}

func (i *ContainerNameIndexer) GetIndices(task *ecs_task_metadata.TaskMetadata) []MetadataIndex {
	var metadata []MetadataIndex

	for _, container := range task.Containers {
		meta := i.formatter.FormatContainerMetadata(task, container.Name)
		metadata = append(metadata, MetadataIndex{
			Index: container.Name,
			Data:  meta,
		})
	}

	return metadata
}
