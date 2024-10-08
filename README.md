# vtrips

## Usage
`ctlptl delete cluster minikube`
`ctlptl create cluster minikube --registry=ctlptl-registry --minikube-start-flags="--cpus=2" --minikube-start-flags="--memory=4gb"`
`tilt up`

## Testing
```sh
curl -X POST "http://localhost:8080/v1/trips?client_id=test" -H "Content-Type: application/json" -d '{
  "status": "draft",
  "volunteer_limit": 10,
  "name": "Cleaning up Nepal",
  "description": "our mountains are piling up with trash from tourists! we need your help to clean them.",
  "housing_type": "camping"
}'
```

```sh
curl -X PUT "http://localhost:8080/v1/trips/:trip_id?client_id=test" -H "Content-Type: application/json" -d '{
  "status": "listed",
}'
```

```sh
curl -X GET "http://localhost:8080/v1/trips/:trip_id?client_id=test" -H "Content-Type: application/json"
```

```sh
curl -X GET "http://localhost:8080/v1/trips?client_id=test" -H "Content-Type: application/json"
curl -X GET "http://localhost:8080/v1/trips?client_id=test&status=listed" -H "Content-Type: application/json"
curl -X GET "http://localhost:8080/v1/trips?client_id=test&status=listed&housing_type=camping" -H "Content-Type: application/json"
```

