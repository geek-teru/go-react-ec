# S3 Buckets
resource "aws_s3_bucket" "alb_access_logs" {
  bucket              = "${local.name_prefix}-alb-access-logs"
  object_lock_enabled = true
  force_destroy       = true

  tags = {
    Name = "${local.name_prefix}-alb-access-logs"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "server_access_logs" {
  bucket = aws_s3_bucket.alb_access_logs.bucket
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256" // SSE-S3
    }
  }
}

resource "aws_s3_bucket_public_access_block" "system_logs_conf" {
  bucket                  = aws_s3_bucket.alb_access_logs.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

data "aws_elb_service_account" "current" {}

# ELB がアクセスログを書き込めるようにする
resource "aws_s3_bucket_policy" "alb_access_logs" {
  bucket = aws_s3_bucket.alb_access_logs.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = { AWS = data.aws_elb_service_account.current.arn }
        Action    = "s3:PutObject"
        Resource  = "${aws_s3_bucket.alb_access_logs.arn}/*"
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.system_logs_conf]
}
