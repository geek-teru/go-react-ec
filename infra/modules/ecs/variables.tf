variable "env" {
  type = string
}

variable "service_name" {
  type = string
}

locals {
  name_prefix = "${var.env}-${var.service_name}"
}
