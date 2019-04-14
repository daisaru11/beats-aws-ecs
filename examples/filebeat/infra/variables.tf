variable aws_region {
  default = "ap-northeast-1"
}

variable "prefix" {
  default = "aws-beats-ecs-example"
}

variable vpc_cidr_block {
  default = "10.0.0.0/16"
}

variable subnet_public_1_cidrblock {
  default = "10.0.1.0/24"
}

variable subnet_public_2_cidrblock {
  default = "10.0.2.0/24"
}

variable subnet_private_1_cidrblock {
  default = "10.0.3.0/24"
}

variable subnet_private_2_cidrblock {
  default = "10.0.4.0/24"
}

variable "task_desired_count" {
  default = 2
}

locals {
  availability_zone_1 = "${var.aws_region}a"
  availability_zone_2 = "${var.aws_region}c"
}
