#!/bin/bash

sha=${GITHUB_SHA::8}
ts=$(date +%s)
build_id=${GITHUB_REF_NAME}-${sha}-${ts}

image_build_tag=${ACR_NAME}.azurecr.io/${IMAGE_NAME}:$build_id
image_latest_tag=${ACR_NAME}.azurecr.io/${IMAGE_NAME}:${GITHUB_REF_NAME}-latest

docker build . -t $image_build_tag -t $image_latest_tag

az acr login --name ${ACR_NAME}
docker push $image_build_tag
docker push $image_latest_tag