output "target_group_arn" {
  value       = aws_lb_target_group.back_containers.arn
  description = "target group arn"
}

output "security_group_id" {
  value       = aws_security_group.alb.id
  description = "security group id"
}
