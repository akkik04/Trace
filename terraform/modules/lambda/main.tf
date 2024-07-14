# IAM Role for Lambda function.
resource "aws_iam_role" "lambda_exec_role" {
  name = var.lambda_exec_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com",
        },
      },
    ],
  })
}

# IAM Policy for Lambda Role
resource "aws_iam_policy" "lambda_policy" {
  name = var.lambda_policy_name

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
        ],
        Effect   = "Allow",
        Resource = "*",
      },
      {
        Action   = "s3:*",
        Effect   = "Allow",
        Resource = "*",
      },
    ],
  })
}

# Attach the policy to the Lambda role.
resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

# Create a CloudWatch Log Group for this Lambda Function.
resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name              = "/aws/lambda/${var.lambda_func_name}"
  retention_in_days = 14 # You can adjust the retention period as needed
}

# Create the inline Lambda code file and package it into a ZIP file.
data "archive_file" "lambda_zip" {
  type        = "zip"
  output_path = "${path.module}/lambda_function_payload.zip"

  source {
    content  = <<EOF
def lambda_handler(event, context):
    # Your Lambda function code here
    return {
        'statusCode': 200,
        'body': 'Hello, World!'
    }
EOF
    filename = "lambda_function.py"
  }
}

# Lambda function
resource "aws_lambda_function" "lambda_func" {
  function_name    = var.lambda_func_name
  runtime          = "python3.12"
  role             = aws_iam_role.lambda_exec_role.arn
  handler          = "index.lambda_handler"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path)
}
