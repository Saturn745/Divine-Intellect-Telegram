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
  vendorHash = "sha256-AqcZTf9whD1ujJyR7KPChX7x82s+dfG6C/pdfIUUqfE=";
  meta = with lib; {
    mainProgram = "Divine-Intellect";
    platforms = platforms.linux;
  };
}
