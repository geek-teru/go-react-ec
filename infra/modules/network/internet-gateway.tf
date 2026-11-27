resource "aws_internet_gateway" "ec" {
  vpc_id = aws_vpc.ec.id

  tags = {
    Name = "${local.name_prefix}-igw"
  }
}
