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
  description = "The .zip filename of the lambda function"
  type        = string
}

variable "lambda_handler" {
  description = "The handler for the lambda function"
  type        = string
}