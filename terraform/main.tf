provider "aws" {
  region = var.aws_region
}

module "aws_s3" {
  source           = "./modules/s3"
  main_bucket_name = var.main_bucket_name
  lambda_func_arn  = module.aws_lambda.lambda_function_arn
  lambda_func_name = var.lambda_func_name
}

module "aws_ecr" {
  source           = "./modules/ecr"
  repository_names = var.repository_names
}

# module "aws_eks" {
#   source = "./modules/aws-eks"
# }

module "aws_cloudwatch" {
  source          = "./modules/cloudwatch"
  log_group_name  = var.log_group_name
  log_stream_name = var.log_stream_name
}

module "aws_lambda" {
  source                = "./modules/lambda"
  lambda_func_name      = var.lambda_func_name
  lambda_exec_role_name = var.lambda_exec_role_name
  lambda_policy_name    = var.lambda_policy_name
  lambda_filename       = var.lambda_filename
  lambda_handler        = var.lambda_handler
}