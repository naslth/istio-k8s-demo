#!/bin/bash

readonly -a arr=(a b c d e f g h)
readonly tag=1.0.0

for i in "${arr[@]}"
do
  docker push "docker.io/naslth/k8s-istio-service-$i:$tag"
done

docker push "docker.io/naslth/k8s-istio-angular-web-service:$tag"
