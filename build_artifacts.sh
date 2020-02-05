#!/bin/bash
set -x
set -ste

rm -rf terraform.d/

# build for each OS
mkdir -p terraform.d/plugins/darwin_amd64
GOOS=darwin go build -mod=vendor -o terraform.d/plugins/darwin_amd64/terraform-provider-fruitpi_${DRONE_TAG}

mkdir -p terraform.d/plugins/linux_amd64
GOOS=linux go build -mod=vendor -o terraform.d/plugins/linux_amd64/terraform-provider-fruitpi_${DRONE_TAG}

# bundle them all up
rm -rf dist
mkdir dist/
tar -zcvf dist/terraform-provider-fruitpi_${DRONE_TAG}.tar.gz terraform.d/

md5sum dist/terraform-provider-fruitpi_${DRONE_TAG}.tar.gz
