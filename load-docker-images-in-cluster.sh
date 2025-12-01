# Build locally
docker build -t poc-fastapi:latest ./apps/fastapi
docker build -t poc-golang:latest ./apps/golang


# Load into kind (so K8s can pull them without a registry)
kind load docker-image poc-fastapi:latest --name poc-compare
kind load docker-image poc-golang:latest --name poc-compare
