terraform {
  required_version = ">= 1.5.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "platform" {
  source = "../../infrastructure/modules/platform-aws"

  name_prefix         = var.name_prefix
  aws_region          = var.aws_region
  vpc_cidr            = var.vpc_cidr
  azs                 = var.azs
  public_subnets      = var.public_subnets
  private_subnets     = var.private_subnets
  enable_nat_gateway  = var.enable_nat_gateway
  single_nat_gateway  = var.single_nat_gateway
  cluster_name        = var.cluster_name
  cluster_version     = var.cluster_version
  node_min_size       = var.node_min_size
  node_max_size       = var.node_max_size
  node_desired_size   = var.node_desired_size
  node_instance_types = var.node_instance_types
  node_capacity_type  = var.node_capacity_type
  tags                = var.tags
  admin_role_name                          = var.admin_role_name
  admin_role_trusted_principal_arns        = var.admin_role_trusted_principal_arns
}

output "cluster_name" {
  value = module.platform.cluster_name
}


