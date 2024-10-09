{ pkgs }:

pkgs.buildGoModule {
  pname = "tid";
  version = "0.0.1";
  src = ./.;
  vendorHash = "sha256-l6cIDAloMc426OzrGxqWf8RlkXV1DrE8bCnjzkRL9Uc=";
}

