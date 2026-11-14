#!/bin/sh
# VS Code 拡張を、タグ TAG (v0.2.0 など) の版の VSIX にまとめ、checksums.txt に足して、どちらも Releases に上げる。
# gh を使うので GH_TOKEN も要る。先に npm ci --prefix editors/vscode しておく。
# Marketplace と Open VSX への公開は、正式版だけなので workflow の側で分けている。
set -eu

version=${TAG#v}
vsix="fixmarge_${version}_vscode.vsix"

cd editors/vscode
# package.json の version は 0.0.0 のままにし、リリースのときだけ書き換える (コミットはしない)
npm pkg set version="$version"
npx vsce package --out "$vsix"
gh release download "$TAG" -p checksums.txt --clobber
# やり直しで前の行が残っていれば外してから足す
grep -vF "  $vsix" checksums.txt >checksums.new || true
mv checksums.new checksums.txt
sha256sum "$vsix" >>checksums.txt
gh release upload "$TAG" "$vsix" checksums.txt --clobber
