output "repository_url" {
  value       = data.aws_ecr_repository.back_container.repository_url
  description = "ECR repository url"
}
