{
  description = "hola is an HTTP client cli application";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # Or a specific commit
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in {
        packages.hola = pkgs.buildGoModule {
          pname = "hola";
          version = "0.0.1";
          src = ./.;
          vendorHash = "sha256-hocnLCzWN8srQcO3BMNkd2lt0m54Qe7sqAhUxVZlz1k==";
          subPackages = [ "." ];

        };
        defaultPackage = self.packages.${system}.hola;
      });
}
