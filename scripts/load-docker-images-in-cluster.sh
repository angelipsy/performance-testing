# Build locally
docker build -t poc-fastapi:latest ./apps/fastapi
docker build -t poc-golang:latest ./apps/golang


# Load into kind
kind load docker-image poc-fastapi:latest --name poc-compare
kind load docker-image poc-golang:latest --name poc-compare
