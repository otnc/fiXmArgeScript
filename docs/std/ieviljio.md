# fxma:ieviljio

プロセス。辞書の意味は「魔女」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

この言語ではプロセスを「魔女」と呼びます。実行中のプログラム自身は「大魔女」で、`diviljio()` で取り出します。大魔女は魔女を生み、魔女を殺す力をただ一人持ちます。Node.jsの `process` に当たります。

```
"fxma:ieviljio" aste ieviljio werer
dv izu ieviljio no diviljio()
```

| 書き方 | 辞書の綴り | 辞書の意味 | 中身 | Node.js |
| --- | --- | --- | --- | --- |
| `dv no "meif"(cmd, args, opts)` | `MEif` | 命 | 魔女 (プロセス) を生む。`opts` は省ける | `child_process.spawn` |
| `w praylee` | `pray·lee` | 祈る | 魔女が終わるのを待ち、出力と終了コードの辞書を得る | `child.on("exit")` |
| `dv no "murm"(w)` | `muRm` | 殺人 | 魔女を殺す (プロセスを止める) | `child.kill()` |
| `dv no "alte"` | `Alte` | 全て | 生んだ魔女の一覧 | なし |
| `dv no "conwa"` | `ConWa` | 会話 | コマンドライン引数 (配列) | `process.argv` |
| `dv no "lowhell"` | `LowhEll` | 牢屋 | 環境変数 (辞書)。`dv no "lowhell" no "HOME"` で読む | `process.env` |
| `dv no "nozuki"("SIGINT", f)` | `nOzuki` | 気付き | シグナル (Ctrl+Cなど) に気付いて関数を呼ぶ | `process.on("SIGINT", f)` |
| `dv no "shid"(番号, 道連れ)` | `shid` | 死んだ | 終了コードを付けて終了する。`道連れ` は省ける | `process.exit()` |
| `dv no "ieviljio"` | `ievil JIO` | 魔女 | プロセス番号 (魔女の番号) | `process.pid` |
| `dv no "wid"` | `WiD` | 世界 | 動いている環境の名前。値は下の表 | `process.platform` |
| `dv no "neyar"` | `nEYar` | 年 | fxmaの版 (`fxma neyar` と同じ語) | `process.version` |

## 動いている環境

`dv no "wid"` の値は、Goの `GOOS` をそのまま使わず、fxmaで決めた一覧から返します。WebAssembly版はブラウザでもNode.jsでも `GOOS=js` になり、区別できないためです。

| 値 | 環境 | 決め方 |
| --- | --- | --- |
| `"windows"`、`"darwin"`、`"linux"` など | CLI、C・Pythonなどへの組み込み | Goの `runtime.GOOS` と同じ |
| `"browser"` | ブラウザ (プレイグラウンドなど) | 実行中に `globalThis.document` があるかで調べる |
| `"node"` / `"deno"` / `"bun"` | npm版をそれぞれで動かしたとき | 実行中に `process.versions.node`、`Deno`、`Bun` があるかで調べる |
| `"wasi"` | WASIランタイム | `GOOS=wasip1` のビルド。WASIの版が変わっても同じ値にする |
| `"baremetal"` | TinyGoでのマイコン | TinyGoのマイコン向けのビルド |

## 魔女を生むときの指定

`meif` の3つ目の引数は、次のキーを持つ辞書です。どれも省けます。キーは、同じ意味のほかの語にそろえています。

| キー | 中身 | Node.jsの `spawn` |
| --- | --- | --- |
| `"zanfession"` | 魔女に渡す入力の文字列 (`zanfession` と同じ語) | `child.stdin.write` |
| `"edo"` | 魔女を動かす場所 (`memoid no edo` と同じ語) | `cwd` |
| `"lowhell"` | 環境変数の辞書 (`dv no "lowhell"` と同じ語)。省くと大魔女と同じ | `env` |
| `"nozuki"` | 出力に気付くたびに呼ぶ関数 | `child.stdout.on("data", f)` |

`"nozuki"` に渡した関数は、出力が届くたびに `{"redogh": 届いた文字列}` か `{"canosiri": 届いた文字列}` の辞書を受け取ります。長く動くコマンドの進み具合を出すときに使います。途中で受け取った出力も、`w praylee` の結果にまとめて入ります。

```
w izu dv no "meif"("sort", [], {"zanfession": "b
a
", "edo": "./data"})
r izu w praylee
```

## 魔女を待った結果

`w praylee` は次のキーを持つ辞書を返します。キーは、出力の語と終了の語に合わせています。

| キー | 中身 | Node.jsの `execFile` |
| --- | --- | --- |
| `"redogh"` | 標準出力の文字列 | `stdout` |
| `"canosiri"` | 標準エラー出力の文字列 | `stderr` |
| `"shid"` | 終了コード | `code` |

```
w izu dv no "meif"("git", ["status"])
r izu w praylee
r no "redogh" redogh
mof nii r no "shid" senaji 0
  r no "canosiri" canosiri
finend
```

## 魔女を殺す力

`murm` は大魔女だけが持つ関数で、パッケージから直接は呼べません。他の関数に魔女を殺させたいときは、大魔女から取り出した `murm` を値として渡します。

```
abiryoku weho cleanup(witches, murm)
  witches aste w azagara
    murm(w)
  finend
finend

cleanup(dv no "alte", dv no "murm")
```

力を渡されていない関数やパッケージは、魔女を殺せません。読み込んだ外部パッケージが勝手にプロセスを止めることも防げます (権限を値として渡す、ケーパビリティという考え方です)。

## そのほかの決まり

- `diviljio` は辞書の `Divil JIO` (大魔女) です。清書では `Divil JIO` と書きます。
- `dv no "conwa"` は、実行するファイルより後ろの引数だけを並べた配列です (Denoの `Deno.args` と同じ)。`fxma keeon main.fxma -x a` なら `["-x", "a"]` で、`-x` のようなオプションもfxmaは読まずにそのまま渡します。Node.jsの `process.argv` と違い、実行ファイルとスクリプトのパスは入れません。
- `meif` の2つ目の引数 (引数の配列) は省けません。引数がなければ `[]` を渡します。
- 魔女が0以外の終了コードで終わっても、エラーではなく `"shid"` に入ります (Node.jsの `spawn` と同じ)。コマンドが見つからなければ `Nii eXi'`、ほかの理由で生めなければ `Ca-nay` エラーです。
- `meif` の3つ目の引数に知らないキーがあると `Nii eXi'` エラーです。
- 終わった魔女を `murm` で殺しても何も起きません。
- 巻き戻しのやり直しでは、魔女を生み直さず、待った結果も記録を使います。
- 出力やシグナルに気付く `nozuki` は、待ちながらほかの処理を進める仕組み (本物の非同期) ができるまで `Ca-nay` エラーです。

- `meif`・`murm`・`shid` は取り消せない操作なので、実行した時点で選択肢を確定します。
- `dv no "shid"` で終了するとき、まだ生きている魔女も一緒に止めます。2つ目の引数に `ulieso` を渡すと、魔女を残して終了します (`dv no "shid"(0, ulieso)`)。
- シグナルの名前は、Node.jsの `process.on` と同じく `"SIGINT"`、`"SIGTERM"` のように書きます。
- ブラウザ版では、`meif`・`murm`・`nozuki`・`shid` が `Ca-nay` エラーです。`diviljio` と情報の取り出しは使えます。
- 前の版にあった `fxma:wid` は廃止し、引数・環境変数・終了は大魔女に移しました。
