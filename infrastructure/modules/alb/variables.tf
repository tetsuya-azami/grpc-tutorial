variable "project_name" {
  type        = string
  description = "project name"
}

variable "vpc_id" {
  type        = string
  description = "vpc id"
}

variable "public_subnet_ids" {
  type        = set(string)
  description = "public subnet ids"
}

variable "certificate_arn" {
  type        = string
  description = "certificate arn"
}
