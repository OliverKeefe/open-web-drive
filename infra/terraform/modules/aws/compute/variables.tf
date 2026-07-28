variable "name" {
  type = string
  description = "Instance name."
}

variable "instance_type" {
  type = string
  default = "t3.micro"
}

# Passed for consistency with interface contract.
variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  type = string
}

variable "security_group_id" {
  type = string
  default = ""
}

variable "environment" {
  type = string
}

variable "tags" {
  type = map(string)
  default = {}
}