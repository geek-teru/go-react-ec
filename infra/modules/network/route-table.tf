# ----------------------------------------
# public: 0.0.0.0/0 → インターネットゲートウェイ
# ----------------------------------------
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.ec.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.ec.id
  }

  tags = {
    Name = "${local.name_prefix}-public-rtb"
  }
}

resource "aws_route_table_association" "public" {
  for_each = aws_subnet.public

  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

# ----------------------------------------
# private-app: 0.0.0.0/0 → NAT Gateway
# ----------------------------------------
resource "aws_route_table" "private_app" {
  vpc_id = aws_vpc.ec.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.egress.id
  }

  tags = {
    Name = "${local.name_prefix}-private-app-rtb"
  }
}

resource "aws_route_table_association" "private_app" {
  for_each = aws_subnet.private_app

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private_app.id
}

# ----------------------------------------
# private-db: VPC 内（local）だけ
# ----------------------------------------
resource "aws_route_table" "private_db" {
  vpc_id = aws_vpc.ec.id

  tags = {
    Name = "${local.name_prefix}-private-db-rtb"
  }
}

resource "aws_route_table_association" "private_db" {
  for_each = aws_subnet.private_db

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private_db.id
}
