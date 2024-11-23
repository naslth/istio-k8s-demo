#!/bin/bash
# purpose: Tear down GKE cluster and associated resources

# Constants - CHANGE ME!
readonly PROJECT='k8s-istio-service-demo'
readonly CLUSTER='k8s-istio-service-demo-cluster'
readonly REGION='us-central1'

# Delete GKE cluster (time in foreground)
yes | gcloud beta container clusters delete ${CLUSTER} --region ${REGION}

# Confirm network resources are also deleted
gcloud compute forwarding-rules list
gcloud compute target-pools list
gcloud compute firewall-rules list

# In case target-pool associated with Cluster is not deleted
yes | gcloud compute target-pools delete  \
  $(gcloud compute target-pools list \
    --filter="region:($REGION)" --project ${PROJECT} \
  | awk 'NR==2 {print $1}')
