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
}

module "shared" {
  source = "../../../shared"
}

module "networking" {
  source = "../../../modules/aws/networking"
  vpc_cidr = var.vpc_cidr
  environment = var.environment
  project_name = module.shared.project_name
  tags = module.shared.common_tags
}