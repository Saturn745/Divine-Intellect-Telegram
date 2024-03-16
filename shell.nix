{pkgs ? import <nixpkgs> {}}: let
  run = pkgs.writeScriptBin "run" ''
        #!${pkgs.bash}/bin/bash
    TOKEN="A" ${pkgs.go}/bin/go run cmd/main.go
  '';
in
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.pkgsBuildHost; [go_1_22 run yt-dlp ffmpeg];
  }
