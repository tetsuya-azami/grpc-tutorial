resource "aws_lb" "main" {
  name               = "${var.project_name}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  tags = {
    Name = "${var.project_name}-alb"
  }
}

resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = var.certificate_arn
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.back_containers.arn
  }
}

resource "aws_lb_target_group" "back_containers" {
  name             = "${var.project_name}-lb-tg"
  port             = 8080
  protocol         = "HTTP"
  protocol_version = "GRPC"
  target_type      = "ip"
  vpc_id           = var.vpc_id

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "alb" {
  name   = "${var.project_name}-alb-sg"
  vpc_id = var.vpc_id

  description = "security group for ALB"
  tags = {
    Name = "${var.project_name}-alb-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "alb_security_group_ingress_rule" {
  security_group_id = aws_security_group.alb.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"

  description = "allow https access from anywhere"
  tags = {
    Name = "${var.project_name}-alb-security-group-igress-rule"
  }
}

resource "aws_vpc_security_group_egress_rule" "alb_security_group_egress_rule" {
  security_group_id = aws_security_group.alb.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"

  description = "allow all access to anywhere"
  tags = {
    Name = "${var.project_name}-alb-security-group-igress-rule"
  }
}
