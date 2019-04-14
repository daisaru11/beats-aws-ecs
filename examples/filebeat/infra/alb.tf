resource "aws_alb" "alb" {
  name                       = "${var.prefix}-alb"
  security_groups            = ["${aws_security_group.alb.id}"]
  subnets                    = ["${aws_subnet.public_1.id}", "${aws_subnet.public_2.id}"]
  internal                   = false
  enable_deletion_protection = false
}

resource "aws_alb_target_group" "alb" {
  name        = "${var.prefix}-alb-tg"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = "${aws_vpc.vpc.id}"
  target_type = "ip"

  health_check {
    interval = 30
    path     = "/"
    port     = 80
    protocol = "HTTP"
    timeout  = 5
    matcher  = 200
  }
}

resource "aws_alb_listener" "http" {
  load_balancer_arn = "${aws_alb.alb.arn}"
  port              = "80"

  default_action {
    target_group_arn = "${aws_alb_target_group.alb.arn}"
    type             = "forward"
  }
}

resource "aws_security_group" "alb" {
  name        = "${var.prefix}-alb"
  vpc_id      = "${aws_vpc.vpc.id}"
  description = "${var.prefix}-alb"

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
