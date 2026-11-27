variable "env" {
  type    = string
  default = "prod"
}

variable "service_name" {
  type    = string
  default = "go-react-ec"
}

data "aws_caller_identity" "current" {}
