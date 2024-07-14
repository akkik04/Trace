variable "main_bucket_name" {
  description = "The name for the S3 bucket"
  type        = string
}

variable "lambda_func_arn" {
  description = "ARN of the Lambda function to be triggered"
  type        = string
}

variable "lambda_func_name" {
  description = "Name of the Lambda function to be triggered"
  type        = string
}
