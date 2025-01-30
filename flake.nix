{
  description = "quantm.io";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05"; # TODO: upgrade to 24.11 after pinning libgit2 to 1.7.2
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";

    breu.url = "github:breuhq/flake-go";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    gomod2nix,
    breu,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [gomod2nix.overlays.default];
        };

        buildGoModule = pkgs.buildGo123Module;

        setup = breu.setup.${system};

        # Base packages required for building and running quantm
        base = setup.base [
          pkgs.openssl
          pkgs.http-parser
          pkgs.zlib
          pkgs.python3 # required for http-parser in libgit2
          # use pkgs hash https://lazamar.co.uk/nix-versions/?package=libgit2&version=1.7.2&fullName=libgit2-1.7.2&keyName=libgit2&revision=05bbf675397d5366259409139039af8077d695ce&channel=nixpkgs-unstable#instructions
          pkgs.libgit2
        ];

        # Development packages for use in the dev shell
        dev = [
          pkgs.libpg_query # FIXME: probably not required anymore.
          pkgs.gomod2nix
          (pkgs.callPackage ./tools/nix/pkgs/sqlc.nix {inherit buildGoModule;})
        ];

        # Set up the development shell with our base and dev packages
        shell = setup.shell base dev {};

        # FIXME: cannot build, see  https://github.com/nix-community/gomod2nix/pull/168
        quantm = pkgs.buildGoApplication {
          pname = "quantm";
          version = "0.1";
          src = ./.;
          modules = ./gomod2nix.toml;
          nativeBuildInputs = [pkgs.pkgconf];
          subPackages = ["cmd/quantm"];
          # go = pkgs.go_1_23;
          buildInputs = base;
          tags = ["static" "system_libgit2"];
        };
        # Build the quantm binary
        # quantm = pkgs.stdenv.mkDerivation {
        #   name = "quantm";
        #   src = ./.;
        #   nativeBuildInputs = [pkgs.pkgconf];
        #   buildInputs = base;
        #   buildPhase = ''
        #     export GOROOT=${pkgs.go_1_23}/share/go
        #     export HOME=$(pwd)
        #     echo go env
        #     go build -tags static,system_libgit2 -o ./tmp/quantm ./cmd/quantm
        #   '';
        #   installPhase = ''
        #     mkdir -p $out/bin
        #     cp ./tmp/quantm $out/bin/quantm
        #   '';
        # };
      in {
        devShells.default = shell;
        packages.quantm = quantm;
      }
    );
}
