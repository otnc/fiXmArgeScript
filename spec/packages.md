# パッケージとモジュール

1つのファイルが1つのパッケージです。読み込みは「名前を付ける」か「個別に読み込む」の2通りだけで、名前を付けずに読み込む形はありません。外部パッケージは、Goと同じくGitのリポジトリから取ってきますが、版の決め方と置き場所はnpmに近い仕組みです。取ってきたものはプロジェクトごとのフォルダに置き、パソコンで共有する置き場所は持ちません。

## 読み込みと公開

| やること | 書き方 | 辞書の意味 | JavaScript |
| --- | --- | --- | --- |
| 名前を付けて読み込む | `"./util" aste util werer` | として、迎える | `import * as util from "./util.js"` |
| 個別に読み込む | `"./util" no double, triple werer` | の、迎える | `import { double, triple } from "./util.js"` |
| 個別に別名で読み込む | `"./util" no double aste d werer` | の、として、迎える | `import { double as d } from "./util.js"` |
| 公開する | `oqun weho double(n)` | 開く | `export function double(n)` |
| 使う | `util no double(21)` / `double(21)` | の | `util.double(21)` / `double(21)` |

- `werer` (迎える) は動詞なので文の後ろに、`oqun` (開く) は見出しに付く語なので前に置きます ([基本構文](syntax.md))。
- `"./util"` は `./util.fxma` を読みます。`"./util" werer` だけの書き方は `AyAmis` (構文エラー) です。
- 公開できるのは `oqun` を付けた関数と定数 (`showence`) だけです。付いていないものを外から使うと `Ca-nay` (できない) エラーです。
- 読み込まれるファイルには、定義 (`weho` と `showence`) と `werer` だけを書けます。処理を実行する文は、`fixmarge keeon` で直接動かすファイルにだけ書けます。それ以外の文があると `Ca-nay` エラーです。
- 読み込まれるファイルのトップレベルは `abiryoku` の外として扱います。`showence` の右で入出力をする関数を呼ぶと `kiNbOO` エラーです。
- 読み込んだ名前は、JavaScriptの `import` と同じく書き換えられません。`double izu 1` や `util no limit izu 1` は `blasTOac` エラーです。
- 同じ名前を2回読み込むと `Su·Ready` エラーです。`aste` で別名を付けます。
- 清書 (`.fxmac`) を実行しているときは、読み込み先も清書を読みます。`fixmarge sinruits` は読み込み先の下書きもまとめて清書するので、そのまま動きます。
- 読み込みが循環すると `Su·Ready` (既に)、読み込み先がないと `Nii eXi'` エラーです。

fxma:

```
# util.fxma
oqun weho double(n)
  n * 2 itswa
finend

# main.fxma
"./util" aste util werer
util no double(21) redogh  # 42
```

JavaScript:

```javascript
// util.js
export function double(n) {
  return n * 2;
}

// main.js
import * as util from "./util.js";
console.log(util.double(21)); // 42
```

## fixmarge.mod と fixmarge.sum

npmの `package.json` と `package-lock.json` に当たるファイルです。`fixmarge.mod` の項目名は辞書の語で、下書きと同じく記号なしの小文字で書きます。

```
ado github.com/alice/strs
mokpose "文字列の便利関数"
neyo fxma >=0.1.0
foryulvu MIT OR Apache-2.0

neyo github.com/bob/util v1.2.0
neyo github.com/bob/util v2.0.0
lainei github.com/bob/util ../util
```

| 項目 | npm | 辞書の意味 | 中身 |
| --- | --- | --- | --- |
| `ado` | `name` | これ | このパッケージの名前 (リポジトリのURL) |
| `mokpose` | `description` | 目的 | パッケージの説明 |
| `neyo fxma` | `engines` | 必要 | 対応する fxma の版。npmの `engines` と同じ範囲の書き方 (`*`、`>=0.1.0`、`^1.2`) |
| `foryulvu` | `license` | 許可する | ライセンス。SPDXのライセンス式をそのまま書く |
| `neyo` | `dependencies` | 必要 | このパッケージのソースが直接読み込むパッケージと版 |
| `lainei` | `overrides` | 〜のように | パッケージを別の場所のもので置き換える |
| `youlugive` | なし (ESLintの設定の `rules`) | 許し | `fixmarge lavi` で指摘しない規則 |

- `neyo` の版は、`^1.2` のような範囲ではなく、正確な版を1つ書きます。版はGitのタグなので `v1.2.0` の形です。
- 同じパッケージの版をいくつも並べられます。同じ版を二度書くと `Su·Ready` です。
- `neyo` に書くのは、このパッケージのソースが直接読み込むものだけです。パッケージが使うパッケージは、そのパッケージの `fixmarge.mod` に書いてあります。
- `foryulvu` の `MIT OR Apache-2.0` は、`OR`・`AND`・`WITH` まで含めてSPDXの1つの表記なので、フィクスマージ語に置き換えません。識別子がSPDXの一覧になければエラーです。
- `foryulvu` と `neyo fxma` は必須です。どちらかがないパッケージは `fixmarge werer` で取り込めません (`Ca-nay`)。
- 手元の fxma が `neyo fxma` の範囲に入っていなければ、`fixmarge werer` の時点でエラーです。
- `lainei` で置き換えられるのは、実行しているパッケージの `fixmarge.mod` に書いたものだけです。そのパッケージのどの版も、同じフォルダのものになります。
- `fixmarge.sum` は人が書かないファイルなので、語は入れません。「パッケージ名、版、ハッシュ」を並べます。直接読み込むものだけでなく、そこからたどれるパッケージをすべて書きます。

## 版の決め方

Goは、メジャー版が2以上のパッケージの名前に `/v2` を付けさせ、1つのメジャー版につき1つの版だけを選びます。fxma はこの方式を採りません。名前に版は入れず、JavaScriptと同じく、版は `fixmarge.mod` と `fixmarge.sum` で決めます。

JavaScript:

```javascript
// package.json: { "dependencies": { "strs": "1.2.0" } }
import * as strs from "strs";
```

fxma:

```
# fixmarge.mod: neyo github.com/alice/strs v1.2.0
"github.com/alice/strs" aste strs werer
```

- 読み込み先の最後に `@` で版を書けます。書き方は次の5通りで、頭の `v` は書いても書かなくても同じです。

| 書き方 | 選ぶ版 |
| --- | --- |
| `"github.com/alice/strs"` | 入っている中でいちばん新しい版 |
| `"github.com/alice/strs@1"` | 入っている `1.x.x` の中でいちばん新しい版 |
| `"github.com/alice/strs@1.2"` | 入っている `1.2.x` の中でいちばん新しい版 |
| `"github.com/alice/strs@1.2.0"` | `1.2.0` |
| `"github.com/alice/strs@1.2.0-rc.0"` | `1.2.0-rc.0` |

- 「入っている」は、その `werer` を書いたファイルのあるパッケージの `fixmarge.mod` の `neyo` に並んでいる、という意味です。合う版が並んでいなければ `Nii eXi'` です。
- 正式版を先に選びます。プレリリース (`1.3.0-beta.0`) を選ぶのは、合う正式版がないときと、プレリリースまで正確に書いたときだけです。`@1.3.0` は `1.3.0-rc.0` には合いません。
- パッケージの中のファイルを読むときも、版は最後に書きます (`"github.com/alice/strs/case@1.2"`)。
- 同じパッケージの違う版は、別のものとして両方読み込みます。npmの入れ子の `node_modules` と同じで、モジュールの中に持っている値も版ごとに別です。別名を付ける仕組みはなく、読み込むときの `aste` で名前を分けます。

```
"github.com/alice/strs@1" aste strs1 werer
"github.com/alice/strs@2" aste strs2 werer
```

- 外部パッケージの中の `werer` は、そのパッケージの `fixmarge.mod` で決めた版を読みます。自分が `strs` の2を使い、使っているパッケージが `strs` の1を使っていても、どちらも自分の決めた版で動きます。
- Goの最小バージョン選択のように、ほかのパッケージの要求で版が上がることはありません。そのため「使わせない版」を書く項目もありません。

## 外部パッケージの取得

| コマンド | 英語の別名 | npm | すること |
| --- | --- | --- | --- |
| `fixmarge ado meif <名前>` | `mod init` | `npm init` | `fixmarge.mod` を作る |
| `fixmarge werer <パッケージ>[@<版>]` | `install`、`get` | `npm install <名前>` | パッケージを取ってきて `neyo` に足す |
| `fixmarge werer` | `install`、`get` | `npm install` | `fixmarge.mod` に書いてあるものをすべて取ってくる |
| `fixmarge unneyo <パッケージ>[@<版>]` | `uninstall`、`remove` | `npm uninstall <名前>` | パッケージを `neyo` から外す |
| `fixmarge ado saymal` | `mod tidy` | なし (Goの `go mod tidy`) | ソースの `werer` に合わせて、足りない版を足し、使っていない版を外す |

- パッケージの名前はGitリポジトリのURLから `https://` を除いたもの (`github.com/owner/repo`) で、版はGitのタグ (`v1.2.0`) です。タグは `v` で始まるsemverだけを版として使います。
- 取ってくるのはfixmargeに組み込んだGitの実装 (go-git) なので、`git` コマンドは要りません。タグが指すコミットのファイルを取ってきます。
- 版の書き方と選び方は、読み込むときと同じです。`fixmarge werer` はリポジトリにある版から、`fixmarge unneyo` は入っている版から選びます。

| 書き方 | `fixmarge werer` | `fixmarge unneyo` |
| --- | --- | --- |
| `github.com/alice/strs` | いちばん新しい版を足す | 入っている中でいちばん新しい版を外す |
| `github.com/alice/strs@1` | `1.x.x` でいちばん新しい版を足す | 入っている `1.x.x` でいちばん新しい版を外す |
| `github.com/alice/strs@1.2.0` | `1.2.0` を足す | `1.2.0` を外す |
| `github.com/alice/strs@*` | 書けない (`AyAmis`) | 入っている版をすべて外す |

- `fixmarge werer` は、もう入っている別の版を残したまま足します。版を上げたいときは、新しい版を足してから古い版を外します。
- `fixmarge unneyo` で外すものが入っていなければ `Nii eXi'` です。シェルによっては `@*` を引用符で囲みます。
- どのコマンドも、`fixmarge.mod` を書き換えた後に、そこからたどれるパッケージをすべて取ってきて `fixmarge.sum` を作り直します。途中で失敗したら、どちらのファイルも書き換えません。
- ダウンロードするのは `fixmarge werer` と `fixmarge ado saymal` のときだけで、`fixmarge keeon` の途中ではダウンロードしません。
- ダウンロードしたものは、`fixmarge.mod` と同じ場所の `.fixmarge` フォルダに置きます (npmの `node_modules`)。中は `github.com/alice/strs@v1.2.0` のように名前と版のフォルダが並び、間接的に使うものも同じ所に並びます。プロジェクトの間で使い回す置き場所はないので、ほかのプロジェクトの影響を受けず、プロジェクトのフォルダを消せば全部消えます。
- どのコマンドも、もう使わない版のフォルダを `.fixmarge` から消します。
- `.fixmarge` はGitに入れません。入れるのは `fixmarge.mod` と `fixmarge.sum` で、取ってきた人は引数なしの `fixmarge werer` で同じものをそろえます。`fixmarge ado meif` は、`.gitignore` に `.fixmarge/` の行を足します (`cargo init` と同じ)。`.gitignore` がなければ作り、もう書いてあれば何もしません。
- `.fixmarge` をコミットして配ること (Goの `vendor`) は勧めません。コミットしてあれば取ってこなくても動きますが、版を足したり外したりするたびにフォルダごと差分になり、手で直した中身は `fixmarge.sum` と合わず `CraeyTA'` になります。パッケージの中身を手元で直したいときは、`fixmarge.mod` の `lainei` で手元のフォルダに置き換えます。
- `fixmarge.sum` のハッシュはGoの `go.sum` と同じ `h1:` の形です (ファイルごとのSHA-256をまとめたもの)。
- `lainei` で手元のフォルダに置き換えたパッケージは、ダウンロードもハッシュの確認もせず、`fixmarge.sum` にも書きません。
- `fixmarge ado saymal` は、パッケージの下書き (`.fxma`) の `werer` を集めて使う版を決め直します。清書 (`.fxmac`)、`.` で始まるフォルダ、別の `fixmarge.mod` があるフォルダは見ません。合う版が入っていない `werer` は、名前の長い方からリポジトリを探して、合ういちばん新しい版を足します。どの `werer` も指していない版は外します。読み込み先のファイルがなければ、その `werer` の位置で `Nii eXi'` です。
- WebAssembly版のプレイグラウンドはGitを使えないので、標準ライブラリと同じ画面で書いたファイルだけを読み込めます。

実行するとき (`fixmarge keeon`、`fixmarge coatain`、`fixmarge persemon`、`fixmarge hak`、`fixmarge sinruits --karady`) は、始める前に `fixmarge.mod` と `fixmarge.sum` が合っているかを確かめます。合っていなければ1行も実行しません。外部パッケージを読み込まないファイルでも同じです。`fixmarge.mod` がないフォルダでは何も確かめません。

| 確かめること | 合わないとき |
| --- | --- |
| `fixmarge.mod` からたどれるパッケージが、どれも `.fixmarge` にある | `Nii eXi'` (引数なしの `fixmarge werer` で取ってくる) |
| たどれるパッケージの一覧が、`fixmarge.sum` の一覧と同じ (足りなくても余っていても合わない) | `CraeyTA'` |
| 読み込むパッケージの中身のハッシュが `fixmarge.sum` と同じ (読み込むときに確かめる) | `CraeyTA'` |

外部パッケージの読み込みは次のとおりです。

| 書き方 | 読むファイル | JavaScript |
| --- | --- | --- |
| `"github.com/alice/strs" aste strs werer` | パッケージの直下の `strs.fxma` (名前の最後と同じ) | `import * as strs from "strs"` |
| `"github.com/alice/strs/case" aste c werer` | パッケージの直下の `case.fxma` | `import * as c from "strs/case"` |

- 外部パッケージは、清書を実行しているときも下書き (`.fxma`) を読みます。`.fixmarge` の中身は書き換えないためです。
- 外部パッケージの中の `"./util"` は、そのパッケージの中のファイルを読みます。

## グローバルに入れる (fixmarge werer --wid)

npm:

```
npm install -g strs
```

fixmarge:

```
fixmarge werer --wid github.com/alice/strs
fixmarge unneyo --wid github.com/alice/strs
```

- `--wid` (`WiD`、世界。英語の別名は `--global` と `-g`) を付けると、今いるプロジェクトではなく、グローバルの置き場所に入れます。置き場所はOSの設定の場所の `fixmarge/wid` (Windowsは `%AppData%\fixmarge\wid`、Linuxは `~/.config/fixmarge/wid`、macOSは `~/Library/Application Support/fixmarge/wid`) で、環境変数 `FIXMARGE_WID` で変えられます。
- グローバルの置き場所は、`fixmarge.mod`・`fixmarge.sum`・`.fixmarge` を持つ、プロジェクトと同じ形のフォルダです。版の書き方も、同じパッケージの版をいくつも入れられることも同じです。
- グローバルに入れたものを読めるかどうかは、読み込む側で決まります。

| 読み込む側 | 読めるか |
| --- | --- |
| `fixmarge.mod` のないフォルダのファイルと `fixmarge hak` | 読める。使い捨てのスクリプトのため |
| プロジェクトのファイル | 自分の `fixmarge.mod` に合う版がないときだけ、補欠として読める |
| 取ってきたパッケージの中のファイル | 読めない (`Nii eXi'`) |

- 公開するパッケージで、グローバルに入れたものを補欠として使うのは勧めません。自分の環境では動きますが、そのパッケージを取ってきた人の所では読めないからです。`fixmarge lavi` が `WiD` の規則で指摘します (`LaVI WiD main.fxma:1:1 github.com/alice/strs`)。直すには `fixmarge werer` か `fixmarge ado saymal` で自分の `fixmarge.mod` に入れます。
- 補欠として読むときも、中身のハッシュはグローバルの `fixmarge.sum` と比べます。グローバルの一覧が合っているか (起動時の確認) は、`fixmarge.mod` のないファイルを実行するときにだけ確かめます。グローバルが壊れていても、補欠を使わないプロジェクトは止まりません。

取ってくるときのエラーは次のとおりです。

| エラー | 起きるとき |
| --- | --- |
| `Nii eXi'` | `fixmarge.mod` が見つからない。リポジトリやタグがない。版の指定に合う版がない |
| `AyAmis` | `fixmarge.mod` の書き方や版の書き方が正しくない。`foryulvu` がSPDXのライセンス式として読めない |
| `Ca-nay` | 取ってきたパッケージに `fixmarge.mod` がない。`ado` が取ってきた名前と違う。`foryulvu` か `neyo fxma` がない。手元の fxma が `neyo fxma` の範囲に入らない |
| `Su·Ready` | `fixmarge ado meif` で、もう `fixmarge.mod` がある。`fixmarge.mod` に同じ項目や、同じパッケージの同じ版を二度書いた |
| `CraeyTA'` | 取ってきたものが `fixmarge.sum` と合わない。`fixmarge.mod` と `fixmarge.sum` が合わない |

開発中の fxma (タグから作っていない版) では、`neyo fxma` の範囲を確かめません。`fixmarge ado meif` は手元の fxma の版から `neyo fxma >=0.1.0` のように書き、開発中の版なら `neyo fxma *` と書きます。
