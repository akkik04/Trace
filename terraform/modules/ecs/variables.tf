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
