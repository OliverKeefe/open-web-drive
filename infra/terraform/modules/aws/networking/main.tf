data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "this" {
  cidr_block = var.vpc_cidr
  enable_dns_support = true
  enable_dns_hostnames = true

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-vpc"
  })
}

resource "aws_subnet" "public" {
  count = 2
  vpc_id = aws_vpc.this.id
  cidr_block = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-public-${count.index}"
    Tier = "public"
  })
}

resource "aws_subnet" "private" {
  count = 2
  vpc_id = aws_vpc.this.id
  cidr_block = cidrsubnet(var.vpc_cidr, 8, count.index + 10) # uses index + 1 to avoid cidr overlap with public subnets.
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-private-${count.index}"
    Tier = "private"
  })
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-igw"
  })
}

# Private instances still need to pull Docker images, package updates etc...
# Thus, NAT gateway is required because the private instances don't have
# public IP addresses.
resource "aws_eip" "nat" {
  domain = "vpc"

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-nat-eip"
  })
}

resource "aws_nat_gateway" "this" {
allocation_id = aws_eip.nat.id
  subnet_id = aws_subnet.public[0].id

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-nat"
  })
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.this.id
}

resource "aws_route_table_association" "public" {
  count = 2
  subnet_id = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count = 2
  subnet_id = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private.id
}

# The default is only for testing / config purposes
# due to it allowing basically all inbound / outbound traffic
# on any ports.
resource "aws_security_group" "default" {
  name = "${var.project_name}-${var.environment}-default-sg"
  description = "Default Security Group for ${var.environment}"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    self = true
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-nat-default-sg"
  })
}

# Security Group for App server / EKS instance (depending on deployment type).
resource "aws_security_group" "app" {
  name = "${var.project_name}-${var.environment}-app-sg"
  vpc_id = aws_vpc.this.id

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # May need to restrict to LB Security Group at some point.
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = var.admin_cidrs
  }

  egress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 53
    to_port = 53
    protocol = "tcp"
    cidr_blocks = [aws_vpc.this.cidr_block]
  }
}

resource "aws_vpc_dhcp_options" "this" {
  domain_name = "${var.environment}.gestalt.internal"
  domain_name_servers = ["AmazonProvidedDNS"]

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-dhcp"
  })
}

resource "aws_vpc_dhcp_options_association" "this" {
  vpc_id = aws_vpc.this.id
  dhcp_options_id = aws_vpc_dhcp_options.this.id
}

resource "aws_route53_zone" "this" {
  name = "${var.environment}.open-web-drive.internal"

  vpc {
    vpc_id = aws_vpc.this.id
  }

  tags = merge(var.tags, {
    Name = "${var.project_name}-${var.environment}-dns"
  })
}