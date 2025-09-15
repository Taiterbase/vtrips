module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = var.name_prefix
  cidr = var.vpc_cidr

  azs             = var.azs
  public_subnets  = var.public_subnets
  private_subnets = var.private_subnets

  enable_nat_gateway = var.enable_nat_gateway
  single_nat_gateway = var.single_nat_gateway

  tags = var.tags
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name                   = var.cluster_name
  cluster_version                = var.cluster_version
  cluster_endpoint_public_access = true

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  enable_irsa = true
  authentication_mode = "API"

  eks_managed_node_groups = {
    default = {
      min_size       = var.node_min_size
      max_size       = var.node_max_size
      desired_size   = var.node_desired_size
      instance_types = var.node_instance_types
      capacity_type  = var.node_capacity_type
    }
  }

  tags = var.tags
}

# Admin IAM role to administer the EKS cluster. Trusted principals may assume it.
resource "aws_iam_role" "eks_admin" {
  name = var.admin_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      for principal in var.admin_role_trusted_principal_arns : {
        Effect = "Allow"
        Principal = {
          AWS = principal
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = var.tags
}

# Attach broad AWS permissions to the admin role (optional but practical for admin tasks)
resource "aws_iam_role_policy_attachment" "eks_admin_adminaccess" {
  role       = aws_iam_role.eks_admin.name
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}

# Grant the admin role Kubernetes admin access to the cluster
resource "aws_eks_access_entry" "eks_admin" {
  cluster_name  = var.cluster_name
  principal_arn = aws_iam_role.eks_admin.arn
  type          = "STANDARD"
}

resource "aws_eks_access_policy_association" "admin_role_cluster_admin" {
  cluster_name  = var.cluster_name
  principal_arn = aws_iam_role.eks_admin.arn
  policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
  access_scope {
    type = "cluster"
  }
}

# Also grant direct access for each trusted principal so they can interact without assuming the role
resource "aws_eks_access_entry" "trusted_principals" {
  for_each      = toset(var.admin_role_trusted_principal_arns)
  cluster_name  = var.cluster_name
  principal_arn = each.value
  type          = "STANDARD"
}

resource "aws_eks_access_policy_association" "trusted_principals_cluster_admin" {
  for_each      = toset(var.admin_role_trusted_principal_arns)
  cluster_name  = var.cluster_name
  principal_arn = each.value
  policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
  access_scope {
    type = "cluster"
  }
}
