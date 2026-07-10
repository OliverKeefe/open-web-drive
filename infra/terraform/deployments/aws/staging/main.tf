terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = aws_region
}

# Will deploy s3 by default unless specified.
module "s3_storage" {
  count = var.storage_backend == "s3" ? 1 : 0
  source = "../../../modules/aws/storage"

  bucket_name = "gestalt-staging-s3-bucket"
  environment = "staging"

  tags = {
    Vendor    = "AWS"
    ManagedBy = "Terraform"
    Project   = "GestaltApp"
  }
}

module "minio_storage" {
  count = var.storage_backend == "minio" ? 1 : 0
  source = "../../../modules/container/"

  environment = "staging"
}

