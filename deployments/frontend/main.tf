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
  count                = var.build_image ? 1 : 0
  name                 = var.ecr_repo_name
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "null_resource" "build_and_push" {
  count = var.build_image ? 1 : 0
  triggers = {
    image_tag   = var.image_tag
    repo_url    = var.build_image ? aws_ecr_repository.this[0].repository_url : ""
    dockerfile  = filemd5("./Dockerfile")
    source_hash = filesha256("../../apps/frontend/cmd/main.go")
  }

  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-lc"]
    command     = <<-EOT
      set -euo pipefail
      AWS_REGION=${var.aws_region}
      REPO_URL=${aws_ecr_repository.this[0].repository_url}
      IMAGE_TAG=${var.image_tag}

      aws ecr get-login-password --region "$AWS_REGION" | docker login --username AWS --password-stdin "$REPO_URL"

      # Use repo root as build context so Dockerfile paths like apps/frontend/... exist
      docker build -t "$REPO_URL:$IMAGE_TAG" -f ./Dockerfile ../../
      docker push "$REPO_URL:$IMAGE_TAG"
    EOT
  }
}

locals {}

module "app" {
  source = "../../infrastructure/modules/app-helm"

  cluster_name     = var.cluster_name
  aws_region       = var.aws_region
  namespace        = var.namespace
  release_name     = var.release_name
  chart_path       = var.chart_path
  values_files     = var.values_files
  set_values       = merge(var.set_values,
    var.build_image ? {
      "image.repository" = aws_ecr_repository.this[0].repository_url
      "image.tag"        = var.image_tag
    } : (
      var.image_repository != "" ? {
        "image.repository" = var.image_repository
        "image.tag"        = var.image_tag
      } : {}
    )
  )
  set_sensitive_values = var.set_sensitive_values
  create_namespace = var.create_namespace
  wait             = var.wait
  timeout          = var.timeout
  atomic           = var.atomic

  build_dependency = var.build_image ? null_resource.build_and_push[0].id : ""
}


