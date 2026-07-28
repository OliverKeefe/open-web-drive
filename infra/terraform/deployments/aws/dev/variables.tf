
variable "aws_region" {
  type = string
  default = "us-east-1"
  description = "AWS region."
}

variable "vpc_cidr" {
  type = string
  default = "10.0.0.0/16"
}

variable "environment" {
  type = string
  default = "dev"
  validation {
    condition = contains(["dev", "test", "staging", "prod"], var.environment)
    error_message = "Must be one of: dev, test, staging, prod."
  }
}