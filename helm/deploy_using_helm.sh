# Cài đặt component metrics để sử dụng thu thập metrics k8s 
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl get deployment metrics-server -n kube-system
# kiểm tra config với flag --dry-run
helm install k8s-istio ./k8s-istio --namespace dev --debug --dry-run
# xoá namespace dev (nếu có)
kubectl delete namespace dev
# tạo namespace dev 
kubectl create namespace dev
# đặt label istio-injection=enabled để inject sidecar khi khởi tạo pod
kubectl label namespace dev istio-injection=enabled
# deploy app với helm chart
helm install k8s-istio ./k8s-istio --namespace dev
