variable "env" {
  type = string
}

variable "service_name" {
  type = string
}

variable "vpc_cidr_block" {
  type = string
}

# キーは AZ の短縮名（1a / 1c）。名前タグの末尾に使う
variable "public_subnets" {
  type = map(object({
    cidr_block        = string
    availability_zone = string
  }))
}

variable "private_app_subnets" {
  type = map(object({
    cidr_block        = string
    availability_zone = string
  }))
}

variable "private_db_subnets" {
  type = map(object({
    cidr_block        = string
    availability_zone = string
  }))
}

# NAT Gateway を置く public サブネットのキー（public_subnets のキー）
variable "nat_gateway_subnet_key" {
  type = string

  validation {
    condition     = contains(keys(var.public_subnets), var.nat_gateway_subnet_key)
    error_message = "nat_gateway_subnet_key は public_subnets のキーのどれかにしてください。"
  }
}

locals {
  name_prefix = "${var.env}-${var.service_name}"
}
