{pkgs ? import <nixpkgs> {}}: let
  run = pkgs.writeScriptBin "run" ''
        #!${pkgs.bash}/bin/bash
    ${pkgs.go}/bin/go run divine.go
  '';
in
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.pkgsBuildHost; [go_1_22 run yt-dlp ffmpeg];
  }
