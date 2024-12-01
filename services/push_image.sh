#!/bin/bash

readonly -a arr=(d )
readonly tag=1.1.2

for i in "${arr[@]}"
do
  docker push "naslth/k8s-istio-service-$i:$tag"
done

