# kubectl -n poc-compare apply -f k8s/k6-job.yaml 
kubectl -n poc-compare logs job/k6-golang-base -f > ./test-load-results/k6-golang-health-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-fastapi-base -f > ./test-load-results/k6-fastapi-health-results.txt
sleep 3s


kubectl -n poc-compare logs job/k6-golang -f > ../test-load-results/k6-golang-cpu-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-fastapi -f > ../test-load-results/k6-fastapi-cpu-results.txt
sleep 3


kubectl -n poc-compare logs job/k6-golang-io -f > ../test-load-results/k6-golang-io-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-fastapi-io -f > ../test-load-results/k6-fastapi-io-results.txt
sleep 3

kubectl -n poc-compare logs job/k6-golang-json -f > ../test-load-results/k6-golang-json-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-fastapi-json -f > ../test-load-results/k6-fastapi-json-results.txt
sleep 3
