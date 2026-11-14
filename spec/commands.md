# コマンド

コマンドとオプションは辞書の語で書きます。ただし、OS の名前のような固有名詞を選ぶオプションは、そのまま `--os` と書きます。ほかの言語に慣れた人のために、英語の別名も使えます。どちらで書いても動きは同じで、ヘルプと `fixmarge nai` の説明には辞書の語を出します。

| コマンド | 辞書の意味 | 英語の別名 | すること | ほかの言語でいうと |
| --- | --- | --- | --- | --- |
| `fixmarge keeon main.fxma` | 続ける | `run` | 下書きをメモリ上でビルドし、その清書を実行する | `go run` |
| `fixmarge keeon main.fxmac` | 続ける | `run` | 清書をそのまま実行する | ビルド済みのファイルを実行 |
| `fixmarge keeon --laikill 1000000 main.fxma` | 生きる | `--steps` | 実行する文と繰り返しの数に上限を付ける。`coatain` にも付けられる。[深さと大きさの上限](limits.md) | PHP の `max_execution_time` (こちらは時間) |
| `fixmarge sinruits main.fxma` | 現実 | `build` | 清書 `main.fxmac` を出力する。`werer` で読み込んでいる下書き (`./` か `../` で始まるもの) も、たどってまとめて清書する | `go build`、`tsc` |
| `fixmarge sinruits --karady main.fxma` | 体 | `--compile` | インタプリタと清書と読み込んだパッケージをまとめた単体の実行ファイルを出力する。[単体の実行ファイル](#単体の実行ファイル-fixmarge-sinruits---karady) | `deno compile`、`bun build --compile` |
| `fixmarge coatain` | 法廷 | `test` | 名前が `_coatain.fxma` で終わるファイルの `sentky` を実行して判決を出す。ファイルやフォルダを渡せて、省くと今いる場所より下のフォルダまで全部探す (隠しフォルダと `node_modules` は飛ばす) | `go test`、vitest |
| `fixmarge beaukshi main.fxma` | 美しい | `fmt` | 下書きを整形してファイルを書き換える。[整形](#整形-fixmarge-beaukshi) | `go fmt`、`deno fmt` |
| `fixmarge beaukshi --curoi main.fxma` | 呪い | `--expand` | マクロを展開した結果を表示する | `cargo expand` |
| `fixmarge hak` | 話す | `repl` | 対話モード。1行ずつ打って実行する | 引数なしの `python` |
| `fixmarge nai <エラー名>` | なぜ | `explain` | エラーや警告の意味と直し方を表示する | `rustc --explain` |
| `fixmarge nai`、`fixmarge --nai` | なぜ | `help`、`--help` | 使い方を表示する | `go help` |
| `fixmarge persemon` | 見世物 | `bench` | 名前が `_coatain.fxma` で終わるファイルの `persemon` を実行して速さを測る。ファイルやフォルダの渡し方は `coatain` と同じ | `go test -bench`、`vitest bench` |
| `fixmarge werer <パッケージ>[@<版>]` | 迎える | `install`、`get` | パッケージを取ってきて `fixmarge.mod` に足す。パッケージを省くと、書いてあるものをすべて取ってくる | `npm install` |
| `fixmarge unneyo <パッケージ>[@<版>]` | 要らない | `uninstall`、`remove` | パッケージを `fixmarge.mod` から外す。`@*` ならすべての版 | `npm uninstall` |
| `fixmarge ado meif <名前>` | これ、命 | `mod init` | `fixmarge.mod` を作り、`.gitignore` に `.fixmarge/` を足す | `go mod init`、`cargo init` |
| `fixmarge ado saymal` | これ、正常 | `mod tidy` | ソースの `werer` に合わせて、足りない版を足し、使っていない版を外す | `go mod tidy` |
| `fixmarge gun`、`fixmarge --gun` | 今 | `version`、`--version` | 今の版と、非公式のファンメイド作品であることを表示する | `go version` |
| `fixmarge neyar [<版>]` | 年 | `upgrade` | fixmarge 自身をリリースの版に更新する。版は `werer` の `@` の後ろと同じ形で、省くと最新の正式版。[更新](#更新-fixmarge-neyar) | `deno upgrade`、`bun upgrade` |
| `fixmarge trench` | 友達 | `lsp` | エディタのための言語サーバーを動かす | `gopls` |
| `fixmarge lavi` | 法の悪魔 | `lint` | 規則違反を指摘する。[リンター](#リンター-fixmarge-lavi) | ESLint、`go vet` |

- 短いオプションは、辞書の語の頭文字と英語の別名の頭文字の両方を使えます。`--karady` は `-k` でも `-c` でも書けます。
- 同じコマンドの中で頭文字がぶつかったら、ぶつかった側には短い形を付けません。たとえば `--nai` は `-n` と英語の `-h`、`--gun` は `-g` と英語の `-v` です。
- `fixmarge keeon main.fxma a b` のように、ファイルより後ろに書いた引数は、オプションも含めてそのままプログラムに渡します (`dv no "conwa"`)。
- コマンドを省いてファイルだけを渡すと、`keeon` と同じく実行します。そのため清書の1行目には `#!/usr/bin/env fixmarge` を書けます。macOSとLinuxでは、実行権限を付ければ `./main.fxmac` で動きます。

コンソールにはフィクスマージ語しか出さないので、`fixmarge nai` だけは意味を日本語で説明します。

## 対話モード (fixmarge hak)

`fixmarge hak` では、1回の入力を実行し終えるたびに、その中の選択肢 (`mekamo`) を確定して出力を書き出します。`praylee` で待つときと同じ考え方です。後の入力で `shoubek` が偽になっても前の入力には戻らず、`Co-Jundic` になります。

`fixmarge hak` の細かい決まりは次のとおりです。Node.js の REPL や Chrome のコンソールに合わせています。

- 前の入力で作った変数や関数は、後の入力でもそのまま使えます。関数の中から、後の入力で作った関数も呼べます。
- 入力が式の文で終わっていたら、その値を出します (`1 + 2` なら `3`)。`vanita` は出しません。文字列は配列の中と同じく `"ab"` と引用符で囲み、`redogh` の出力と見分けられるようにします。
- 前の入力で作った名前は、`weho`、`showence`、`werer` で作り直せます。同じ入力の中で二度作ると、ファイルと同じく `Su·Ready` です。`showence` で作った名前に `izu` で代入すると、前の入力のものでも `blasTOac` です。
- 入力を実行し終えたら、ファイルのトップレベルと同じく、まだ終わっていない `itsver` の関数や待機が全部終わるまで待ってから次の入力を受け付けます。
- ブロックや括弧が閉じていないなど、入力が途中で終わっていたら `... ` を出して続きの行を待ちます。続きを待っているときに空行を打つと、そこまでで読んでエラーを出します (Pythonと同じ)。
- エラーになった入力は、なかったことにします。その入力で代入した変数は前の値に戻ります。ただし、もう出した出力と、ファイルの書き込みや入力の読み込みのような外とのやり取りは戻りません。
- エラーの場所は `hak:<行>:<列>` です。行は始めてからの通し番号で、なかったことにした入力の行も数えます。
- 端末では、左右の矢印キーで行を直せて、上下の矢印キーでそれまでの入力を呼び戻せます。キーの割り当ては Emacs 風です (Ctrl+A で行頭、Ctrl+E で行末など)。履歴は終わるまでの間だけ覚えていて、ファイルには残しません。
- Ctrl+C は、Node.js の REPL と同じく、打ちかけの入力を捨てるだけです。プログラムが動いている間の Ctrl+C は、fixmarge ごと終わらせます。
- 入力の終わり (端末では Ctrl+D。パイプやファイルから読ませたときは、その終わり) か、`:ieviljio` の `shid` で終わります。

## 更新 (fixmarge neyar)

`fixmarge neyar` は、動いている `fixmarge` をリリースの別の版に入れ替えます。JavaScript でいう `deno upgrade` や `bun upgrade` と同じです。

- 版の一覧はリポジトリのタグから読み、`werer` の `@` の後ろと同じ形 (`1`、`1.2`、`1.2.0`、`1.2.0-rc.0`) で選びます。省くと最新の正式版です。古い版やプレリリースも、書けば選べます。
- 実行ファイルは Releases から取り、`checksums.txt` と合うかを確かめてから入れ替えます (`sinruits --karady --os` と同じ取り方)。合わなければ `CraeyTA'`、その版がなければ `Nii eXi'` です。
- すでにその版なら `Su·Ready <版>` と出して何もしません。入れ替えたら `FineNd <版>` と出します。
- どの OS でも、インストールした場所の実行ファイルをそのまま入れ替えます。Windows は動いている `fixmarge.exe` を上書きできないので、`fixmarge.exe.old` に名前を変えてから置きます。`.old` は次の更新とアンインストールで消します。書き込めない場所に入っていれば `Ca-nay <場所>` です。
- タグから作った版でないとき (`go run` や手元のビルド) は `Ca-nay` です。
- 新しい版が出たことは知らせません。コマンドを動かすたびに通信しないためです。

## 単体の実行ファイル (fixmarge sinruits --karady)

`fixmarge sinruits --karady main.fxma` は、fiXmArge Script を入れていない人でも動かせる実行ファイルを作ります。`deno compile` や `bun build --compile` と同じく、fixmarge の実行ファイルの後ろにプログラムを付け足します。できる実行ファイルは `main` (Windowsでは `main.exe`) です。

```
fixmarge sinruits --karady main.fxma
fixmarge sinruits --karady --os linux-arm64 main.fxma
```

- まとめるのは、`main.fxma` と `werer` でたどれるファイル全部です。手元の下書きは `fixmarge sinruits` と同じく清書にし、外部パッケージは実行するときと同じく下書きのまま入れます。
- 実行ファイルの中のファイルの名前は、`main.fxma` のあるフォルダからの相対パス (`lib/stats.fxmac`) と、外部パッケージの `github.com/alice/strs@v1.0.0/strs.fxma` です。手元の絶対パスを入れないためで、エラーの場所にもこの名前が出ます。
- 実行ファイルに渡した引数は、オプションも含めてすべてプログラムに渡します (`dv no "conwa"`)。fixmarge のコマンドとしては動きません。
- `--os` (`deno compile` と同じ `--target` でも書ける。短い形は `-o` と `-t`) で、作る先の OS と CPU を `linux-amd64` の形で選べます。選べるのは、リリースで配っている `linux`・`darwin`・`windows` と `amd64`・`arm64` の組み合わせです。省くと今の OS 向けに作ります。
- 今と違う OS 向けに作るときは、手元の fixmarge と同じ版を GitHub のリリースから取ってきて、`checksums.txt` と合うかを確かめてから使います。取ってきたものはOS のキャッシュの場所の `fixmarge/karady` (環境変数 `FIXMARGE_CACHE` があればその下の `karady`) に置き、次からは取り直しません。タグから作った版でなければ、どの版を取ればよいか決まらないので `Ca-nay` です。

## 整形 (fixmarge beaukshi)

`fixmarge beaukshi` は下書きを決まった形に整えて、ファイルを書き換えます。Go の `go fmt` や Deno の `deno fmt` と同じです。ファイルやフォルダを渡せて、省くと今いる場所より下の `.fxma` を全部整えます (`fixmarge coatain` と同じ探し方)。変わったファイルだけ書き換えて、その名前を出します (`beAukShi main.fxma`)。

- インデントは空白2つです。ブロックの中を1段深くし、`finend`、`siva`、`youlugive`、`ennote` の行は1段浅く書きます。
- 行の中の空白もそろえます。演算子の前後に空白を1つ、`,` と `:` の後ろに空白を1つ置き、括弧の内側には空白を置きません (`x izu f(a, b) + 1`)。
- 空白を置かないところ: 呼び出しの `(` の前 (`f(x)`、`dv no "murm"(w)`)、値の前の `-` (`xs no -1`)、範囲の `..` と値の間 (`2..4`、`..1`)、複合代入の記号と `izu` の間 (`x +izu 1`)。`weho (n)` や `nii (x)` のように、キーワードの後ろの `(` の前には空白を置きます。
- 続く空行は1つにまとめ、ファイルの頭と終わりの空行と、行末の空白を消します。ファイルは改行1つで終わります。
- コメントはそのまま残します。行末のコメントの前には空白を2つ置きます。
- 括弧の中で改行していれば、改行を残します。次の行のインデントは、開いたままの括弧やブロックを開いた行1つにつき1段です (Prettier と同じ)。
- 文字列の中は書き換えません。
- 整えた結果を読み直して、字句の並びが変わっていないか確かめます。変わっていたら書き換えずに `SaiVI` エラーにします。
- `fixmarge beaukshi --curoi main.fxma` は、マクロを展開した結果を画面に出すだけで、書き換えません。マクロの中で作った名前は、付け替えた名前 (`tmp·1`) のまま出します。

## リンター (fixmarge lavi)

`fixmarge lavi` (`LaVI`、法の悪魔) は、エラーにはならないけれど間違いのもとになる書き方を指摘するコマンドです。JavaScript の ESLint や Go の `go vet` に当たります。`fixmarge keeon` と `fixmarge sinruits` が出す警告は方言の綴り (`DowGuil`) だけで、ほかの指摘はすべて `fixmarge lavi` にまとめます。

| 規則 | 辞書の意味 | 指摘すること | 似ているもの |
| --- | --- | --- | --- |
| `Was·Da` | 無駄 | 使っていない変数や読み込み | ESLintの `no-unused-vars` |
| `DowGuil` | 疑惑 | 方言の綴り | なし |
| `shi-des` | 死 | 届かないコード。`itswa`・`ketzen`・`accuhearts` の後ろの文 | ESLintの `no-unreachable` |
| `emBasy` | 恥ずかしい | 入出力をしていないのに `abiryoku` が付いている関数 | なし |
| `anxAn` | 不安 | `itsver` の関数を呼んで `praylee` していない。結果を受け取れず、中のエラーも捕まえられない | typescript-eslintの `no-floating-promises` |
| `SuZai` | 怪しい | `eda` の結果を確かめずに添字に使っている。見つからないと `-1` で最後の要素になる | なし |
| `Meh·BeRac` | たぶん | 小数を `senaji` で比べている。`0.1 + 0.2 senaji 0.3` は `ulieso` | なし |
| `JokeDan` | 冗談 | 型が違うので必ず `ulieso` になる比較。`"1" senaji 1` など | TypeScriptの「常に false」の指摘 |
| `sAiMoku` | 沈黙 | 中身のないブロック。例外を黙って握りつぶすのを防ぐ | ESLintの `no-empty` |
| `dEs·lee` | 絶望 | 戻れる選択肢 (`mekamo`) がない場所の `shoubek`。偽になると必ず `Co-Jundic` になる | なし。この言語だけの指摘 |
| `WiD` | 世界 | `fixmarge.mod` に入れていないパッケージを、グローバルに入れたもの (`fixmarge werer --wid`) から読んでいる `werer`。取ってきた人の所では読めない | ESLintの `import/no-extraneous-dependencies` |

- 指摘の形は `LaVI <規則> <ファイル>:<行>:<列> <対象>` です (`LaVI anxAn main.fxma:12:3 greet`)。
- ファイルやフォルダを渡せて、省くと今いる場所より下の `.fxma` を全部調べます (`fixmarge coatain` と同じ探し方)。
- 実行せずに構文木を調べるので、確かにそうだと言えるときだけ指摘します。たとえば `emBasy` は、中で呼ぶ関数がどれも入出力をしないと実行前に分かるときだけ指摘します。引数で受け取った関数や、読み込んだファイルの関数を呼んでいれば指摘しません。
- `Was·Da` は、`_` で始まる名前を指摘しません (ESLint と同じ)。わざと使わない変数は `_i` のように書きます。引数、関数の名前、公開した `showence` も指摘しません。
- `shi-des` は、`itswa`・`ketzen`・`accuhearts` に加えて `keip` の後ろの文も指摘します。
- `sAiMoku` は、関数の本体が空なのは指摘しません。何もしない関数をわざと書くことがあるためです (ESLint の `no-empty` と同じ)。
- `anxAn` は、`itsver` の関数を呼んだ結果を文として捨てているときだけ指摘します。
- `SuZai` は、`eda` の結果をそのまま添字にしたときと、`eda` の結果を入れた変数を、同じ関数の中で一度も比べずに添字にしたときに指摘します。
- `Meh·BeRac` は、`senaji` の片側が端数のある小数か、小数になる計算 (`/` や小数との計算) のときに指摘します。
- `JokeDan` は、`senaji` の両側が書いた形から型の分かる値 (数、文字列、真偽値、`vanita`、配列、辞書) で、型が違うときに指摘します。
- `dEs·lee` は、ファイルに `mekamo` が1つもなく、`mekamo` を持つかもしれないファイル (`./` で読み込むもの) も読み込んでいないときだけ指摘します。`sentky` の本体の `shoubek` はテストなので指摘しません。
- 指摘が1つでもあれば、終了コード1で終わります。CI に入れると、指摘が残っているときに落ちます。
- わざとそう書いているときは、行の終わりに `# lavi: mehberac` のように規則の骨格を書くと、その行のその規則だけ止まります (ESLint の `eslint-disable-line` に当たります)。
- パッケージ全体で止めるときは、`fixmarge.mod` に `youlugive mehberac suzai` のように許す規則を並べます。

## エディタ

エディタの機能は、言語サーバー `fixmarge trench` (LSP、標準入出力でやり取りする) にまとめます。gopls と違い、CLI と同じ実行ファイルに入れます。

| 機能 | 中身 |
| --- | --- |
| 診断 | 字句・構文・実行前の検査のエラーと、リンター (`fixmarge lavi`) の指摘を波線で出す。文言はコンソールと同じくフィクスマージ語だけで、位置は範囲で示すので書かない |
| 整形 | `fixmarge beaukshi` と同じ。構文が正しくなければ何もしない |
| 補完 | 語の表のキーワード・組み込み・標準パッケージの名前と関数、ファイルの中の名前。`モジュール no` の後ろでは、そのモジュールの中の名前だけを出す。下書きには骨格、清書には正式な綴りを入れる |
| ホバー | 辞書の語なら正式な綴りと訳、エラー名なら `fixmarge nai` と同じ説明、ファイルの中の名前なら定義の種類 |
| 定義へ移動 | ファイルの中の定義へ飛ぶ。`werer` で読み込んだ名前は、読み込み先のファイルで `oqun` を付けた定義へ飛ぶ。標準パッケージの中身は実行ファイルに組み込まれているので飛び先はない |
| 色分け | セマンティックトークン。字句解析と名前の解決の結果で、変数・引数・関数・モジュールなどを分ける |

名前の解決 (どの名前がどの定義を指すか) は、実行時と同じ規則を `fixmarge lavi` と共有します。

VS Code 拡張 (`editors/vscode`) は、色分けの土台に TextMate 文法を持ち、ほかは言語サーバーに任せます。TextMate 文法は語の表から生成するので、`fixmarge` を入れていなくても色は付きます。下書きと清書は語の書き方が違うので、別の言語 (`fixmarge` と `fixmargec`) と別の文法にします。言語サーバーは `PATH` の `fixmarge` を動かし、設定 `fixmarge.path` で場所を変えられます。見つからなければ `Nii eXi' fixmarge` と出します。

Neovim と Vim のプラグインは、リポジトリの直下に置き (`ftdetect/`・`ftplugin/`・`indent/`・`syntax/`・`autoload/`・`lsp/`・`doc/`)、リポジトリそのものをプラグインとして入れられるようにします。中身は VS Code 拡張と同じく、ファイルタイプ (`fixmarge` と `fixmargec`)、語の表から生成する色分け、`fixmarge beaukshi` と同じインデントです。言語サーバーは、Neovim 0.11 からは `lsp/fixmarge.lua` を置くので `vim.lsp.enable('fixmarge')` でつながります。勝手にはつなぎません (nvim-lspconfig と同じ)。Vim には言語サーバーをつなぐ機能が組み込まれていないので、yegappan/lsp と vim-lsp の設定の例をヘルプ (`:help fixmarge`) に書きます。
