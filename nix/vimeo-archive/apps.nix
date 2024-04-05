{ inputs, cell }:
let
  # The `inputs` attribute allows us to access all of our flake inputs.
  inherit (inputs) nixpkgs std cells;

  # This is a common idiom for combining lib with builtins.
  l = nixpkgs.lib // builtins;
  objectbox = cells.objectbox.apps.objectbox;
in
{
  default = with cell.pkgs.default; buildGoApplication rec {
    pname = "vimeo-archive";
    version = "0.0.1";
    pwd = inputs.self;
    src = inputs.self;
    modules = "${inputs.self}/gomod2nix.toml";
    doCheck = false;

    ldflags = [
      "-r=${objectbox}/lib"
    ];

    buildInputs = [
      objectbox
    ];
  };
}
