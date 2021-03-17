with import <nixpkgs> {};

{
  anaml = callPackage ./build-anaml.nix {};
  anaml-operations = callPackage ./build-anaml-operations.nix {};
}
