#!/bin/bash

readonly -a arr=( d f )
readonly tag=1.0.5

for i in "${arr[@]}"
do
  docker push "docker.io/naslth/k8s-istio-service-$i:$tag"
done

# docker push "docker.io/naslth/k8s-istio-angular-web-service:$tag"
