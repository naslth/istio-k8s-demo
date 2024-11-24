# Reference Application Platform Helm Chart

Prerequisite: Kubernetes Metrics Server for HPA

<https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html>

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl get deployment metrics-server -n kube-system
```

Install Helm Chart

```shell
# perform dry run
helm install k8s-istio ./k8s-istio --namespace dev --debug --dry-run

# apply chart resources
kubectl delete namespace dev
kubectl create namespace dev
kubectl label namespace dev istio-injection=enabled
helm install k8s-istio ./k8s-istio --namespace dev
