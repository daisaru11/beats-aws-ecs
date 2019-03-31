package add_aws_ecs_task_metadata

import (
	"fmt"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Matcher interface {
	MetadataIndex(event common.MapStr) string
}

type Matchers struct {
	matchers []Matcher
}

func NewMatchers(config *addAwsEcsTaskMetadataConfig) *Matchers {
	matchers := []Matcher{}

	for _, matcherConfigs := range config.Matchers {
		for name, matcherConfig := range matcherConfigs {
			switch name {
			case "container_name":
				m, err := NewContainerNameMatcher(&matcherConfig)
				if err != nil {
					logp.Err("Unable to configure the matcher. %s, %v", name, matcherConfig)
					continue
				}
				matchers = append(matchers, m)
			default:
				logp.Warn("Unable to find indexer plugin %s", name)
			}
		}
	}

	return &Matchers{
		matchers: matchers,
	}
}

func (m *Matchers) MetadataIndex(event common.MapStr) string {
	for _, matcher := range m.matchers {
		index := matcher.MetadataIndex(event)
		if index != "" {
			return index
		}
	}

	return ""
}

type ContainerNameMatcher struct {
	Name string
}

func NewContainerNameMatcher(cfg *common.Config) (*ContainerNameMatcher, error) {
	config := struct {
		Name string `config:"name"`
	}{}

	err := cfg.Unpack(&config)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack container name matcher configuration: %s", err)
	}

	return &ContainerNameMatcher{
		Name: config.Name,
	}, nil
}

func (m *ContainerNameMatcher) MetadataIndex(event common.MapStr) string {
	return m.Name
}
