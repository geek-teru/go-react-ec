# ----------------------------------------
# Network
# ----------------------------------------
module "network" {
  source = "../../modules/network"

  env          = var.env
  service_name = var.service_name

  vpc_cidr_block = "10.30.0.0/16"

  public_subnets = {
    "1a" = { cidr_block = "10.30.0.0/24", availability_zone = "ap-northeast-1a" }
    "1c" = { cidr_block = "10.30.1.0/24", availability_zone = "ap-northeast-1c" }
  }

  private_app_subnets = {
    "1a" = { cidr_block = "10.30.10.0/24", availability_zone = "ap-northeast-1a" }
    "1c" = { cidr_block = "10.30.11.0/24", availability_zone = "ap-northeast-1c" }
  }

  private_db_subnets = {
    "1a" = { cidr_block = "10.30.20.0/24", availability_zone = "ap-northeast-1a" }
    "1c" = { cidr_block = "10.30.21.0/24", availability_zone = "ap-northeast-1c" }
  }

  nat_gateway_subnet_key = "1a"
}

# ----------------------------------------
# ECS
# ----------------------------------------
module "ecs" {
  source = "../../modules/ecs"

  env          = var.env
  service_name = var.service_name
}
