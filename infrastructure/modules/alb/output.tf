output "name" {
  value       = aws_lb.main.name
  description = "alb name"
}

output "target_group_arn" {
  value       = aws_lb_target_group.back_containers.arn
  description = "target group arn"
}

output "security_group_id" {
  value       = aws_security_group.alb.id
  description = "security group id"
}

output "access_logs_bucket_name" {
  value       = aws_s3_bucket.elb_access_log.bucket
  description = "access logs bucket"
}
