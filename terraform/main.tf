provider "aws" {
  region = var.aws_region
}

module "aws_s3" {
  source           = "./modules/aws-s3"
  main_bucket_name = var.main_bucket_name
}

module "aws_ecr" {
  source = "./modules/aws-ecr"
}

module "aws_eks" {
  source = "./modules/aws-eks"
}
