aws_region  = "us-west-2"
cluster_name = "vtrips-dev-eks"

namespace    = "frontend"
release_name = "frontend"

chart_path = "./"

values_files = [
  "./values.yaml",
]

set_values = {}

set_sensitive_values = {}

create_namespace = true
wait             = true
timeout          = 600
atomic           = true

ecr_repo_name = "vtrips-frontend"
image_tag     = "dev"

# Set to true to build/push locally; otherwise provide image_repository
build_image       = false
image_repository  = ""

