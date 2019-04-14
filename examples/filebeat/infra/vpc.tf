resource "aws_vpc" "vpc" {
  cidr_block           = "${var.vpc_cidr_block}"
  instance_tenancy     = "default"
  enable_dns_support   = "true"
  enable_dns_hostnames = "true"

  tags {
    Name = "${var.prefix}-vpc"
  }
}

resource "aws_default_route_table" "default" {
  default_route_table_id = "${aws_vpc.vpc.default_route_table_id}"

  tags {
    Name = "${var.prefix}-default-route"
  }
}

resource "aws_subnet" "public_1" {
  vpc_id                  = "${aws_vpc.vpc.id}"
  cidr_block              = "${var.subnet_public_1_cidrblock}"
  availability_zone       = "${local.availability_zone_1}"
  map_public_ip_on_launch = true

  tags {
    Name = "${var.prefix}-public-subnet-1"
  }
}

resource "aws_subnet" "public_2" {
  vpc_id                  = "${aws_vpc.vpc.id}"
  cidr_block              = "${var.subnet_public_2_cidrblock}"
  availability_zone       = "${local.availability_zone_2}"
  map_public_ip_on_launch = true

  tags {
    Name = "${var.prefix}-public-subnet-2"
  }
}

resource "aws_route_table_association" "public_1" {
  route_table_id = "${aws_route_table.public_route.id}"
  subnet_id      = "${aws_subnet.public_1.id}"
}

resource "aws_route_table_association" "public_2" {
  route_table_id = "${aws_route_table.public_route.id}"
  subnet_id      = "${aws_subnet.public_2.id}"
}

resource "aws_route_table" "public_route" {
  vpc_id = "${aws_vpc.vpc.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gateway.id}"
  }

  tags {
    Name = "${var.prefix}-public-route"
  }
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = "${aws_vpc.vpc.id}"
}

resource "aws_subnet" "private_1" {
  vpc_id                  = "${aws_vpc.vpc.id}"
  cidr_block              = "${var.subnet_private_1_cidrblock}"
  availability_zone       = "${local.availability_zone_1}"
  map_public_ip_on_launch = false

  tags {
    "Name" = "${var.prefix}-private-subnet-1"
  }
}

resource "aws_subnet" "private_2" {
  vpc_id                  = "${aws_vpc.vpc.id}"
  cidr_block              = "${var.subnet_private_2_cidrblock}"
  availability_zone       = "${local.availability_zone_2}"
  map_public_ip_on_launch = false

  tags {
    "Name" = "${var.prefix}-private-subnet-2"
  }
}

resource "aws_route_table" "private_route" {
  vpc_id = "${aws_vpc.vpc.id}"

  tags {
    Name = "${var.prefix}-private-route"
  }
}

resource "aws_route" "private_to_nat" {
  route_table_id         = "${aws_route_table.private_route.id}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.nat_1.id}"
}

resource "aws_route_table_association" "private_1" {
  route_table_id = "${aws_route_table.private_route.id}"
  subnet_id      = "${aws_subnet.private_1.id}"
}

resource "aws_route_table_association" "private_2" {
  route_table_id = "${aws_route_table.private_route.id}"
  subnet_id      = "${aws_subnet.private_2.id}"
}

resource "aws_eip" "nat_1" {
  vpc = true
}

resource "aws_nat_gateway" "nat_1" {
  allocation_id = "${aws_eip.nat_1.id}"
  subnet_id     = "${aws_subnet.public_1.id}"

  tags {
    Name = "${var.prefix}-nat-1"
  }
}
