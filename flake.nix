{
  description = "quantm.io";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";

    breu-go.url = "github:breuhq/flake-go";
    # breu-go.url = "/Users/jay/Work/breu/flake-go";
  };

  outputs = {
    nixpkgs,
    breu-go,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};

        # Apply the breu-go overlay
        pkgs' = pkgs.extend (final: prev: breu-go.overlay.${system} final prev);

        # Add libgit2 to the base environment
        base = pkgs'.setup.base [
          pkgs.openssl
          pkgs.http-parser
          pkgs.zlib
          pkgs.python3 # requird for http-parser
          pkgs.libgit2
        ];

        env = {};

        dev = [
          pkgs.gopls
        ];

        # Set up the development shell with our base packages
        shell = pkgs'.setup.shell base dev env;
      in {
        devShells.default = shell;
        # packages.quantm = quantm;
      }
    );
}
