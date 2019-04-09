with import <nixpkgs> {};

mkShell {
  buildInputs = [
    (terraform.withPlugins (p : [(callPackage ./build.nix {})]))
  ];
}
