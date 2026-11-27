# ALB・NAT Gateway
resource "aws_subnet" "public" {
  for_each = var.public_subnets

  vpc_id            = aws_vpc.ec.id
  cidr_block        = each.value.cidr_block
  availability_zone = each.value.availability_zone

  tags = {
    Name = "${local.name_prefix}-public-${each.key}"
    Tier = "public"
  }
}

# ECS タスク
resource "aws_subnet" "private_app" {
  for_each = var.private_app_subnets

  vpc_id            = aws_vpc.ec.id
  cidr_block        = each.value.cidr_block
  availability_zone = each.value.availability_zone

  tags = {
    Name = "${local.name_prefix}-private-app-${each.key}"
    Tier = "private-app"
  }
}

# Aurora
resource "aws_subnet" "private_db" {
  for_each = var.private_db_subnets

  vpc_id            = aws_vpc.ec.id
  cidr_block        = each.value.cidr_block
  availability_zone = each.value.availability_zone

  tags = {
    Name = "${local.name_prefix}-private-db-${each.key}"
    Tier = "private-db"
  }
}
