variable "project_name" {
  type        = string
  description = "project name"
}

variable "database_name" {
  type        = string
  description = "athena database name"
}

variable "source_s3_location" {
  type        = string
  description = "source s3 location which athena is querying from"
}
