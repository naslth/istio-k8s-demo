#!/bin/bash


# set helm repo
helm repo add istio.io https://storage.googleapis.com/istio-prerelease/daily-build/master-latest-daily/charts
helm repo list

# install with helm
kubectl create namespace istio-system
helm install istio-base istio/base -n istio-system  
helm install istiod istio/istiod -n istio-system  --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_CAPTURE="true" --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_AUTO_ALLOCATE="true"
helm install istio-ingressgateway istio/gateway -n istio-system  
helm install istio-egressgateway istio/gateway -n istio-system  --set service.type=ClusterIP 
# kubectl create -n istio-system secret tls naslth-credential   --key=myddns_cert/naslth.myddns.me.key   --cert=myddns_cert/naslth.myddns.me.crt 

#check istio
kubectl get svc -n istio-system
kubectl get pods -n istio-system
helm list