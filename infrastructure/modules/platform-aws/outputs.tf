output "vpc_id" {
  value       = module.vpc.vpc_id
  description = "VPC ID"
}

output "private_subnets" {
  value       = module.vpc.private_subnets
  description = "Private subnet IDs"
}

output "public_subnets" {
  value       = module.vpc.public_subnets
  description = "Public subnet IDs"
}

output "cluster_name" {
  value       = module.eks.cluster_name
  description = "EKS cluster name"
}

output "cluster_endpoint" {
  value       = module.eks.cluster_endpoint
  description = "EKS cluster endpoint"
}

output "cluster_ca_certificate" {
  value       = module.eks.cluster_certificate_authority_data
  description = "EKS cluster CA cert (base64)"
}

output "cluster_security_group_id" {
  value       = module.eks.cluster_security_group_id
  description = "EKS cluster security group ID"
}

output "oidc_provider_arn" {
  value       = module.eks.oidc_provider_arn
  description = "OIDC provider ARN for IRSA"
}

output "admin_role_arn" {
  value       = aws_iam_role.eks_admin.arn
  description = "ARN of the EKS admin IAM role"
}


