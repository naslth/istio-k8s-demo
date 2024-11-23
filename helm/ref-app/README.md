# Reference Application Platform Helm Chart

This Helm 3 chart will install all Kubernetes resources to the `dev` namespace for the Reference Application Platform. Before proceeding, add your environment specific values in the chart's `values.yaml`. Note that this chart includes container resource requests and limits, along with the use `HorizontalPodAutoscaler` resources.

Prerequisite: Kubernetes Metrics Server for HPA

<https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html>

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl get deployment metrics-server -n kube-system
```

Install Helm Chart

```shell
# perform dry run
helm install ref-app ./ref-app --namespace dev --debug --dry-run
kubectl label namespace dev istio-injection=enabled
# apply chart resources
helm install ref-app ./ref-app --namespace dev --create-namespace
