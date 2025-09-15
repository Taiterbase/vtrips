aws_region      = "us-west-2"
name_prefix     = "vtrips-staging"
cluster_name    = "vtrips-staging-eks"
cluster_version = "1.30"

vpc_cidr = "10.30.0.0/16"
azs = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]
public_subnets = [
  "10.30.0.0/20",
  "10.30.16.0/20",
  "10.30.32.0/20",
]
private_subnets = [
  "10.30.128.0/20",
  "10.30.144.0/20",
  "10.30.160.0/20",
]

enable_nat_gateway = true
single_nat_gateway = true

node_min_size       = 2
node_desired_size   = 3
node_max_size       = 5
node_instance_types = ["t3.large"]
node_capacity_type  = "ON_DEMAND"

tags = {
  env   = "staging"
  app   = "vtrips"
  owner = "platform"
}


