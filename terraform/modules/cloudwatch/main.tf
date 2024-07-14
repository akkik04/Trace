resource "aws_cloudwatch_log_group" "log_group" {
  name              = var.log_group_name
  retention_in_days = 7
}

resource "aws_cloudwatch_log_stream" "example" {
  name           = var.log_stream_name
  log_group_name = aws_cloudwatch_log_group.log_group.name
}
