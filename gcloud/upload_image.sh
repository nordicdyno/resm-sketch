#!/bin/bash
#gcloud config set compute/zone <your-cluster-zone>
#gcloud config set container/cluster <your-cluster-name>
GCP_PROJECT=resm-sketch
GCR_NAMESPACE=${GCP_PROJECT//-/_}
GCR_NAMESPACE=resm-sketch
IMAGE_NAME=gcloud-resm
docker tag -f $IMAGE_NAME gcr.io/${GCR_NAMESPACE}/$IMAGE_NAME
gcloud preview docker push gcr.io/${GCR_NAMESPACE}/$IMAGE_NAME
