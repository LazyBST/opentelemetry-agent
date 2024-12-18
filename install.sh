#!/bin/bash

echo -e "Initiating KMagent setup for mac arm64.\n"

mkdir /tmp/otlp-agent
mkdir /tmp/otlp-agent/collector

cp ./builder-config.yaml /tmp/otlp-agent/builder-config.yaml

curl --proto '=https' --tlsv1.2 -fL -o /tmp/otlp-agent/ocb \
https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv0.116.0/ocb_0.116.0_darwin_arm64

chmod -R 777 /tmp/otlp-agent/

cd ./agent && GOOS=darwin GOARCH=arm64 go build -o agent .

./agent