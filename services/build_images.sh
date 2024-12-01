#!/bin/bash
# purpose: Build Go microservices for demo

readonly -a arr=(d)
readonly tag=1.1.2

for i in "${arr[@]}"
do
  cp -f Dockerfile "service-$i"
  pushd "service-$i" || exit
  docker build -t "naslth/k8s-istio-service-$i:$tag" . --no-cache
  rm -rf Dockerfile
  popd || exit
done

docker image ls | grep 'naslth/k8s-istio-service-'