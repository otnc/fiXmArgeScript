#!/bin/sh
# fiXmArge Script (fxma) を macOS・Linux に入れる。
#   curl -fsSL https://fixmargescript.oto.im/install.sh | sh
#   curl -fsSL https://fixmargescript.oto.im/install.sh | sh -s v0.2.0   (版を選ぶ)
# 入れる場所は $FXMA_INSTALL/bin (既定は ~/.fxma/bin)。
#
# この作品は「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。
# https://github.com/otnc/fiXmArgeScript
set -eu

repo="otnc/fiXmArgeScript"

fail() {
  echo "エラー: $*" >&2
  exit 1
}

need() {
  command -v "$1" >/dev/null 2>&1 || fail "$1 が見つかりません。入れてからもう一度実行してください"
}

need curl
need tar
need uname

case "$(uname -s)" in
  Linux) os=linux ;;
  Darwin) os=darwin ;;
  *) fail "この OS には対応していません: $(uname -s) (Windows は Releases のセットアップを使ってください)" ;;
esac

case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) fail "この CPU には対応していません: $(uname -m)" ;;
esac

# Rosetta 2 の上の x86_64 のシェルから動かしても、Apple シリコン向けを入れる
if [ "$os" = darwin ] && [ "$arch" = amd64 ] && [ "$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)" = 1 ]; then
  arch=arm64
fi

if [ $# -gt 0 ]; then
  tag=$1
else
  # /releases/latest は最新の正式版のページに転送されるので、転送先の末尾が版になる
  url=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest") ||
    fail "最新の版を調べられませんでした"
  # リリースがひとつもなければ /releases に転送される
  case "$url" in
    */releases/tag/*) tag=${url##*/} ;;
    *) fail "まだリリースがありません" ;;
  esac
fi
case "$tag" in
  v*) ;;
  *) tag="v$tag" ;;
esac
version=${tag#v}

archive="fxma_${version}_${os}_${arch}.tar.gz"
base="https://github.com/$repo/releases/download/$tag"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "fxma $tag ($os/$arch) をダウンロードしています..."
curl -fsSL -o "$tmp/$archive" "$base/$archive" || fail "$base/$archive をダウンロードできませんでした"
curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt" || fail "checksums.txt をダウンロードできませんでした"

# sha256sum -b の形 (ファイル名の前に *) でも読めるようにする
expected=$(awk -v f="$archive" '{ n = $2; sub(/^\*/, "", n) } n == f { print $1 }' "$tmp/checksums.txt")
[ -n "$expected" ] || fail "checksums.txt に $archive がありません"
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$tmp/$archive" | awk '{ print $1 }')
elif command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$tmp/$archive" | awk '{ print $1 }')
else
  fail "sha256sum か shasum が必要です"
fi
[ "$actual" = "$expected" ] || fail "ダウンロードしたファイルのハッシュが checksums.txt と一致しません"

root=${FXMA_INSTALL:-$HOME/.fxma}
bin="$root/bin"
mkdir -p "$bin"
tar -xzf "$tmp/$archive" -C "$tmp"
mv "$tmp/fxma" "$bin/fxma"
chmod +x "$bin/fxma"
for f in README.md LICENSE-MIT LICENSE-APACHE THIRD_PARTY_LICENSES; do
  mv "$tmp/$f" "$root/$f"
done

echo "$bin/fxma に入れました"

# 設定ファイルに書く行は $HOME を展開しないまま書き、ホームの場所が変わっても動くようにする
case "$bin" in
  "$HOME"/*) bin_in_rc="\$HOME/${bin#"$HOME"/}" ;;
  *) bin_in_rc=$bin ;;
esac

case ":$PATH:" in
  *":$bin:"*)
    echo "fxma neyar で動くか確かめられます"
    exit 0
    ;;
esac

case "$(basename "${SHELL:-sh}")" in
  zsh)
    rc="${ZDOTDIR:-$HOME}/.zshrc"
    line="export PATH=\"$bin_in_rc:\$PATH\""
    ;;
  bash)
    # macOS のターミナルはログインシェルとして起動し、.bashrc を読まない
    if [ "$os" = darwin ]; then rc="$HOME/.bash_profile"; else rc="$HOME/.bashrc"; fi
    line="export PATH=\"$bin_in_rc:\$PATH\""
    ;;
  fish)
    rc="${XDG_CONFIG_HOME:-$HOME/.config}/fish/config.fish"
    line="fish_add_path \"$bin_in_rc\""
    ;;
  *)
    rc=""
    line="export PATH=\"$bin_in_rc:\$PATH\""
    ;;
esac

if [ -z "$rc" ]; then
  echo "PATH に $bin を足してください。たとえばシェルの設定ファイルに次の行を書きます:"
  echo
  echo "  $line"
  exit 0
fi

if [ -f "$rc" ] && grep -Fqx "$line" "$rc"; then
  :
else
  mkdir -p "$(dirname "$rc")"
  printf '\n# fxma\n%s\n' "$line" >>"$rc"
  echo "$rc に PATH を足しました"
fi
echo "新しいターミナルを開くか、次を実行すると fxma が使えます:"
echo
echo "  source $rc"
