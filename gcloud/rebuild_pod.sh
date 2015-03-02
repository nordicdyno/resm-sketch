#!/bin/bash
set -x
gcloud preview container kubectl delete pod resmpd

gcloud preview container kubectl create -f resmpd.yml
