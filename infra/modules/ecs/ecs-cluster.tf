# -----------------------------
# api/ web共通 ECS cluster
# -----------------------------
resource "aws_ecs_cluster" "app" {
  name = local.name_prefix

  # 費用を抑えるため Container Insights は使わない
  setting {
    name  = "containerInsights"
    value = "disabled"
  }
}

resource "aws_ecs_cluster_capacity_providers" "app" {
  cluster_name = aws_ecs_cluster.app.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
  }
}
