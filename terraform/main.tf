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

module "aws_ecs" {
  source                          = "./modules/ecs"
  ecs_cluster_name                = var.ecs_cluster_name
  ecs_service_name                = var.ecs_service_name
  ecs_task_definition_family_name = var.ecs_task_definition_family_name
  collector_image_uri             = var.collector_image_uri
  ingestor_image_uri              = var.ingestor_image_uri
  ecs_security_group_name         = var.ecs_security_group_name
  vpc_name                        = var.vpc_name
  ecs_vpc_cidr_block              = var.ecs_vpc_cidr_block
  ecs_subnet_name                 = var.ecs_subnet_name
  ecs_subnet_cidr_block           = var.ecs_subnet_cidr_block
  ecs_route_table_cidr_block      = var.ecs_route_table_cidr_block
  collector_image_digest          = var.collector_image_digest
  ingestor_image_digest           = var.ingestor_image_digest
}
