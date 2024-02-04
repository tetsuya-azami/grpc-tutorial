variable "project_name" {
  type        = string
  description = "The name of the project"
}

variable "vpc_id" {
  type        = string
  description = "vpc id for the vpc endpoint"
}

variable "vpc_cidr_block" {
  type        = string
  description = "The CIDR block of the VPC"
}

variable "backend_container_subnet_ids" {
  type        = list(string)
  description = "The subnet ids for the backend containers"
}

variable "route_table_id" {
  type        = string
  description = "The route table id for the s3 vpc endpoint"
}
