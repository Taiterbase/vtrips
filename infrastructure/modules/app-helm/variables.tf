variable "cluster_name" {
  description = "EKS cluster name to target"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace to install into"
  type        = string
}

variable "release_name" {
  description = "Helm release name"
  type        = string
}

variable "chart_path" {
  description = "Local path to Helm chart (or remote chart reference)"
  type        = string
}

variable "values_files" {
  description = "List of values files to pass to Helm"
  type        = list(string)
  default     = []
}

variable "set_values" {
  description = "Map of simple set values for the Helm release"
  type        = map(string)
  default     = {}
}

variable "set_sensitive_values" {
  description = "Map of sensitive set values for the Helm release"
  type        = map(string)
  default     = {}
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "wait" {
  description = "Whether to wait for resources to become ready"
  type        = bool
  default     = true
}

variable "timeout" {
  description = "Wait timeout in seconds (e.g., 600 for 10m)"
  type        = number
  default     = 600
}

variable "atomic" {
  description = "Rollback changes if install/upgrade fails"
  type        = bool
  default     = true
}

# Optional: used to create an implicit dependency from callers (e.g., image build)
variable "build_dependency" {
  description = "Opaque dependency handle to enforce ordering from caller"
  type        = string
  default     = ""
}


