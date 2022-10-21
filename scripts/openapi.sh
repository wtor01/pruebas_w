#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

mkdir -p "$output_dir"

oapi-codegen -generate types -templates api/templates -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yaml"
oapi-codegen -generate gin -templates api/templates -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yaml"
