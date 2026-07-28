# Elastic Kubernetes Service (EKS)

# AWS Account ID for IAM policy ARNs.
data "aws_caller_identity" "current" {

}

# EKS Cluster Control Plane IAM Role
resource "aws_iam_role" "eks" {
  name = "${var.project_name}-${var.environment}-eks-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
  })

  tags = var.tags
}