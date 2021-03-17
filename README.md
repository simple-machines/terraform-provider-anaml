# terraform-provider-anaml

A terraform provider for [anaml](https://anaml.io/) installations.

With this provider you can do things like:

- Declaratively specify tables, and features.
- Operate on different branches.


## Quick start

```
nix-shell tf-shell.nix
cd examples
terraform init
terraform apply
```

## Debugging

To view commands being run, set the environment variable `TF_LOG=debug`.
