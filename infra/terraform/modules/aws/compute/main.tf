data "aws_ami" "latest" {
  most_recent = true
  owners = ["amazon"]

  filter = {
    name = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

resource "aws_instance" "this" {
  ami = data.aws_ami.latest.id
  instance_type = var.instance_type
  subnet_id = var.subnet_id
  vpc_security_group_ids = [var.security_group_id]

  tags = merge(var.tags, {
    Name = var.name
    Environment = var.environment
  })
}

