output "name_servers" {
  value       = data.aws_route53_zone.main.name_servers
  description = "name servers"
}

output "certificate_arn" {
  value       = data.aws_acm_certificate.main.arn
  description = "certificate arn"
}
