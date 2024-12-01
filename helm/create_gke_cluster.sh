#!/bin/bash

# Biến sử dụng (thay đổi theo config)
readonly PROJECT="final-thesis-20021354"
readonly CLUSTER="k8s-istio-cluster"
readonly ZONE="us-central1"
readonly GKE_VERSION="1.28.14-gke.1376000"
readonly MACHINE_TYPE="e2-standard-2"
readonly SERVICE_ACCOUNT="346848910511-compute@developer.gserviceaccount.com"

# Tạo gcloud clusters
gcloud beta container \
  --project "${PROJECT}" clusters create "${CLUSTER}" \
  --zone "${ZONE}" \
  --no-enable-basic-auth \
  --cluster-version "${GKE_VERSION}" \
  --no-enable-insecure-kubelet-readonly-port \
  --release-channel "regular" \
  --machine-type "${MACHINE_TYPE}" \
  --image-type "COS_CONTAINERD" \
  --disk-type "pd-standard" \
  --disk-size "100" \
  --metadata disable-legacy-endpoints=true \
  --service-account "${SERVICE_ACCOUNT}" \
  --num-nodes "1" \
  --logging=SYSTEM,WORKLOAD \
  --monitoring=SYSTEM \
  --enable-ip-alias \
  --network "projects/${PROJECT}/global/networks/default" \
  --subnetwork "projects/${PROJECT}/regions/${ZONE}/subnetworks/default" \
  --no-enable-intra-node-visibility \
  --default-max-pods-per-node "100" \
  --enable-autoscaling \
  --min-nodes "0" --max-nodes "1" \
  --enable-dataplane-v2 \
  --no-enable-master-authorized-networks \
  --addons HorizontalPodAutoscaling,HttpLoadBalancing,Istio,GcePersistentDiskCsiDriver \
  --enable-autoupgrade \
  --enable-autorepair \
  --max-surge-upgrade 1 \
  --max-unavailable-upgrade 0 \
  --enable-autoprovisioning \
  --min-cpu 1 --max-cpu 2 \
  --min-memory 4 --max-memory 4 \
  --autoprovisioning-locations=us-central1-a,us-central1-b,us-central1-c \
  --autoprovisioning-service-account=${SERVICE_ACCOUNT} \
  --enable-autoprovisioning-autorepair \
  --enable-autoprovisioning-autoupgrade \
  --autoprovisioning-max-surge-upgrade 1 \
  --autoprovisioning-max-unavailable-upgrade 0 \
  --enable-vertical-pod-autoscaling \
  --enable-shielded-nodes \
  --node-locations "us-central1-a","us-central1-b","us-central1-c"

# lấy credital cho clusters 
gcloud container clusters get-credentials "${CLUSTER}" \
  --region "${ZONE}" --project "${PROJECT}"

# sử dụng context gcloud cluster cho command kubectl
kubectl config current-context