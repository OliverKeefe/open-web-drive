resource "aws_s3_bucket" "this" {
  bucket = var.bucket_name

  tags = merge(
    var.tags,
    {
      Name = var.bucket_name
      Environment = var.environment
    }
  )
}

output "bucket_id" {
  value = aws_s3_bucket.this.id
}