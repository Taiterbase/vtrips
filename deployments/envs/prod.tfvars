aws_region      = "us-west-2"
name_prefix     = "vtrips-prod"
cluster_name    = "vtrips-prod-eks"
cluster_version = "1.30"

vpc_cidr = "10.40.0.0/16"
azs = [
  "us-west-2a",
  "us-west-2b",
  "us-west-2c",
]
public_subnets = [
  "10.40.0.0/20",
  "10.40.16.0/20",
  "10.40.32.0/20",
]
private_subnets = [
  "10.40.128.0/20",
  "10.40.144.0/20",
  "10.40.160.0/20",
]

enable_nat_gateway = true
single_nat_gateway = false

node_min_size       = 3
node_desired_size   = 4
node_max_size       = 8
node_instance_types = ["m6i.large"]
node_capacity_type  = "ON_DEMAND"

tags = {
  env   = "prod"
  app   = "vtrips"
  owner = "platform"
}


