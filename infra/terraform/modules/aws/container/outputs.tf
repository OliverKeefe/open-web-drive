output "cluster_id" {
  value = aws_eks_cluster.this.id
}

# Kubernetes API Server endpoint uris are formatted as follows:
# https://XXXXX.eks.amazonaws.com
output "cluster_endpoint" {
  value = aws_eks_cluster.this.endpoint
}

# Base64-encoded CA cert for the k8s cluster.
# Required for kubectl to verify the API server's ID.
output "cluster_ca_certificate" {
  value = aws_eks_cluster.this.certificate_authority[0].data
}