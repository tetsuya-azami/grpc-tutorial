
# ecr repository
data "aws_ecr_repository" "back_container" {
  name = "${var.project_name}-back-container"
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = var.vpc_id
  service_name      = "com.amazonaws.ap-northeast-1.s3"
  vpc_endpoint_type = "Gateway"

  tags = {
    Name = "${var.project_name}-s3-vpc-endpoint"
  }
}

resource "aws_vpc_endpoint_route_table_association" "private_s3" {
  vpc_endpoint_id = aws_vpc_endpoint.s3.id
  route_table_id  = var.route_table_id
}

resource "aws_security_group" "vpc_endpoint" {
  name        = "${var.project_name}-vpc-endpoint-sg"
  description = "allow inbound traffic from backend containers"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.project_name}-vpc-endpoint-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "vpc_endpoint" {
  security_group_id = aws_security_group.vpc_endpoint.id
  cidr_ipv4         = var.vpc_cidr_block
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
}

resource "aws_vpc_security_group_egress_rule" "vpc_endpoint" {
  security_group_id = aws_security_group.vpc_endpoint.id
  cidr_ipv4         = var.vpc_cidr_block
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
}

resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.ap-northeast-1.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.backend_container_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoint.id]
  private_dns_enabled = true

  tags = {
    Name = "${var.project_name}-ecr-dkr-vpc-endpoint"
  }
}

resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.ap-northeast-1.ecr.api"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.backend_container_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoint.id]
  private_dns_enabled = true

  tags = {
    Name = "${var.project_name}-ecr-api-vpc-endpoint"
  }
}
