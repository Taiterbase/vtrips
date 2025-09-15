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

resource "aws_ecr_repository" "this" {
  name                 = var.ecr_repo_name
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "null_resource" "build_and_push" {
  triggers = {
    image_tag   = var.image_tag
    repo_url    = aws_ecr_repository.this.repository_url
    dockerfile  = filemd5("./Dockerfile")
    source_hash = filesha256("../../apps/backend/cmd/main.go")
  }

  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-lc"]
    command     = <<-EOT
      set -euo pipefail
      AWS_REGION=${var.aws_region}
      REPO_URL=${aws_ecr_repository.this.repository_url}
      IMAGE_TAG=${var.image_tag}

      aws ecr get-login-password --region "$AWS_REGION" | docker login --username AWS --password-stdin "$REPO_URL"

      # Use repo root as build context so Dockerfile paths like apps/backend/... exist
      docker build -t "$REPO_URL:$IMAGE_TAG" -f ./Dockerfile ../../
      docker push "$REPO_URL:$IMAGE_TAG"
    EOT
  }
}

locals {
  ecr_image_repository = aws_ecr_repository.this.repository_url
  ecr_image_tag        = var.image_tag
}

module "app" {
  source = "../../infrastructure/modules/app-helm"

  cluster_name     = var.cluster_name
  aws_region       = var.aws_region
  namespace        = var.namespace
  release_name     = var.release_name
  chart_path       = var.chart_path
  values_files     = var.values_files
  set_values       = merge(var.set_values, {
    "image.repository" = local.ecr_image_repository
    "image.tag"        = local.ecr_image_tag
  })
  set_sensitive_values = var.set_sensitive_values
  create_namespace = var.create_namespace
  wait             = var.wait
  timeout          = var.timeout
  atomic           = var.atomic

  # Implicit dependency using an input consumed in the module
  build_dependency = null_resource.build_and_push.id
}


