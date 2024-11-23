readonly -a arr1=(a b c d e f g h)
for i in "${arr1[@]}"
do
  pushd "json-rest/service-$i" || exit
  rm -rf go.mod
  rm -rf go.sum
  go mod init "k8s-istio-service-$i"
  go mod tidy -v
  popd || exit
done

