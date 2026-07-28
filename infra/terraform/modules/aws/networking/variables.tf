
variable "vpc_cidr" {
  type = string
  default = "10.0.0.0/16"
  description = "CIDR block for the VPC."
}

# The environment name used for resource naming and tagging.
# Validation block ensures only valid values are accepted, thus
# preventing typos like "stging" or "text" creating misspelled
# resources.
variable "environment" {
  type = string
  description = "Deployment environment."
  validation {
    condition = contains(["dev", "test", "staging", "prod"], var.environment)
    error_message = "Must be on of: dev, test, staging, prod."
  }
}

variable "tags" {
  type = map(string)
  description = "Tags to apply to all resources."
  default = {}
}

variable "project_name" {
  type = string
  description = "Project name used as prefix for resource names."
}

variable "admin_cidrs" {
  type = list(string)
  description = "CIDR blocks assigned to admins."
  default = ["10.0.0.0/2"]
}