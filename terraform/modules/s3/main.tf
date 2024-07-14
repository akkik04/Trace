resource "aws_s3_bucket" "main_bucket" {
  bucket = var.main_bucket_name
}