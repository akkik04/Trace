variable "aws_region" {
  description = "The AWS region to deploy to"
  type        = string
}

variable "main_bucket_name" {
  description = "The name of the main S3 bucket"
  type        = string
}

variable "lambda_func_arn" {
  description = "The ARN of the Lambda function to be triggered"
  type        = string
}

variable "log_group_name" {
  description = "The name of the log group to create"
  type        = string
}

variable "log_stream_name" {
  description = "The name of the log stream to create"
  type        = string
}

variable "lambda_func_name" {
  description = "The name of the lambda function"
  type        = string
}

variable "lambda_exec_role_name" {
  description = "The name of the IAM role for the lambda function"
  type        = string
}

variable "lambda_policy_name" {
  description = "The name of the IAM policy for the lambda function"
  type        = string
}

variable "lambda_filename" {
  description = "The filename of the lambda function"
  type        = string
}

variable "lambda_handler" {
  description = "The handler for the lambda function"
  type        = string
}

variable "repository_names" {
  type        = list(string)
  description = "List of ECR repository names"

}

variable "ecs_cluster_name" {
  description = "The name of the ECS cluster."
  type        = string
}

variable "ecs_service_name" {
  description = "The name of the ECS service."
  type        = string
}

variable "ecs_task_definition_family_name" {
  description = "The family name of the ECS task definition."
  type        = string
}

variable "collector_image_uri" {
  description = "The URI of the collector microservice's Docker image within AWS ECR."
  type        = string
}

variable "ingestor_image_uri" {
  description = "The URI of the ingestor microservice's Docker image within AWS ECR."
  type        = string
}

variable "ecs_security_group_name" {
  description = "The name of the security group to create for ECS"
  type        = string
}

variable "vpc_name" {
  description = "The name of the VPC where ECS resources will be created."
  type        = string
}

variable "ecs_vpc_cidr_block" {
  description = "The CIDR block for the VPC."
  type        = string
}

variable "ecs_subnet_name" {
  description = "The name of the subnet where ECS resources will be created."
  type        = string
}

variable "ecs_subnet_cidr_block" {
  description = "The CIDR block for the subnet."
  type        = string
}

variable "ecs_route_table_cidr_block" {
  description = "The CIDR block for the route table."
  type        = string
}

variable "collector_image_digest" {
  description = "The digest of the collector microservice's Docker image within AWS ECR."
  type        = string
}

variable "ingestor_image_digest" {
  description = "The digest of the ingestor microservice's Docker image within AWS ECR."
  type        = string
}