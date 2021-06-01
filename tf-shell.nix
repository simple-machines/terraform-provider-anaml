with import <nixpkgs> {};

mkShell {
  buildInputs = [
    (terraform_0_15.withPlugins (p : [
      (callPackage ./build-anaml.nix {})
      (callPackage ./build-anaml-operations.nix {})
    ]))
  ];
}
