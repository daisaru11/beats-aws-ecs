locals {
  nginx_container_name    = "nginx"
  filebeat_container_name = "filebeat"
}

resource "aws_ecs_task_definition" "nginx" {
  family = "${var.prefix}-nginx"

  #  task_role_arn            = "${aws_iam_role.nginx_task.arn}"
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = "${aws_iam_role.nginx_task_execution.arn}"

  volume {
    name = "nginx_logs"
  }

  container_definitions = <<DEFINITION
[
  {
    "essential": true,
    "image": "nginx:alpine",
    "name": "${local.nginx_container_name}",
    "portMappings": [
      {
        "containerPort": 80
      }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "${aws_cloudwatch_log_group.ecs_task_logs.name}",
          "awslogs-region": "${var.aws_region}",
          "awslogs-stream-prefix": "nginx"
      }
    },
    "command": [
      "/bin/sh", "-c",
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
    "name": "${local.filebeat_container_name}",
    "secrets": [
      {
        "name": "FILEBEAT_CONFIG",
        "valueFrom": "${aws_ssm_parameter.filebeat_config.arn}"
      }
    ],
    "command": [
      "/bin/sh", "-c",
      "echo \"$FILEBEAT_CONFIG\" > /tmp/filebeat.yml && filebeat --plugin beats-aws-ecs.so -e -c /tmp/filebeat.yml"
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
          "awslogs-group": "${aws_cloudwatch_log_group.ecs_task_logs.name}",
          "awslogs-region": "${var.aws_region}",
          "awslogs-stream-prefix": "filebeat"
      }
    },
    "volumesFrom": [
      {
        "sourceContainer": "${local.nginx_container_name}"
      }
    ]
  }
]
DEFINITION
}

resource "aws_ssm_parameter" "filebeat_config" {
  name = "/data/${var.prefix}/filebeat-config"
  type = "String"

  value = <<VALUE
filebeat.inputs:
  - type: log
    paths:
      - /var/log/nginx/access.log

processors:
  - add_aws_ecs_task_metadata:
      indexers:
        - container_name:
      matchers:
        - container_name:
            name: ${local.nginx_container_name}
output.console:
  pretty: true
VALUE
}

resource "aws_ecs_service" "nginx" {
  name          = "${var.prefix}-nginx"
  cluster       = "${aws_ecs_cluster.cluster.id}"
  desired_count = "${var.task_desired_count}"
  launch_type   = "FARGATE"

  task_definition = "${aws_ecs_task_definition.nginx.arn}"

  load_balancer {
    target_group_arn = "${aws_alb_target_group.alb.arn}"
    container_name   = "${local.nginx_container_name}"
    container_port   = 80
  }

  network_configuration {
    subnets = ["${aws_subnet.private_1.id}", "${aws_subnet.private_2.id}"]

    security_groups = [
      "${aws_security_group.nginx_task_security_group.id}",
    ]
  }
}

resource "aws_security_group" "nginx_task_security_group" {
  name        = "${var.prefix}-nginx-task"
  description = "${var.prefix}-nginx-task"
  vpc_id      = "${aws_vpc.vpc.id}"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
