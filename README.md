# terraform-provider-anaml

A Terraform provider for [anaml](https://anaml.io/) resources, such as entities, tables, and features.

## Installation

Configure the `anaml` and/or the `anaml-operations` Terraform providers through
your Terraform manifest:

```terraform
terraform {
  required_version = "~> 0.14"
  required_providers {
    anaml = {
      source  = "simple-machines/anaml"
    }
    anaml-operations = {
      source  = "simple-machines/anaml-operations"
    }
  }
}

provider "anaml" {
  host       = "http://localhost:8080"
  username   = "admin"
  password   = "admin-password"
  branch     = "official"
}

provider "anaml-operations" {
  host       = "http://localhost:8080"
  username   = "admin"
  password   = "admin-password"
}
```

When there is a new release, run `terraform init -upgrade` to upgrade to the
latest version.

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

To release a new version of the Terraform provider, create a new _Git tag_ with
a tag that matches the regular expression: `^v[0-9]+\.[0-9]+\.[0-9]+$`. Do not
create a GitHub release, the build process will automatically create a GitHub
release when the build is successful.

For example, to release version 1.0.0, create a GitHub release with the tag `v1.0.0`. 

Terraform assumes that providers follows [Semantic Versioning](https://semver.org/).

The versioning semantics are used to select the latest version of a Terraform
provider when no exact version number is specified in the provider
configuration, such as if no version number is specified, or if version
constraints are specified instead.

The Terraform providers are publshed to the public Terraform Registry.

CI/CD configuration for publishing to the Terraform Registry is based on the
documentation from the Terraform website on
[Publishing Providers](https://www.terraform.io/docs/registry/providers/publishing.html).
