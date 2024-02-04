resource "aws_athena_workgroup" "main" {
  name = "${var.project_name}-workgroup"

  configuration {
    enforce_workgroup_configuration    = true
    publish_cloudwatch_metrics_enabled = true
    result_configuration {
      output_location = "s3://${aws_s3_bucket.athena-results.bucket}"
    }
  }
  tags = {
    Name = "main"
  }
}

resource "aws_athena_database" "main" {
  name   = var.database_name
  bucket = aws_s3_bucket.athena-results.id
}

resource "aws_athena_named_query" "create_table" {
  name      = "create_table"
  workgroup = aws_athena_workgroup.main.id
  database  = aws_athena_database.main.name
  query = templatefile(
    "${path.module}/sql/create_table.tftpl",
    { source_s3_location = var.source_s3_location }
  )
}

resource "aws_s3_bucket" "athena-results" {
  bucket        = "${var.project_name}-athena-results"
  force_destroy = true
  tags = {
    Name = "athena-results"
  }
}

resource "aws_s3_bucket_public_access_block" "elb_access_log" {
  bucket = aws_s3_bucket.athena-results.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
