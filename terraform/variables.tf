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