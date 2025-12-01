# kubectl -n poc-compare apply -f k8s/k6-job-stress.yaml 
kubectl -n poc-compare logs job/k6-stress-golang-base -f > ../test-stress-results/k6-stress-golang-health-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-stress-fastapi-base -f > ../test-stress-results/k6-stress-fastapi-health-results.txt
sleep 3s


kubectl -n poc-compare logs job/k6-stress-golang -f > ../test-stress-results/k6-stress-golang-cpu-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-stress-fastapi -f > ../test-stress-results/k6-stress-fastapi-cpu-results.txt
sleep 3


kubectl -n poc-compare logs job/k6-stress-golang-io -f > ../test-stress-results/k6-stress-golang-io-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-stress-fastapi-io -f > ../test-stress-results/k6-stress-fastapi-io-results.txt
sleep 3

kubectl -n poc-compare logs job/k6-stress-golang-json -f > ../test-stress-results/k6-stress-golang-json-results.txt
sleep 3s
kubectl -n poc-compare logs job/k6-stress-fastapi-json -f > ../test-stress-results/k6-stress-fastapi-json-results.txt
sleep 3
