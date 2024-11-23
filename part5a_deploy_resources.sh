#!/bin/bash
# purpose: Optional: Deploy Kubernetes/Istio resources

# Constants - CHANGE ME!
readonly NAMESPACES=(dev)
readonly SERVICES=(a b c d e f g h)


# Create Namespaces
kubectl apply -f ./resources/other/namespaces.yaml
kubectl apply -f ./resources/other/istio-gateway.yaml
kubectl apply -f ../golang-srv-demo-secrets/other/external-mesh-mongodb-atlas.yaml
kubectl apply -f ../golang-srv-demo-secrets/other/external-mesh-cloudamqp.yaml

for namespace in "${NAMESPACES[@]}"; do
  # Enable automatic Istio sidecar injection
  kubectl label namespace "${namespace}" istio-injection=enabled

  kubectl apply -n "${namespace}" -f ../golang-srv-demo-secrets/secret/k8s-istio-service-demo.yaml

  for service in "${SERVICES[@]}"; do
    kubectl apply -n "${namespace}" -f "./resources/services/service-$service.yaml"
  done
  kubectl apply -n "${namespace}" -f ./resources/services/angular-ui.yaml
done
