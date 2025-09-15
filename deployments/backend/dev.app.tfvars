aws_region  = "us-west-2"
cluster_name = "vtrips-dev-eks"

namespace    = "backend"
release_name = "backend"

chart_path = "./"

values_files = [
  "./values.yaml",
]

set_values = {
  "logLevel" = "INFO"
}

set_sensitive_values = {
  //"db.connectionUrl" = "jdbc:postgresql://postgres:password@localhost:5432/vtrips"
}

create_namespace = true
wait             = true
timeout          = 600
atomic           = true

ecr_repo_name = "vtrips-backend"
image_tag     = "dev"
