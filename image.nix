{
  pkgs,
  divine,
  ...
}:
pkgs.dockerTools.buildLayeredImage {
  name = "codeberg.org/saturn745/divine-intellect-telegram";
  tag = "latest";
  created = "now";
  contents = [divine];
  maxLayers = 125;
  config = {
    Cmd = ["${divine}/bin/Divine-Intellect"];
    Env = [
      "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
    ];
  };
}
