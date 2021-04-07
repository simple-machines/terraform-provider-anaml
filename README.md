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

To release a new version of the Terraform provider, create a new GitHub release
with a tag that matches the regular expression: `^release-v[0-9]+\.[0-9]+\.[0-9]+$`.

For example, to release version 1.0.0, create a GitHub release with the tag `release-v1.0.0`. 

Terraform assumes that providers follows [Semantic Versioning](https://semver.org/).

The versioning semantics are used, to select the latest version of a Terraform
provider when no exact version number is specified in the provider
configuration, such as if no version number is specified, or if version
constraints are specified instead.
