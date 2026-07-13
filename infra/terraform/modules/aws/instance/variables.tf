variable "aws_region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region in which your infrastructure is located."
}

variable "aws_instance_name" {
  type = string
  default = "ExampleInstance"
  description = "AWS instance name."
}