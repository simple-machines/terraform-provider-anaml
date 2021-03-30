#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR=$(dirname $0)

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

mkdir -p "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml/$VERSION/darwin_amd64/"
mkdir -p "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml-operations/$VERSION/darwin_amd64/"

cp terraform-provider-anaml "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml/$VERSION/darwin_amd64/terraform-provider-anaml_v$VERSION"
cp terraform-provider-anaml-operations "$HOME/.terraform.d/plugins/registry.anaml.io/anaml/anaml-operations/$VERSION/darwin_amd64/terraform-provider-anaml-operations_v$VERSION"

echo "Successfully installed Anaml Terraform providers version $VERSION."
echo "Remember to run 'terraform init -upgrade' to upgrade to the latest provider."
