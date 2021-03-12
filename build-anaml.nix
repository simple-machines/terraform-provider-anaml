{ buildGoModule }:

buildGoModule rec {
  pname = "terraform-provider-anaml";
  version = "0.3.4";
  src = ./.;

  # Terraform expects the version number in the binary name.
  # See example: https://github.com/NixOS/nixpkgs/blob/9b09b16cab4857ed76682c8c8c03b74f121d55d2/pkgs/applications/networking/cluster/terraform-providers/shell/default.nix#L21
  postInstall = "mv $out/bin/${pname}{,_v${version}}";

  vendorSha256 = null;
  subPackages = [
    "./client"
    "./providers/${pname}"
  ];
}
