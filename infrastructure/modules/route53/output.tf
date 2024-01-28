output "name_servers" {
  value       = aws_route53_zone.main.name_servers
  description = "name servers"
}

output "certificate_arn" {
  value       = aws_acm_certificate.main.arn
  description = "certificate arn"
}
