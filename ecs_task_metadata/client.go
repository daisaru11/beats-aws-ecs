package ecs_task_metadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type TaskMetadataClient struct {
	client                 *http.Client
	taskMetadataEndpoint   string
	maxRetries             int
	durationBetweenRetries time.Duration
}

func NewTaskMetadataClient(c *http.Client) *TaskMetadataClient {
	return &TaskMetadataClient{
		client:                 c,
		taskMetadataEndpoint:   os.Getenv("ECS_CONTAINER_METADATA_URI") + "/task",
		maxRetries:             3,
		durationBetweenRetries: 1 * time.Second,
	}
}

func (c *TaskMetadataClient) GetTaskMetadata() (*TaskMetadata, error) {
	data, err := c.request(c.taskMetadataEndpoint)
	if err != nil {
		return nil, err
	}

	metadata, err := ParseTaskMetadata(data)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (c *TaskMetadataClient) request(endpoint string) ([]byte, error) {
	var resp []byte
	var err error
	for i := 0; i < c.maxRetries; i++ {
		resp, err = c.requestOnce(endpoint)
		if err == nil {
			return resp, nil
		}
		fmt.Fprintf(os.Stderr, "Attempt [%d/%d]: unable to get metadata response for from '%s': %v", i, c.maxRetries, endpoint, err)
		time.Sleep(c.durationBetweenRetries)
	}

	return nil, err
}

func (c *TaskMetadataClient) requestOnce(endpoint string) ([]byte, error) {
	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Unable to get response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Incorrect status code  %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response body: %v", err)
	}

	return body, nil
}
