#!/bin/bash

kubectl delete namespace dev test

helm delete --purge istio
helm delete --purge istio-init

kubectl get all -n dev
kubectl get all -n test
istioctl get all
