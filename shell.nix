{ pkgs }:
pkgs.mkShell {
  packages = [
    python3
    jq
  ];
}
