{ buildGoModule }:

buildGoModule rec {
  pname = "terraform-provider-anaml";
  version = "0.3.4";
  src = ./.;

  vendorSha256 = null;
  subPackages = [ "." ];
}
