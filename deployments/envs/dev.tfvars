aws_region      = "us-west-2"
name_prefix     = "vtrips-dev"
cluster_name    = "vtrips-dev-eks"
cluster_version = "1.30"

vpc_cidr = "10.20.0.0/16"
azs = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]
public_subnets = [
  "10.20.0.0/20",
  "10.20.16.0/20",
  "10.20.32.0/20",
]
private_subnets = [
  "10.20.128.0/20",
  "10.20.144.0/20",
  "10.20.160.0/20",
]

enable_nat_gateway = true
single_nat_gateway = true

node_min_size       = 1
node_desired_size   = 2
node_max_size       = 3
node_instance_types = ["t3.medium"]
node_capacity_type  = "ON_DEMAND"

admin_role_name = "eks-admin"
admin_role_trusted_principal_arns = [
  "arn:aws:iam::188530693390:user/taite",
]

tags = {
  env   = "dev"
  app   = "vtrips"
  owner = "platform"
}


