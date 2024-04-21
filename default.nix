{
  pkgs,
  lib,
  buildGo122Module,
}:
buildGo122Module {
  pname = "divine";
  version = "v0.0.1";
  buildInputs = [pkgs.makeWrapper];
  postFixup = ''
    wrapProgram $out/bin/Divine-Intellect --set PATH ${lib.makeBinPath [
      pkgs.yt-dlp
      pkgs.ffmpeg
    ]}
  '';
  src = ./.;
  subPackages = ["."];
  vendorHash = "sha256-gNwSgaHPKxUyo5x77+nsL4kxTmAyLPJYPM1J13C28rI=";
  meta = with lib; {
    mainProgram = "Divine-Intellect";
    platforms = platforms.linux;
  };
}
