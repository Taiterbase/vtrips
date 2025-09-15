variable "name_prefix" { type = string }
variable "aws_region" { type = string }
variable "vpc_cidr" { type = string }
variable "azs" { type = list(string) }
variable "public_subnets" { type = list(string) }
variable "private_subnets" { type = list(string) }
variable "enable_nat_gateway" { type = bool }
variable "single_nat_gateway" { type = bool }
variable "cluster_name" { type = string }
variable "cluster_version" { type = string }
variable "node_min_size" { type = number }
variable "node_max_size" { type = number }
variable "node_desired_size" { type = number }
variable "node_instance_types" { type = list(string) }
variable "node_capacity_type" { type = string }
variable "tags" { type = map(string) }
variable "admin_role_name" { type = string }
variable "admin_role_trusted_principal_arns" { type = list(string) }


