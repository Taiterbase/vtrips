variable "name_prefix" {
  description = "Name prefix for created resources"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
}

variable "azs" {
  description = "List of availability zones to use"
  type        = list(string)
}

variable "public_subnets" {
  description = "Public subnet CIDRs"
  type        = list(string)
}

variable "private_subnets" {
  description = "Private subnet CIDRs"
  type        = list(string)
}

variable "enable_nat_gateway" {
  description = "Whether to enable NAT gateways"
  type        = bool
  default     = true
}

variable "single_nat_gateway" {
  description = "Whether to use a single shared NAT gateway"
  type        = bool
  default     = true
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "cluster_version" {
  description = "EKS cluster Kubernetes version"
  type        = string
}

variable "node_min_size" {
  description = "Node group min size"
  type        = number
  default     = 1
}

variable "node_max_size" {
  description = "Node group max size"
  type        = number
  default     = 3
}

variable "node_desired_size" {
  description = "Node group desired size"
  type        = number
  default     = 1
}

variable "node_instance_types" {
  description = "Instance types for managed node group"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "node_capacity_type" {
  description = "Capacity type for node group (ON_DEMAND or SPOT)"
  type        = string
  default     = "ON_DEMAND"
}

variable "tags" {
  description = "Common tags applied to all resources"
  type        = map(string)
  default     = {}
}

variable "admin_role_name" {
  description = "Name for the IAM role that will administer the EKS cluster"
  type        = string
  default     = "eks-admin"
}

variable "admin_role_trusted_principal_arns" {
  description = "List of IAM principal ARNs (users/roles) allowed to assume the admin role"
  type        = list(string)
  default     = []
}


