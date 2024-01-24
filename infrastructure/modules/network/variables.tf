variable "project_name" {
  type        = string
  description = "project name"
}

variable "vpc_cidr_block" {
  type        = string
  description = "vpc cidr block"
}

variable "public_subnets" {
  type = map(object({
    availability_zone = string
  }))
  description = "public subnet settings"
}

variable "private_subnets" {
  type = map(object({
    availability_zone = string
  }))
  description = "private subnet settings"
}
