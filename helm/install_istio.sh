#!/bin/bash
# thêm repo istio vào helm
helm repo add istio.io https://storage.googleapis.com/istio-prerelease/daily-build/master-latest-daily/charts
# tạo namespace istio-system (root namespace)
kubectl create namespace istio-system
# cài đặt istio control plane
helm install istio-base istio/base -n istio-system
helm install istiod istio/istiod -n istio-system
helm install istio-cni istio/cni -n kube-system --set components.cni.enabled=true
# cài đặt gateway istio
helm install istio-ingressgateway istio/gateway -n istio-system
helm install istio-egressgateway istio/gateway -n istio-system  --set service.type=ClusterIP
#tạo ssl certificate trong secrets
kubectl create -n istio-system secret tls naslth-credential   --key=myddns_cert/naslth.myddns.me.key   --cert=myddns_cert/naslth.myddns.me.crt 
k6 run load-test.js --insecure-skip-tls-verify 

