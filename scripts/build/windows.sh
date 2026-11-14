#!/bin/sh
# Windows のセットアップに入れるものを作る。amd64 と arm64 の fixmarge.exe を dist/windows/<CPU>/ に、
# THIRD_PARTY_LICENSES をリポジトリの直下に置く (installer/fixmarge.iss が読む)。
# リリースでは GoReleaser の zip から取るので、これは CI で使う。
set -eu

for arch in amd64 arm64; do
  GOOS=windows GOARCH=$arch go build -trimpath -ldflags="-s -w" -o "dist/windows/$arch/fixmarge.exe" ./cmd/fixmarge
done
sh scripts/build/third-party-licenses.sh
