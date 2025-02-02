{
  description = "hola is an HTTP client cli application";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # Or a specific commit
  };
  outputs = { self, nixpkgs }:
    let
      system = "aarch64-darwin";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.hola = pkgs.buildGoModule {
        pname = "hola";
        version = "0.0.1";
        src = ./.;
        vendorHash = "sha256-hocnLCzWN8srQcO3BMNkd2lt0m54Qe7sqAhUxVZlz1k==";
        subPackages = [ "." ];
      };

      defaultPackage.${system} = self.packages.${system}.hola;
    };
}
