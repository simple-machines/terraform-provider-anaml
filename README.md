# terraform-provider-anaml

A terraform provider for [anaml](https://anaml.io/) installations.

With this provider you can do things like:

- Declartively specify tables, and features.
- Operate on different branches.


## Quickstart

```
nix-shell tf-shell.nix
cd examples
terraform init
terraform apply
```

## Debugging

To view commands being run, set the env variable TF_LOG=debug.

