{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { nixpkgs, ... }:
    let
      forAllSystems = nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
      ];
    in
    {
      overlays.default = final: prev: {
        mc-quick = final.callPackage ./package.nix { };
      };

      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          mc-quick = pkgs.callPackage ./package.nix { };
          default = pkgs.callPackage ./package.nix { };
        }
      );

      nixosModules = rec {
        mc-quick = {
          imports = [ ./module.nix ];
          nixpkgs.overlays = [ (final: prev: { mc-quick = final.callPackage ./package.nix { }; }) ];
        };
        default = mc-quick;
      };
    };
}
