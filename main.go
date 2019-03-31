package main

import (
	ecs "github.com/daisaru11/beats-aws-ecs/processors/add_aws_ecs_task_metadata"
	"github.com/elastic/beats/libbeat/plugin"
	"github.com/elastic/beats/libbeat/processors"
)

var Bundle = plugin.Bundle(
	processors.Plugin("add_aws_ecs_task_metadata", ecs.NewAddAwsEcsTaskMetadata),
)
