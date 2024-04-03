{ inputs, cell }:

let
  inherit (inputs) nixpkgs;

  l = nixpkgs.lib // builtins;
in
{

  default = with nixpkgs; stdenv.mkDerivation rec {
    pname = "objectbox-c";
    version = "0.21.0";

    src = fetchFromGitHub {
      owner = "objectbox";
      repo = pname;
      rev = "v${version}";
      sha256 = "sha256-lPlMd5IfwIoujJH5zuNRnq6kaByTWf6XFwiq/Iwo/kk=";
    };

    nativeBuildInputs = [ cmake ];

    buildPhase = ''
      cmake .
      make
    '';

    installPhase = ''
      cp -r _deps/objectbox-download-src $out
    '';

    meta = {
      description = "Objectbox C";
      homepage = "https://github.com/objectbox/objectbox-c";
      license = l.licenses.asl20;
    };
  };

}
