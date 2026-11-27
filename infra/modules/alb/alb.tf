# ホストゾーンは Terraform で作らず参照する
data "aws_route53_zone" "public" {
  name         = var.zone_name
  private_zone = false
}

# ACM はドメイン名から証明書を探して ARN を返す
data "aws_acm_certificate" "ec" {
  domain      = var.domain_name
  statuses    = ["ISSUED"]
  most_recent = true
}

resource "aws_alb" "public_alb" {
  name                       = "${local.name_prefix}-public-alb"
  subnets                    = var.public_subnet_ids
  security_groups            = [var.security_group_id]
  internal                   = false
  enable_deletion_protection = false
  load_balancer_type         = "application"

  access_logs {
    bucket  = aws_s3_bucket.alb_access_logs.bucket
    prefix  = "${var.env}/ALB/${local.name_prefix}-public-alb"
    enabled = var.enable_access_log
  }

  tags = {
    Name = "${local.name_prefix}-public-alb"
  }

  # バケットポリシーより先に ALB を作ると、アクセスログの設定で Access Denied になる
  depends_on = [aws_s3_bucket_policy.alb_access_logs]
}

# 443: 既定は web、/api/* は api
resource "aws_alb_listener" "public_alb_listener_https" {
  load_balancer_arn = aws_alb.public_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = data.aws_acm_certificate.ec.arn

  default_action {
    type = "forward"

    # TODO: フロント実装後にwebに変更
    # target_group_arn = aws_alb_target_group.web.arn
    target_group_arn = aws_alb_target_group.api.arn
  }
}

resource "aws_alb_listener_rule" "api" {
  listener_arn = aws_alb_listener.public_alb_listener_https.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.api.arn
  }

  condition {
    path_pattern {
      values = ["/api/*"]
    }
  }
}

# 80: 443 へ 301
resource "aws_alb_listener" "public_alb_listener_http" {
  load_balancer_arn = aws_alb.public_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}
