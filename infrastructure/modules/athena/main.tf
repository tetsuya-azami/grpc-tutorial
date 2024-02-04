resource "aws_athena_workgroup" "main" {
  name = "${var.project_name}-workgroup"

  configuration {
    enforce_workgroup_configuration    = true
    publish_cloudwatch_metrics_enabled = true
    result_configuration {
      output_location = "s3://${aws_s3_bucket.athena-results.bucket}"
    }
  }

  force_destroy = true
  tags = {
    Name = "main"
  }
}

resource "aws_athena_database" "main" {
  name   = var.database_name
  bucket = aws_s3_bucket.athena-results.id

  force_destroy = true
}

resource "terraform_data" "execute_create_table" {
  provisioner "local-exec" {
    command = <<-EOF
     aws athena start-query-execution \
     --profile admin \
     --work-group ${aws_athena_workgroup.main.name} \
     --query-execution-context Database=${aws_athena_database.main.name} \
     --query-string '${templatefile("${path.module}/sql/create_table.tftpl", { source_s3_location = var.source_s3_location })}'
    EOF
  }
}

resource "aws_s3_bucket" "athena-results" {
  bucket = "${var.project_name}-athena-results"

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
