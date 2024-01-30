resource "aws_ecs_task_definition" "back_container" {
  family                   = "${var.project_name}-back-container"
  requires_compatibilities = ["FARGATE"]

  network_mode = "awsvpc"
  cpu          = 256
  memory       = 512

  container_definitions = jsonencode([
    {
      name      = "${var.project_name}-back-container"
      image     = "${aws_ecr_repository.back_container.repository_url}:${var.image_tag}"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${var.project_name}-back-container"
          "awslogs-region"        = "ap-northeast-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  execution_role_arn = data.aws_iam_role.my_ecs_task_execution_role.arn
  task_role_arn      = data.aws_iam_role.my_ecs_task_role.arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "ARM64"
  }

  tags = {
    Name = "${var.project_name}-back-container-task-definition"
  }
}

resource "aws_ecs_service" "back_containers" {
  name            = "${var.project_name}-back-containers"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.back_container.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = "${var.project_name}-back-container"
    container_port   = 8080
  }

  network_configuration {
    subnets          = var.backend_container_subnet_ids
    security_groups  = []
    assign_public_ip = false
  }

  tags = {
    Name = "${var.project_name}-back-containers"
  }
}

resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-ecs-cluster"

  tags = {
    Name = "${var.project_name}-ecs-cluster"
  }
}

resource "aws_security_group" "backend_containers" {
  name        = "${var.project_name}-backend-containers"
  description = "allow inbound traffic from ALB"
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.project_name}-backend-containers-sg"
  }
}

resource "aws_vpc_security_group_ingress_rule" "backend_containers" {
  security_group_id            = aws_security_group.backend_containers.id
  referenced_security_group_id = var.alb_security_group_id
  ip_protocol                  = "tcp"
  from_port                    = 8080
  to_port                      = 8080
}

resource "aws_vpc_security_group_egress_rule" "back_containers" {
  security_group_id = aws_security_group.backend_containers.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

# ecr repository
resource "aws_ecr_repository" "back_container" {
  name                 = "${var.project_name}-back-container"
  image_tag_mutability = "IMMUTABLE"

  tags = {
    Name = "${var.project_name}-back-container-repository"
  }
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

resource "aws_vpc_endpoint" "logs" {
  vpc_id              = var.vpc_id
  service_name        = "com.amazonaws.ap-northeast-1.logs"
  vpc_endpoint_type   = "Interface"
  subnet_ids          = var.backend_container_subnet_ids
  security_group_ids  = [aws_security_group.vpc_endpoint.id]
  private_dns_enabled = true

  tags = {
    Name = "${var.project_name}-logs-vpc-endpoint"
  }
}

# roles
data "aws_iam_role" "my_ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"
}

data "aws_iam_role" "my_ecs_task_role" {
  name = "ecsTaskRole"
}

# log group
resource "aws_cloudwatch_log_group" "ecs_container_log" {
  name              = "/ecs/${var.project_name}-back-container"
  retention_in_days = 7

  tags = {
    Name = "${var.project_name}-back-container-log-group"
  }
}
