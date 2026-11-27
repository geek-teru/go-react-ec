# -----------------------------
# API
# -----------------------------
resource "aws_ecr_repository" "api" {
  name = "${local.name_prefix}/api"

  image_scanning_configuration {
    scan_on_push = true
  }
}

# 最新 10 世代だけ残し、それより古いイメージは消す
resource "aws_ecr_lifecycle_policy" "api" {
  repository = aws_ecr_repository.api.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# -----------------------------
# Web
# -----------------------------
resource "aws_ecr_repository" "web" {
  name = "${local.name_prefix}/web"

  image_scanning_configuration {
    scan_on_push = true
  }
}

# 最新 10 世代だけ残し、それより古いイメージは消す
resource "aws_ecr_lifecycle_policy" "web" {
  repository = aws_ecr_repository.web.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
