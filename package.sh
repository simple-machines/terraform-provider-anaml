#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

# Parse a version number from the Git release tag.
set +e
GIT_TAG=$(git describe --tags --abbrev=0)
if [[ $? -ne 0 ]]; then
  echo "No Git tag set on current commit, skipping publication"
  exit 1
fi
set -e

if ! [[ "$GIT_TAG" =~ ^release-v([0-9]+\.[0-9]+\.[0-9]+)$ ]]; then
  echo "No release Git tag set on current commit, skipping publication. Actual tag: $GIT_TAG" 
  exit 1
fi

VERSION_RELEASE=${BASH_REMATCH[1]}

# Create a directory containing the files to package
PACKAGE_DIR=$(mktemp -d)

cat <<EOT > "$PACKAGE_DIR/version.sh"
VERSION=$VERSION_RELEASE
EOT

cp install.sh terraform-provider-anaml terraform-provider-anaml-operations "$PACKAGE_DIR"
chmod +x "$PACKAGE_DIR/install.sh"

# Create a self-extracting archive as an installer
makeself "$PACKAGE_DIR" "terraform-provider-anaml-install-${VERSION_RELEASE}.run" "Anaml Terraform provider $VERSION_RELEASE" ./install.sh
