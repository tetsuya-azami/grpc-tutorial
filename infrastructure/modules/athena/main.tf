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

resource "aws_s3_bucket" "athena-results" {
  bucket = "${var.project_name}-athena-results"
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
