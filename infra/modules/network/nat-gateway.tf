# 外向き通信は NAT Gateway 1台（冗長化しない）
resource "aws_eip" "nat" {
  domain = "vpc"

  tags = {
    Name = "${local.name_prefix}-nat-eip"
  }

  depends_on = [aws_internet_gateway.ec]
}

resource "aws_nat_gateway" "egress" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public[var.nat_gateway_subnet_key].id

  tags = {
    Name = "${local.name_prefix}-nat-${var.nat_gateway_subnet_key}"
  }

  depends_on = [aws_internet_gateway.ec]
}
