resource "kubernetes_namespace" "ns" {
  count = var.create_namespace ? 1 : 0
  metadata {
    name = var.namespace
  }
}

resource "helm_release" "app" {
  name       = var.release_name
  namespace  = var.namespace

  chart      = var.chart_path

  wait       = var.wait
  timeout    = var.timeout
  atomic     = var.atomic

  set = [for k, v in var.set_values : {
    name  = k
    value = v
  }]

  set_sensitive = [for k, v in var.set_sensitive_values : {
    name  = k
    value = v
  }]

  values = [for f in var.values_files : file(f)]
}


