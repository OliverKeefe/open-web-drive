terraform {
  required_version = ">= 1.15.8"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  endpoints {
    sqs     = "http://localhost:4566"
    s3      = "http://localhost:4566"
    iam     = "http://localhost:4566"
    sts     = "http://localhost:4566"
    ec2     = "http://localhost:4566"
    route53 = "http://localhost:4566"
  }

  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

module "shared" {
  source = "../../../shared"
}

module "networking" {
  source       = "../../../modules/aws/networking"
  vpc_cidr     = var.vpc_cidr
  environment  = var.environment
  project_name = module.shared.project_name
  admin_cidrs  = ["0.0.0.0/0"]
  tags         = module.shared.common_tags
}