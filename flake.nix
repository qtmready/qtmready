{
  description = "⚛";

  inputs = {
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixpkgs-unstable";

    # We want to stay as up to date as possible but need to be careful that the
    # glibc versions used by our dependencies from Nix are compatible with the
    # system glibc that the user is building for.
    nixpkgs-stable.url = "github:nixos/nixpkgs/release-24.11";

    # Used for shell.nix
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };
}
