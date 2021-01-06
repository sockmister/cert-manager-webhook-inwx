#!/bin/sh
os=$(uname -s)
curl -sSL  https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.15.5-"$(echo "$os" | tr '[:upper:]' '[:lower:]')"-amd64.tar.gz | tar -zvxf -

