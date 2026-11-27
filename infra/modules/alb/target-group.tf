# ECS Fargate（awsvpc）なので target_type は ip
resource "aws_alb_target_group" "api" {
  name                 = "${local.name_prefix}-api"
  vpc_id               = var.vpc_id
  port                 = 8080
  protocol             = "HTTP"
  target_type          = "ip"
  deregistration_delay = 30

  health_check {
    path                = "/api/health"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  tags = {
    Name = "${local.name_prefix}-api"
  }
}
