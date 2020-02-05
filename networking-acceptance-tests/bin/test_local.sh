#!/bin/bash

set -euo pipefail


set +u
if [[ -z $1 ]]; then
  echo "Usage: ./test.sh <test_config_path> [kube_config_path]"
  exit 1
fi
set -u

kube_config_path=${2:="${HOME}/.kube/config"}
test_config_path="$1"

CONFIG="${test_config_path}" KUBECONFIG="${kube_config_path}" ginkgo -v .
