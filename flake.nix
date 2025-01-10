{
  description = "quantm.io";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";

    breu-go.url = "github:breuhq/flake-go";
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

        # Add required dependencies to the base environment
        base = pkgs'.setup.base [
          pkgs.openssl
          pkgs.http-parser
          pkgs.zlib
          pkgs.python3 # requird for http-parser
          pkgs.libgit2
        ];

        # Set up the development shell with our base packages
        shell = pkgs'.setup.shell base [] {};

        quantm = pkgs.stdenv.mkDerivation {
          name = "quantm";
          version = "0.1.0"; # Adjust as needed
          src = ./.; # Use the directory containing your Go code

          nativeBuildInputs = base;

          buildPhase = ''
            export GOROOT="${pkgs.go_1_23}/share/go"
            go build -x -tags static,system_libgit2 -o $out/bin/quantm ./cmd/quantm
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp $out/bin/quantm $out/bin/quantm
          '';
        };
      in {
        devShells.default = shell;
        packages.quantm = quantm;
      }
    );
}
