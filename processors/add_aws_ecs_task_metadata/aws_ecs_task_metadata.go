package add_aws_ecs_task_metadata

import (
	"net"
	"net/http"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"

	"github.com/daisaru11/beats-aws-ecs/ecs_task_metadata"
	"github.com/daisaru11/beats-aws-ecs/processors/cache"
)

type addAwsEcsTaskMetadata struct {
	config   addAwsEcsTaskMetadataConfig
	indexers *Indexers
	matchers *Matchers
	updater  Updater
	cache    *cache.Cache
}

func NewAddAwsEcsTaskMetadata(c *common.Config) (processors.Processor, error) {
	config := defaultAddAwsEcsTaskMetadataConfig()
	err := c.Unpack(&config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unpack add_cloud_metadata config")
	}

	p := &addAwsEcsTaskMetadata{
		config: config,
		cache:  cache.NewCache(),
	}

	p.indexers = NewIndexers(&config)
	p.matchers = NewMatchers(&config)

	// Create HTTP client with our timeouts and keep-alive disabled.
	hc := http.Client{
		Timeout: config.EndpointTimeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout:   config.EndpointTimeout,
				KeepAlive: 0,
			}).DialContext,
		},
	}
	clientCfg := ecs_task_metadata.GetDefaultConfig()
	clientCfg.TaskMetadataEndpoint = config.EcsTaskMetadataUri
	clientCfg.MaxRetries = config.EndpointMaxRetries

	client := ecs_task_metadata.NewTaskMetadataClient(&hc, clientCfg)

	p.updater = NewOnceAtStartupUpdater(client, p.updateTaskMetadata)
	p.updater.Run()

	return p, nil
}

func (p *addAwsEcsTaskMetadata) Run(event *beat.Event) (*beat.Event, error) {
	index := p.matchers.MetadataIndex(event.Fields)
	logp.Debug("add_aws_ecs_task_metadata", "Matched index: %s", index)

	if index == "" {
		return event, nil
	}

	metadata := p.cache.Get(index)
	if metadata == nil {
		return event, nil
	}

	event.Fields.DeepUpdate(common.MapStr{
		"aws_ecs_task": metadata.Clone(),
	})

	return event, nil
}

func (p *addAwsEcsTaskMetadata) String() string {
	return "add_aws_ecs_task_metadata="
}

func (p *addAwsEcsTaskMetadata) updateTaskMetadata(task *ecs_task_metadata.TaskMetadata) {
	metadata := p.indexers.GetIndices(task)
	for _, m := range metadata {
		logp.Debug("add_aws_ecs_task_metadata", "Indexing metadata. index: %s, metadata: %v", m.Index, m.Data)
		p.cache.Set(m.Index, m.Data)
	}
}
