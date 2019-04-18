#!/bin/bash
# Read tobuild file line by line and docker build and push the image.
# Usage: PROJECT_ID=krn-dev SHORT_SHA=$(git rev-parse --short HEAD~) ./build-containers.sh tobuild

gcr_base="gcr.io/$PROJECT_ID"
short_sha=$SHORT_SHA

while IFS='' read -r line || [[ -n "$line" ]]; do
    image_name=$line
    image_tag=$(date +"%Y%m%d.%H%M%S").$short_sha
    image_path=$gcr_base/$image_name:$image_tag
    echo "Building and pushing image $image_path"
    docker build -t $image_path $image_name
    docker push $image_path
done < "$1"
