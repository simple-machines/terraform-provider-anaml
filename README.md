# terraform-provider-anaml

A Terraform provider for [anaml](https://anaml.io/) resources, such as entities, tables, and features.

## Installation

Download the latest [release on GitHub](https://github.com/simple-machines/terraform-provider-anaml/releases).

Run the installer script:

```
ANAML_TERRAFORM_PROVIDER_VERSION=1.2.3
./terraform-provider-anaml-install-$ANAML_TERRAFORM_PROVIDER_VERSION.run
```

## Development

### Prerequisites

* [Go](https://golang.org/dl/) 1.14 or later

Alternatively, use [Nix](https://nixos.org/download.html) to get a quick,
reproducible development environment.

### Building

To start a Nix shell with Go installed:

```
nix-shell -p go
```

To build the Go modules:

```
go build ./...
```

### Testing with Terraform

To start a Nix shell that builds and installs the Terraform provider:

```
nix-shell tf-shell.nix
```

To run example Terraform manifests:

```
cd examples
terraform init
terraform apply
```

Set the environment variable `TF_LOG=debug` when running Terraform to set debug
level logging.

### Releasing

Terraform assumes that providers follows [Semantic Versioning](https://semver.org/).

The versioning semantics are used to select the latest version of a Terraform
provider when no exact version number is specified in the provider
configuration, such as if no version number is specified, or if version
constraints are specified instead.

#### Private Release

A private release publishes a self-extracting archive as a GitHub release, that
can be downloaded and run to install the Terraform provider locally. The
provider isn't published to the public Terraform registry, so each version must
be downloaded and installed manually.

To release a new version of the Terraform provider, create a new GitHub release
with a tag that matches the regular expression: `^release-v[0-9]+\.[0-9]+\.[0-9]+$`.

For example, to release version 1.0.0, create a GitHub release with the tag `release-v1.0.0`. 

#### Public Release

A public release publishes to the public Terraform registry, so each version
can be installed automatically through Terraform configuration and `terraform
init`.

CI/CD configuration for publishing to the Terraform registry is based on the
documentation from the Terraform website on
[Publishing Providers](https://www.terraform.io/docs/registry/providers/publishing.html).
