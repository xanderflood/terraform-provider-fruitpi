#!/bin/bash
set -x
set -ste

apt-get update
apt-get install -y curl jq

export RELEASE_URL=`curl -s -H "Authorization: token $GITHUB_TOKEN" \
	https://api.github.com/repos/xanderflood/terraform-provider-fruitpi/releases \
	| jq -r --arg tag "$DRONE_TAG" '.[] | select(.["tag_name"] == $tag) | .["html_url"]'`

cat << EOF
Successfully created a draft release of {{repo.name}} for tagged version {{build.tag}}.

This release can be viewed and finalized here: $RELEASE_URL
EOF
