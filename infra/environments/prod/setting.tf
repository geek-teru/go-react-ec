# ----------------------------------------
# Provider Settings
# ----------------------------------------
provider "aws" {
  region = "ap-northeast-1"

  default_tags {
    tags = {
      Project     = var.service_name
      Environment = var.env
      ManagedBy   = "terraform"
    }
  }
}

# ----------------------------------------
# Terraform Settings
# ----------------------------------------
terraform {
  required_version = "~> 1.17.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {
    # bucket は「[本番]Terraform 実行環境の準備」で決める（既存のバケットを使うならここを書き換える）
    bucket       = "go-react-ec-prod-tfstate"
    key          = "go-react-ec/prod/terraform.tfstate"
    region       = "ap-northeast-1"
    encrypt      = true
    use_lockfile = true
  }
}
