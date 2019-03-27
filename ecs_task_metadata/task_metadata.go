package ecs_task_metadata

import (
	"encoding/json"
	"fmt"
	"time"
)

type TaskMetadata struct {
	Cluster       string
	TaskARN       string
	Family        string
	Revision      string
	DesiredStatus string `json:",omitempty"`
	KnownStatus   string
	Containers    []ContainerMetadata `json:",omitempty"`
	Limits        Limits              `json:",omitempty"`
}

type ContainerMetadata struct {
	ID            string `json:"DockerId"`
	Name          string
	DockerName    string
	Image         string
	ImageID       string
	Ports         []Port            `json:",omitempty"`
	Labels        map[string]string `json:",omitempty"`
	DesiredStatus string
	KnownStatus   string
	ExitCode      int `json:",omitempty"`
	Limits        Limits
	CreatedAt     time.Time `json:",omitempty"`
	StartedAt     time.Time `json:",omitempty"`
	FinishedAt    time.Time `json:",omitempty"`
	Type          string
	Health        HealthStatus `json:"health,omitempty"`
	Networks      []Network    `json:",omitempty"`
}

type HealthStatus struct {
	Status   string     `json:"status,omitempty"`
	Since    *time.Time `json:"statusSince,omitempty"`
	ExitCode int        `json:"exitCode,omitempty"`
	Output   string     `json:"output,omitempty"`
}

type Limits struct {
	CPU    uint
	Memory uint
}

type Port struct {
	ContainerPort uint16
	Protocol      string
	HostPort      uint16 `json:",omitempty"`
}

type Network struct {
	NetworkMode   string   `json:"NetworkMode,omitempty"`
	IPv4Addresses []string `json:"IPv4Addresses,omitempty"`
	IPv6Addresses []string `json:"IPv6Addresses,omitempty"`
}

func ParseTaskMetadata(data []byte) (*TaskMetadata, error) {
	var m TaskMetadata
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("Unable to parse TaskMetadata response: %v", err)
	}

	return &m, nil
}
