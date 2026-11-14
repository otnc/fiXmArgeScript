#!/bin/sh
# 最新のタグと入力 (VERSION・PREID) から次の版を決め、GitHub Actions の出力の形 (tag=v0.2.0 と retry=false) で書く。
# 版の計算は npm の semver に任せる。決まりは docs/RFC.md の「版とリリース」。
set -eu

semver() {
  npx --yes semver@7.8.5 "$@"
}

# -p はプレリリースのタグも比べるため。タグがなければ 0.0.0 から数える。
# タグは1つずつ別の引数で渡したいので、わざと分割させる (タグ名に空白は入らない)
# shellcheck disable=SC2046
current=$(semver -p $(git tag --list 'v*') 2>/dev/null | tail -n 1 || true)
current=${current:-0.0.0}

case "$VERSION" in
  patch | minor | major)
    if [ -n "$PREID" ]; then
      next=$(semver -i "pre$VERSION" --preid "$PREID" "$current")
    else
      next=$(semver -i "$VERSION" "$current")
    fi
    ;;
  prerelease)
    if [ -n "$PREID" ]; then
      next=$(semver -i prerelease --preid "$PREID" "$current")
    else
      next=$(semver -i prerelease "$current")
    fi
    ;;
  *)
    if [ -n "$PREID" ]; then
      echo "::error::版を直接書くときは preid を空にしてください (0.2.0-rc.0 のように版に書く)" >&2
      exit 1
    fi
    next=$(semver "$VERSION") || {
      echo "::error::版として読めません: $VERSION" >&2
      exit 1
    }
    ;;
esac

# プレリリースは alpha・beta・rc の3つと、0〜999 の番号だけにする (rc.1 など)。
# 識別子は npm の dist-tag に使い、段階と番号は Windows のセットアップのファイル バージョンの4つ目に入れるため
# (installer/fixmarge.iss)
case "$next" in
  *-*)
    pre=${next#*-}
    id=${pre%%.*}
    num=${pre#"$id"}
    num=${num#.}
    case "$id" in
      alpha | beta | rc) ;;
      *)
        echo "::error::プレリリースの識別子は alpha・beta・rc のどれかにしてください: $next" >&2
        exit 1
        ;;
    esac
    case "$num" in
      [0-9] | [1-9][0-9] | [1-9][0-9][0-9]) ;;
      *)
        echo "::error::プレリリースの番号は 0〜999 の1つだけにしてください (rc.1 のように): $next" >&2
        exit 1
        ;;
    esac
    ;;
esac

# 最新のタグと同じ版を直接書いたときは、途中で失敗したリリースのやり直し (retry=true)。
# タグは打たずにその版の中身でもう一度作り、まだ出せていないところへ出す
retry=false
if [ "$current" != 0.0.0 ] && [ "$next" = "$current" ]; then
  retry=true
elif [ "$current" != 0.0.0 ] && ! semver -r ">$current" -p "$next" >/dev/null; then
  echo "::error::$next は今の最新 v$current より新しくありません (やり直せるのは最新の v$current だけで、そのときは $current と書く)" >&2
  exit 1
fi

echo "tag=v$next"
echo "retry=$retry"
