#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR=$(dirname $0)

ARCH="darwin_amd64"

if [[ $(arch) == "arm64" ]]; then
  ARCH="darwin_arm64"
fi

echo $ARCH

if ! [[ -f version.sh ]]; then
  echo "Version configuration not found. Aborting installation."
  exit 1
fi

if ! [[ -f terraform-provider-anaml ]]; then
  echo "Anaml provider executable 'terraform-provider-anaml' not found. Aborting installation."
  exit 1
fi

if ! [[ -f terraform-provider-anaml-operations ]]; then
  echo "Anaml provider executable 'terraform-provider-anaml-operations' not found. Aborting installation."
  exit 1
fi

source "$SCRIPT_DIR/version.sh"

echo "Installing Anaml Terraform providers version $VERSION under '$HOME/.terraform.d/plugins/registry.anaml.io'"

mkdir -p "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml/$VERSION/$ARCH/"
mkdir -p "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml-operations/$VERSION/$ARCH/"

cp terraform-provider-anaml "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml/$VERSION/$ARCH/terraform-provider-anaml_v$VERSION"
cp terraform-provider-anaml-operations "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml-operations/$VERSION/$ARCH/terraform-provider-anaml-operations_v$VERSION"

echo "Successfully installed Anaml Terraform providers version $VERSION."
echo "Remember to run 'terraform init -upgrade' to upgrade to the latest provider."
