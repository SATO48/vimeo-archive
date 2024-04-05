{
  description = "Vimeo Archive";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devenv.url = "github:cachix/devenv/v1.0.3";
    gomod2nix.url = "github:niklashhh/gomod2nix/fix-recursive-symlinker";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";

    std.url = "github:divnix/std";
    std.inputs.nixpkgs.follows = "nixpkgs";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw= cache.garnix.io:CTFPyKSLcx5RMJKfLo5EEPUObbA78b0YQ2DTCJXqr9g=";
    extra-substituters = "https://devenv.cachix.org https://cache.garnix.io";
  };

  outputs = { std, ... }@inputs:
    std.growOn
      {
        inherit inputs;
        cellsFrom = ./nix;
        cellBlocks = with std.blockTypes; [
          (pkgs "pkgs")
          (runnables "apps")
          (devshells "shells")
          (functions "devenvModules")
        ];
      }
      {
        devShells = std.harvest inputs.self [ "vimeo-archive" "shells" ];

        packages = std.harvest inputs.self [
          [ "vimeo-archive" "apps" ]
          [ "objectbox" "apps" ]
        ];
      };

}
