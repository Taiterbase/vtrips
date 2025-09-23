#k8s_yaml(helm('deployments/localstack', name='localstack'), allow_duplicates=False)
#k8s_resource('localstack', port_forwards=['4566:4566'], labels=["infra"])

# App Tiltfiles
load_dynamic('./apps/trips/Tiltfile')
load_dynamic('./apps/users/Tiltfile')
load_dynamic('./apps/frontend/Tiltfile')

#local_resource("desktop", "cd ./apps/desktop && wails dev", deps=['./apps/frontend', './apps/desktop'])

# Optional load testing tool depends on trips service
#local_resource("load_test", "go run ./apps/load_test/main.go", resource_deps=['trips'])
