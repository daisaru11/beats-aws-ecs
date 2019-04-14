resource "aws_cloudwatch_log_group" "ecs_task_logs" {
  name              = "${var.prefix}-ecs-tasks"
  retention_in_days = 1
}
