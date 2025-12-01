kubectl apply -f k8s/namespace.yaml

kubectl apply -f k8s/prometheus/prometheus-configmap.yaml
kubectl apply -f k8s/prometheus/prometheus-deployment.yaml
kubectl apply -f k8s/prometheus/prometheus-service.yaml

kubectl apply -f k8s/prometheus/grafana-deployment.yaml

kubectl -n poc-compare apply -f k8s/fastapi-service/
kubectl -n poc-compare apply -f k8s/go-service/
kubectl -n poc-compare get pods -w
