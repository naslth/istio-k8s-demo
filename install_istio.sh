#!/bin/bash

# helm repo add istio.io https://storage.googleapis.com/istio-prerelease/daily-build/master-latest-daily/charts
# helm repo list

# install helm
helm install istio-base istio/base -n istio-system 
helm install istiod istio/istiod -n istio-system --wait
helm install istio-ingressgateway istio/gateway -n istio-system
helm install istio-egressgateway istio/gateway -n istio-system --set service.type=ClusterIP

#check istio
kubectl get svc -n istio-system
kubectl get pods -n istio-system
helm list istio


#apply 3rd party tools
#kiali
helm install \
    --set cr.create=true \
    --set cr.namespace=istio-system \
    --set cr.spec.auth.strategy="anonymous" \
    --namespace kiali-operator \
    --create-namespace \
    kiali-operator \
    kiali/kiali-operator
# check kiali svc
kubectl -n istio-system get svc kiali

#open dashboard kiali
istioctl dashboard kiali