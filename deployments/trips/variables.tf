variable "aws_region" { type = string }
variable "cluster_name" { type = string }
variable "namespace" { type = string }
variable "release_name" { type = string }
variable "chart_path" { type = string }
variable "values_files" { type = list(string) }
variable "set_values" { type = map(string) }
variable "set_sensitive_values" { type = map(string) }
variable "create_namespace" { type = bool }
variable "wait" { type = bool }
variable "timeout" { type = number }
variable "atomic" { type = bool }

variable "image_tag" { type = string }
variable "ecr_repo_name" { type = string }
