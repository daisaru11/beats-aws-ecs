# resource "aws_iam_role" "nginx_task" {
#   name = "${var.prefix}-nginx-task"
# 
#   assume_role_policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Sid": "",
#       "Effect": "Allow",
#       "Principal": {
#         "Service": "ecs-tasks.amazonaws.com"
#       },
#       "Action": "sts:AssumeRole"
#     }
#   ]
# }
# EOF
# }
# 
# data "aws_iam_policy_document" "nginx_task" {
# }
# 
# resource "aws_iam_policy" "nginx_task" {
#   name        = "${var.prefix}-nginx-task"
#   path        = "/"
#   description = ""
#   policy      = "${data.aws_iam_policy_document.nginx_task.json}"
# }
# 
# resource "aws_iam_role_policy_attachment" "nginx_task" {
#   role       = "${aws_iam_role.nginx_task.name}"
#   policy_arn = "${aws_iam_policy.nginx_task.arn}"
# }

data "aws_iam_policy_document" "nginx_task_execution" {
  statement {
    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "ssm:GetParameters",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy" "nginx_task_execution" {
  name        = "${var.prefix}-nginx-task-execution"
  path        = "/"
  description = ""
  policy      = "${data.aws_iam_policy_document.nginx_task_execution.json}"
}

resource "aws_iam_role" "nginx_task_execution" {
  name = "${var.prefix}-nginx-task-execution"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "nginx_task_execution" {
  role       = "${aws_iam_role.nginx_task_execution.name}"
  policy_arn = "${aws_iam_policy.nginx_task_execution.arn}"
}
