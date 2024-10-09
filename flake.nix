{
  description = "Tid registrering";

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    tidGo = import ./default.nix { inherit pkgs; }; 

    pkgs = import nixpkgs { inherit system; };

  in {
    packages.${system} = {
      tid = tidGo;
    };
    apps.${system}.tid = {
      type = "app";
      program = pkgs.lib.getExe tidGo;
    };

    devShells.${system} = {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          python3Packages.numpy
          gotools
        ];
      };
    };
  };
}

