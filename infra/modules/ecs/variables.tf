variable "env" {
  type = string
}

variable "service_name" {
  type = string
}

locals {
  name_prefix = "${var.service_name}-${var.env}"

  # ECR・タスクロール・ロググループをこの単位で作る
  apps = toset(["api", "web"])
}
