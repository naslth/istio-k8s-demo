#!/bin/bash

# helm repo add istio.io https://storage.googleapis.com/istio-prerelease/daily-build/master-latest-daily/charts
# helm repo list

# install helm
helm install istio-base istio/base -n istio-system meshConfig.outboundTrafficPolicy.mode=ALLOW_ANY --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_CAPTURE="true" --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_AUTO_ALLOCATE="true"
helm install istiod istio/istiod -n istio-system --wait meshConfig.outboundTrafficPolicy.mode=ALLOW_ANY --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_CAPTURE="true" --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_AUTO_ALLOCATE="true"
helm install istio-ingressgateway istio/gateway -n istio-system meshConfig.outboundTrafficPolicy.mode=ALLOW_ANY --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_CAPTURE="true" --set meshConfig.defaultConfig.proxyMetadata.ISTIO_META_DNS_AUTO_ALLOCATE="true"

#check istio
kubectl get svc -n istio-system
kubectl get pods -n istio-system
helm list istio

kubectl create -n istio-system secret tls naslth-credential   --key=myddns_cert/naslth.myddns.me.key   --cert=myddns_cert/naslth.myddns.me.crt 

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

#prometheus
