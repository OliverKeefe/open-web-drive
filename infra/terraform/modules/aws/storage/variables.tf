variable "bucket_name" {
  type        = string                     # The type of the variable, in this case a string
  default     = "t2.micro"                 # Default value for the variable
  description = "The type of EC2 instance" # Description of what this variable represents
}

variable "environment" {
  type        = string
  description = "The environment for this bucket (e.g. dev, staging, prod). Used for tagging."

# Enforce specific variable names.
  validation {
    condition = contains(["dev", "test", "staging", "prod"], var.environment)
    error_message = "The environment name must be one of the following: dev, test, staging or prod."
  }
}

variable "tags" {
  type = map(string)
  description = "Additional tags to apply to the S3 bucket."
  default = {}
}