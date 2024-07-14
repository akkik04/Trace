resource "aws_s3_bucket" "main_bucket" {
  bucket = var.main_bucket_name
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.main_bucket.id

  lambda_function {
    lambda_function_arn = var.lambda_func_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "logs/"
  }
}

resource "aws_lambda_permission" "allow_s3_invoke_lambda" {
  statement_id  = "AllowS3InvokeLambda"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_func_name
  principal     = "s3.amazonaws.com"

  source_arn = aws_s3_bucket.main_bucket.arn
}

