variable "env" {
  type    = string
  default = "prod"
}

variable "service_name" {
  type    = string
  default = "sample-ec"
}

variable "zone_name" {
  type    = string
  default = "geek-teru.com"
}

data "aws_caller_identity" "current" {}
