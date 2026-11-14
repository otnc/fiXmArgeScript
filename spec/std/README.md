# 標準パッケージ

fiXmArge Script (fxma) の標準パッケージの仕様です。言語そのものの仕様は [../README.md](../README.md) から章ごとに読めます。

標準パッケージは `":名前"` で読み込みます。Node.jsの `node:fs` から `node` を省いた形で、`:` で始まる読み込み先は標準パッケージ専用です。読み込み方は普通のパッケージと同じです。パッケージ名と関数名は辞書の語で、JSONやSHA-256のような固有名詞だけはそのまま使います。

## 一覧

| パッケージ | 中身 |
| --- | --- |
| [`:meme`](meme.md) | 数学 |
| [`:coword`](coword.md) | 文字列 |
| [`:varizama`](varizama.md) | 配列 |
| [`:toim`](toim.md) | 時間 |
| [`:destimay`](destimay.md) | 乱数 |
| [`:keyga`](keyga.md) | 暗号 |
| [`:json`](json.md) | JSON |
| [`:ieviljio`](ieviljio.md) | プロセス |
| [`:gemon`](gemon.md) | 通信 |
| [`:memoid`](memoid.md) | ファイル |
| [`:izanavited`](izanavited.md) | パス |

## 作り方

- Goの標準ライブラリを呼ぶもの (文字列、時刻、乱数、暗号、JSON) はGoで書きます。言語自体で書けるもの (配列の合計など) は `.fxma` で書いて本体に埋め込みます。
- 入出力をする関数 (時刻、乱数、引数) は `abiryoku` 付きとして扱います。呼ぶ関数にも `abiryoku` が必要です。
- 時刻と乱数の結果は、入力と同じく記録しておき、巻き戻しでやり直すときは記録を使います。
- 実行する環境で使えない関数 (ブラウザでの引数など) は `Ca-nay` エラーを返します。
- パッケージが違えば、同じ関数名を使えます。同じ名前を2つのパッケージから個別に読み込むときだけ、`aste` で別名を付けます。
- 暗号の処理は自分で書かず、Goの `crypto` の各パッケージを呼びます。

## 組み込み (読み込み不要)

| 書き方 | 清書 | 辞書の意味 | 中身 | Python |
| --- | --- | --- | --- | --- |
| `nang(xs)` | `nAng` | 長い | 配列・文字列・辞書の長さ | `len(xs)` |
| `sanand(s)` | `sAnAnd` | 砂 | 文字列をUTF-8のバイト列にする。`sanand(s, "hex")`、`sanand(s, "base64")` で16進・base64の文字列から戻す | `s.encode()`、`bytes.fromhex(s)`、`base64.b64decode(s)` |

## 標準パッケージのために足した言語の機能

- **辞書型**: `{"lang": "fxma"}` で作り、`d no "lang"` で取り出し、`d no "lang" izu "x"` で書き換えます。キーがあるかは `"lang" exi d` で調べます。キーは入れた順に並びます。JavaScriptの `Map` やPythonの `dict` と同じです。`d aste k azagara` で回すと、`k` にはキーが入った順に入ります。JavaScriptの `for (const k in obj)` と同じで、値は `d no k` で取り出します。JSONのオブジェクトを受け取るために必要です。
- **バイト列型**: 型の名前は `sAnAnd` (砂) で、1バイトを砂の1粒に見立てます。PythonのbytesやGoの `[]byte` に当たる型です。`madrazySanand`、`sha256`、`sha512`、`hmac` はバイト列を返し、`base64` と `hex` はバイト列か文字列を受け取って文字列を返します。文字列はUTF-8のバイト列として扱います。バイト列を直接書く書き方 (リテラル) はなく、組み込みの `sanand` で文字列から作ります。型の名前で変換するのは、Goの `[]byte(s)` と同じ発想です。文字列に戻すのは `coword no salook` です。
- **関数を値として渡す**: `weho` で定義した関数は、名前で他の関数に渡せます (`varizama no tarip(xs, double)`)。名前のない関数は、`weho` の名前を省いて式の中に書きます。JavaScriptの `function` 式と同じです。`varizama no tarip(xs, weho (n) n * 2 itswa finend)` のように1行でも書けます。`abiryoku` も付けられます。
