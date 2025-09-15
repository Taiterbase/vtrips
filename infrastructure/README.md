## Infrastructure (Terraform)

This directory contains reusable Terraform modules and example stacks to provision the shared platform (AWS VPC + EKS) and to deploy each app via Helm into the cluster.

### Structure

- `modules/platform-aws`: Provisions VPC and EKS (using official community modules)
- `modules/app-helm`: Deploys a Helm chart into a namespace with values from Terraform
- Roots and env var-files now live under `deployments/`:
  - `deployments/platform`: Root to create the shared platform
  - `deployments/<service>`: Roots to deploy each app via Helm (next to each chart)
  - `deployments/envs`: Shared environment var-files (dev/staging/prod)

### Usage Overview

1. Create the platform once per environment:

```bash
cd deployments/platform
terraform init
terraform plan -var-file="../envs/dev.tfvars"
terraform apply -var-file="../envs/dev.tfvars"
```

2. Deploy an app (example: backend) to an existing EKS cluster:

```bash
cd deployments/backend
terraform init
terraform plan -var-file="../envs/dev.tfvars" -var-file="./dev.app.tfvars"
terraform apply -var-file="../envs/dev.tfvars" -var-file="./dev.app.tfvars"
```

Notes:

- Platform stack expects AWS credentials in your environment and will create VPC + EKS.
- App stacks read the cluster by name via AWS data sources and configure `kubernetes`/`helm` providers automatically.
- Helm charts are sourced from the repository paths under `deployments/<service>` by default (can be swapped for a remote repo).
