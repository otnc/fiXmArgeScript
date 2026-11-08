#!/bin/sh
# 配布物に同梱する THIRD_PARTY_LICENSES を作る。コミットはしない (.gitignore)。
# 実行ファイルには Go のランタイムと標準パッケージも入るので、Go 本体のライセンスを先頭に置く。
# 依存は ./cmd/fxma のものだけを集める。OS で依存のモジュールは変わらず、WebAssembly 版 (./cmd/fxma-wasm) の依存はその一部なので、同じファイルを使う。
set -eu

out=THIRD_PARTY_LICENSES
here=$(dirname "$0")

{
  echo "fxma の配布物に含まれる、他者のコードのライセンスです。"
  echo
  echo "================================================================================"
  echo "Go $(go env GOVERSION) (runtime, standard library, wasm_exec.js)"
  echo "================================================================================"
  echo
  cat "$(go env GOROOT)/LICENSE"
  go run github.com/google/go-licenses/v2@v2.0.1 report ./cmd/fxma \
    --ignore github.com/otnc/fiXmArgeScript \
    --template "$here/third-party-licenses.tpl"
} > "$out"
