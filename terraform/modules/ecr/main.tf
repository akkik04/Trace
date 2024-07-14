resource "aws_ecr_repository" "repos" {
  count = length(var.repository_names)
  name  = var.repository_names[count.index]
}