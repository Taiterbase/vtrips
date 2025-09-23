aws_region  = "us-west-2"
cluster_name = "vtrips-dev-eks"

namespace    = "users"
release_name = "users"

chart_path = "./"

values_files = [
  "./values.yaml",
]

set_values = {
  "logLevel" = "INFO"
}

set_sensitive_values = {}

create_namespace = true
wait             = true
timeout          = 600
atomic           = true

ecr_repo_name = "vtrips-users"
image_tag     = "dev"
