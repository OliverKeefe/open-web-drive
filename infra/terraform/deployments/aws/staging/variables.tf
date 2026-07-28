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
  default = "staging"
  validation {
    condition = contains(["dev", "test", "staging", "prod"], var.environment)
    error_message = "Must be one of: dev, test, staging, prod."
  }
}



variable "admin_cidrs" {
  type        = list(string)
  description = "CIDR blocks allowed SSH access."
}

variable "k8s_deployment_type" {
  type = string
  default = "eks"
  description = "Choose the k*s node type: 'eks (Managed) or 'ec2' (Self-Managed Nodes)."

  validation {
    condition = contains(["eks", "ec2"], var.k8s_deployment_type)
    error_message = "Allowed values are 'eks' or 'ec2'."
  }
}

variable "storage_backend" {
  type = string
  default = "s3"
  description = "Choose the object storage backend: 's3' or self managed 'minio'."

  validation {
    condition = contains(["s3", "minio"], var.storage_backend)
    error_message = "Allowed values are 's3' or 'minio'"
  }
}