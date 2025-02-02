#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ./vendor/k8s.io/code-generator)}

# At some environments (eg. GitHub Actions), due to $GOPATH setting, the codegen output might not be at the expected path
# in the project repo, therefore we should force the output to the specific directory (see --output-base)
# we need to handle (move) the generated files to the correct location in the repo then
CODEGEN_OUTPUT_BASE="${SCRIPT_ROOT}"/output
CODEGEN_OUTPUT_GENERATED="${CODEGEN_OUTPUT_BASE}"/github.com/rabbitmq/messaging-topology-operator/pkg/generated

# Move the API files to here so that codegen can find them properly (wants the group as part of the name)
mkdir -p api/rabbitmq.com/v1alpha2
cp api/v1alpha2/* api/rabbitmq.com/v1alpha2/

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
# bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
# Deepcopy is generated by operator-sdk
bash "${CODEGEN_PKG}"/generate-groups.sh "client,informer,lister" \
  github.com/rabbitmq/messaging-topology-operator/pkg/generated github.com/rabbitmq/messaging-topology-operator/api \
  rabbitmq.com:v1alpha2 \
  --go-header-file "${SCRIPT_ROOT}"/hack/NOTICE.go.txt --output-base "${CODEGEN_OUTPUT_BASE}"

if [ -d "${CODEGEN_OUTPUT_GENERATED}" ]; then
  mkdir -p "${SCRIPT_ROOT}"/pkg/generated
  rm -rf "${SCRIPT_ROOT}"/pkg/generated
  mv "${CODEGEN_OUTPUT_GENERATED}" "${SCRIPT_ROOT}"/pkg/generated
  rm -rf "${SCRIPT_ROOT}"/output
fi

# Clean up the tmp move
rm -rf api/rabbitmq.com

# Kubebuilder project layout has api under 'api/v1alpha2'
# client-go codegen expects group name in the path ie. 'api/rabbitmq.com/v1alpha2'
# Because there's no way how to modify any of these settings,
# we need to hack things a little bit (replace the name of package)
find pkg/generated -type f -name "*.go" |\
xargs sed -i".out" -e "s#github.com/rabbitmq/messaging-topology-operator/api/rabbitmq.com/v1alpha2#github.com/rabbitmq/messaging-topology-operator/api/v1alpha2#g"
find pkg/generated -type f -name "*.go.out" | xargs rm -rf

