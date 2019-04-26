## Beats plugin for AWS ECS

Note: This is still a work in progress.

[Beats](https://github.com/elastic/beats) processor plugin for AWS ECS Task.

This plugin appends AWS ECS Task and container metadata to beats events.

```
{
  "@timestamp": "2019-04-14T16:17:49.374Z",
  ...
  "aws_ecs_task": {
    "task": {
      "cluster": "ecs-local-cluster",
      "arn": "arn:aws:ecs:us-west-2:111111111111:task/ecs-local-cluster/37e873f6-37b4-42a7-af47-eac7275c6152",
      "family": "esc-local-task-definition",
      "revision": "1"
    },
    "container": {
      "image_id": "sha256:031c45582fce6e8234175ed01cfea828a8f096e5b1ed3cdd41142d2a40244d27",
      "name": "test-nginx",
      "docker_name": "test-nginx",
      "image": "nginx:alpine"
    }
  },
  "message": "192.168.176.1 - - [14/Apr/2019:15:26:07 +0000] \"GET / HTTP/1.1\" 200 612 \"-\" \"curl/7.54.0\" \"-\"",
  ...
}
```

## Usage

### Filebeat

Add `add_aws_ecs_task_metadata` processor to your filebeat config.

```yaml:filebeat.yml
processors:
  - add_aws_ecs_task_metadata:
      indexers:
        - container_name:
      matchers:
        - container_name:
            name: nginx
```

Run a filebeat container as a sidecar with the container which outputs logs you want to collect. The logs is shared through volume.

```json:taskDefinition.json
{
  "family": "nginx",
  "taskRoleArn": "xxxx",
  "executionRoleArn": "xxxx",
  "networkMode": "awsvpc",
  "cpu": 256,
  "memory": 512,
  "requiresCompatibilities": ["FARGATE"],
  "volumes": [
    {
      "name": "nginx_logs"
    }
  ],
  "containerDefinitions": [
    {
      "essential": true,
      "image": "nginx:alpine",
      "name": "nginx",
      "portMappings": [
        {
          "containerPort": 80
        }
      ],
      "command": [
        "/bin/sh",
        "-c",
        "rm /var/log/nginx/*.log && exec nginx -g 'daemon off;'"
      ],
      "mountPoints": [
        {
          "sourceVolume": "nginx_logs",
          "containerPath": "/var/log/nginx"
        }
      ]
    },
    {
      "essential": true,
      "image": "daisaru11/beats-aws-ecs:filebeat",
      "name": "filebeat",
      "secrets": [
        {
          "name": "FILEBEAT_CONFIG",
          "valueFrom": "xxxx"
        }
      ],
      "command": [
        "/bin/sh",
        "-c",
        "echo \"$FILEBEAT_CONFIG\" > /tmp/filebeat.yml && filebeat --plugin beats-aws-ecs.so -e -c /tmp/filebeat.yml"
      ],
      "volumesFrom": [
        {
          "sourceContainer": "nginx"
        }
      ]
    }
  ]
}
```
