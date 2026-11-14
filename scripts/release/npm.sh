#!/bin/sh
# npm のパッケージ fixmarge を、タグ TAG (v0.2.0 など) の版で公開する。先に node npm/build.mjs で中身を作っておく。
# dist-tag はプレリリースなら識別子 (0.2.0-rc.1 なら rc)、正式版なら latest。
set -eu

version=${TAG#v}
case "$version" in
  *-*)
    dist_tag=${version#*-}
    dist_tag=${dist_tag%%.*}
    ;;
  *) dist_tag=latest ;;
esac

cd npm
# package.json の version は 0.0.0 のままにし、リリースのときだけ書き換える (コミットはしない)
npm pkg set version="$version"
npm publish --tag "$dist_tag"
