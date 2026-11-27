variable "env" {
  type = string
}

variable "service_name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "public_subnet_ids" {
  type = list(string)
}

variable "security_group_id" {
  type = string
}

# Route 53 のホストゾーン名（例: example.com）。ec.<ゾーン名> で公開する
variable "zone_name" {
  type = string
}

# 公開するドメイン名（例: sample-ec.geek-teru.com）。zone_name のゾーン内の名前にする
variable "domain_name" {
  type = string
}

variable "enable_access_log" {
  type = bool
}

locals {
  name_prefix = "${var.env}-${var.service_name}"
}
