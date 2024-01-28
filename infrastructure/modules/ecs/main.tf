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
    }
  ])

  execution_role_arn = data.aws_iam_role.my_ecs_task_execution_role.arn
  task_role_arn      = data.aws_iam_role.my_ecs_task_role.arn

  tags = {
    Name = "${var.project_name}-back-container-task-definition"
  }
}

resource "aws_ecr_repository" "back_container" {
  name                 = "${var.project_name}-back-container"
  image_tag_mutability = "IMMUTABLE"

  tags = {
    Name = "${var.project_name}-back-container-repository"
  }
}

data "aws_iam_role" "my_ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"
}

data "aws_iam_role" "my_ecs_task_role" {
  name = "ecsTaskRole"
}
