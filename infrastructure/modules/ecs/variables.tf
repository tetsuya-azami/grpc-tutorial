variable "project_name" {
  type        = string
  description = "project name"
}

variable "image_tag" {
  type        = string
  description = "image tag"
}

variable "target_group_arn" {
  type        = string
  description = "target group arn"
}

variable "vpc_id" {
  type        = string
  description = "vpc id"
}

variable "vpc_cidr_block" {
  type        = string
  description = "vpc cidr block"
}

variable "backend_container_subnet_ids" {
  type        = set(string)
  description = "backend container subnet ids"
}

variable "alb_security_group_id" {
  type        = string
  description = "alb security group id"
}

variable "route_table_id" {
  type        = string
  description = "route table id"
}
