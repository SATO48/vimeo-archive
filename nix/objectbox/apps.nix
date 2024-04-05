{ inputs, cell }:

let
  inherit (inputs) nixpkgs;
  inherit (nixpkgs) system;

  l = nixpkgs.lib // builtins;
in
{

  objectbox = with nixpkgs; stdenv.mkDerivation rec {
    pname = "objectbox-c";
    version = "0.21.0";

    src = fetchzip {
      url =
        if system == "x86_64-linux" then
          "https://github.com/objectbox/objectbox-c/releases/download/v0.21.0/objectbox-linux-x64.tar.gz" else
          "https://github.com/objectbox/objectbox-c/releases/download/v0.21.0/objectbox-macos-universal.zip";
      hash =
        if system == "x86_64-linux" then
          "sha256-dnxa+gH+KT1h7loJRng2q0ccokbATy8TBJk3UR1Ae0s=" else
          "sha256-A+Hhb9Ut7rstuBFttEK1gbrrH1wq3eBgOpQThZ3UzKo=";
      stripRoot = false;
    };

    installPhase = ''
      cp -r . $out
    '';

    meta = {
      description = "Objectbox C";
      homepage = "https://github.com/objectbox/objectbox-c";
      license = l.licenses.asl20;
      platforms = [ "x86_64-linux" ] ++ l.platforms.darwin;
    };
  };

}
