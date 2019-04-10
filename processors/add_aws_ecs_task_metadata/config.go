package add_aws_ecs_task_metadata

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type addAwsEcsTaskMetadataConfig struct {
	EndpointTimeout    time.Duration `config:"endpoint_timeout"` // Amount of time to wait for responses from the metadata services.
	EndpointMaxRetries int           `config:"endpoint_max_retries"`
	EndpointBaseUrl    string        `config:"endpoint_base_url"`
	EndpointVersion    string        `config:"endpoint_version"`
	Indexers           IndexerConfig `config:"indexers"`
	Matchers           MatcherConfig `config:"matchers"`
}

type IndexerConfig []map[string]common.Config
type MatcherConfig []map[string]common.Config

func defaultAddAwsEcsTaskMetadataConfig() addAwsEcsTaskMetadataConfig {
	return addAwsEcsTaskMetadataConfig{
		EndpointTimeout:    3 * time.Second,
		EndpointMaxRetries: 3,
		EndpointVersion:    "v3",
	}
}
