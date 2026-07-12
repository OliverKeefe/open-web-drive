resource "aws_instance" "app_server" {
  ami = ""
  instance_type = "t2.micro"

  tags = {
    Name = "StagingInstance"
  }
}