{ pkgs }:
pkgs.mkShell {
  packages = with pkgs; [
    python3
    jq
  ];
}
