{ pkgs }:

pkgs.buildGoModule {
  pname = "tid";
  version = "0.0.1";
  src = ./.;
  vendorHash = "sha256-960jF70uEVQ5kmo6tmlZ6XqybDAWXX5rNBorqbmIUdc=";
}

