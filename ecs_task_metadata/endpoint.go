package ecs_task_metadata

import "os"

type EndpointVersion int

const (
	EndpointV2 EndpointVersion = iota
	EndpointV3
)

func GetDefaultEndpointBaseUrl(v EndpointVersion) string {
	switch v {
	case EndpointV2:
		return "http://169.254.170.2/v2"
	case EndpointV3:
		return os.Getenv("ECS_CONTAINER_METADATA_URI")
	default:
		return ""
	}
}

func GetTaskMetadataEndpointPath(v EndpointVersion, baseUrl string) string {
	switch v {
	case EndpointV2:
		return baseUrl + "/metadata"
	case EndpointV3:
		return baseUrl + "/task"
	default:
		return ""
	}
}

func GetContainerStatsEndpointPath(v EndpointVersion, baseUrl string) string {
	switch v {
	case EndpointV2:
		return baseUrl + "/stats"
	case EndpointV3:
		return baseUrl + "/stats"
	default:
		return ""
	}
}
