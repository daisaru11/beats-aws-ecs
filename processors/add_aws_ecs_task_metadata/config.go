package add_aws_ecs_task_metadata

import (
	"os"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
)

type addAwsEcsTaskMetadataConfig struct {
	EndpointTimeout    time.Duration `config:"endpoint_timeout"` // Amount of time to wait for responses from the metadata services.
	EndpointMaxRetries int           `config:"endpoint_max_retries"`
	EcsTaskMetadataUri string        `config:"ecs_task_metadata_uri"`
	Indexers           IndexerConfig `config:"indexers"`
	Matchers           MatcherConfig `config:"matchers"`
}

type IndexerConfig []map[string]common.Config
type MatcherConfig []map[string]common.Config

func defaultAddAwsEcsTaskMetadataConfig() addAwsEcsTaskMetadataConfig {
	uri := strings.TrimRight(os.Getenv("ECS_CONTAINER_METADATA_URI"), "/") + "/task"
	return addAwsEcsTaskMetadataConfig{
		EndpointTimeout:    3 * time.Second,
		EndpointMaxRetries: 3,
		EcsTaskMetadataUri: uri,
	}
}
