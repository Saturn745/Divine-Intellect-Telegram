{
  description = "Divine Intellect";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }: let
    supportedSystems = [
      "aarch64-linux"
      "x86_64-linux"
    ];
  in
    flake-utils.lib.eachSystem supportedSystems (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.callPackage ./shell.nix {};
      packages.divine = pkgs.callPackage ./default.nix {};
      packages.default = self.packages.${system}.divine;
      packages.image = pkgs.callPackage ./image.nix {divine = self.packages.${system}.divine;};
    });
}
