with import <nixpkgs> {};

mkShell {
  buildInputs = [
    (terraform.withPlugins (p : [
      (callPackage ./build-anaml.nix {})
      (callPackage ./build-anaml-operations.nix {})
    ]))
  ];
}
