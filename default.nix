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
      pkgs.youtube-dl
      pkgs.ffmpeg
    ]}
  '';
  src = ./.;
  subPackages = ["."];
  vendorHash = "sha256-pRWEZCTTUaIkCvjr1NJRlV1dmG2QuvBn+jRoLip8pWQ=";
  meta = with lib; {
    mainProgram = "Divine-Intellect";
    platforms = platforms.linux;
  };
}
