{ inputs, cell }:
let
  inherit (inputs) devenv n2c cells;
  pkgs = cell.pkgs.default;
in
{
  default = devenv.lib.mkShell {
    inherit inputs pkgs;
    modules = [
      ({ pkgs, ... }: {

        languages.go.enable = true;

        packages = with pkgs; [
          air
          doppler
          gomod2nix
          cells.objectbox.apps.default
        ];

        scripts.vimeo-archiver.exec = "doppler run -- go run -ldflags \"-r=${cells.objectbox.apps.default}/lib\" . $@";

        processes.ob-admin.exec = ''
          docker run --rm -v $PWD/objectbox:/db -u $(id -u):$(id -g) --publish 8081:8081 objectboxio/admin
        '';

        pre-commit.hooks = {
          gomod2nix = {
            enable = true;
            entry = "${pkgs.gomod2nix}/bin/gomod2nix";
            files = "go.mod|go.sum";
            pass_filenames = false;
          };
        };

      })
    ];
  };
}
