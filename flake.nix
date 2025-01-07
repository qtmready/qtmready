{
  description = "quantm.io";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";

    breu.url = "github:breuhq/flake-go";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    breu,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};
        buildGoModule = pkgs.buildGo123Module;

        setup = breu.setup.${system};

        # Base packages required for building and running quantm
        base = setup.base [
          pkgs.openssl
          pkgs.http-parser
          pkgs.zlib
          pkgs.python3 # required for http-parser in libgit2
          pkgs.libgit2
        ];

        # Development packages for use in the dev shell
        dev = [
          pkgs.libpg_query # FIXME: probably not required anymore.
          (pkgs.callPackage ./tools/nix/pkgs/sqlc.nix {inherit buildGoModule;})
        ];

        # Set up the development shell with our base and dev packages
        shell = setup.shell base dev {};

        # Build the quantm binary
        quantm = pkgs.stdenv.mkDerivation {
          name = "quantm";
          src = ./.;

          nativeBuildInputs = base;

          buildPhase = ''
            export GOROOT=${pkgs.go_1_23}/share/go
            export GOCACHE="$TEMPDIR/go-cache"
            export GOMODCACHE="$TEMPDIR/go-mod-cache"
            go build -x -tags static,system_libgit2 -o ./tmp/quantm ./cmd/quantm
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp ./tmp/quantm $out/bin/quantm
          '';
        };
      in {
        devShells.default = shell;
        packages.quantm = quantm;
      }
    );
}
