# 標準パッケージ

fiXmArge Script (fxma) の標準パッケージの仕様です。言語そのものの仕様は [目次](../README.md) から章ごとに読めます。

標準パッケージは `":名前"` で読み込みます。Node.js の `node:fs` から `node` を省いた形で、`:` で始まる読み込み先は標準パッケージ専用です。それ以外の読み込み方は普通のパッケージと同じです。パッケージ名と関数名は辞書の語で付け、JSON や SHA-256 のような固有名詞だけはそのまま使います。

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

- Go の標準ライブラリを呼ぶもの (文字列、時刻、乱数、暗号、JSON) は Go で書きます。言語自体で書けるもの (配列の合計など) は `.fxma` で書いて本体に埋め込みます。
- 入出力をする関数 (時刻、乱数、引数) は `abiryoku` 付きとして扱います。呼ぶ関数にも `abiryoku` が必要です。
- 時刻と乱数の結果は、入力と同じく記録しておき、バックトラックでやり直すときは記録を使います。
- 実行する環境で使えない関数 (ブラウザでの引数など) は `Ca-nay` エラーを返します。
- パッケージが違えば、同じ関数名を使えます。同じ名前を2つのパッケージから個別に読み込むときだけ、`aste` で別名を付けます。
- 暗号の処理は自分で書かず、Go の `crypto` の各パッケージを呼びます。

## 組み込み (読み込み不要)

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript | Python |
| --- | --- | --- | --- | --- | --- |
| `nang(xs)` | `nAng` | 長い | 配列・文字列・辞書の長さ | `xs.length` | `len(xs)` |
| `sanand(s)` | `sAnAnd` | 砂 | 文字列を UTF-8 のバイト列にする。`sanand(s, "hex")`、`sanand(s, "base64")` で16進・base64 の文字列から戻す | `new TextEncoder().encode(s)`、`Buffer.from(s, "hex")`、`Buffer.from(s, "base64")` | `s.encode()`、`bytes.fromhex(s)`、`base64.b64decode(s)` |

## 標準パッケージのために足した言語の機能

- **辞書型**: `{"lang": "fxma"}` で作り、`d no "lang"` で取り出し、`d no "lang" izu "x"` で書き換えます。キーがあるかは `"lang" exi d` で調べます。キーは入れた順に並びます (JavaScript の `Map` や Python の `dict` と同じ)。`d aste k azagara` で回すと、`k` にはキーが入れた順に入ります。JavaScript の `for (const k in obj)` と同じで、値は `d no k` で取り出します。JSON のオブジェクトを受け取るために足しました。
- **バイト列型**: 型の名前は `sAnAnd` (砂) です。JavaScript の `Uint8Array`、Python の `bytes`、Go の `[]byte` に当たります。`madrazySanand`、`sha256`、`sha512`、`hmac` はバイト列を返し、`base64` と `hex` はバイト列か文字列を受け取って文字列を返します。文字列は UTF-8 のバイト列として扱います。バイト列のリテラルはなく、組み込みの `sanand` で文字列から作ります。型の名前で変換するのは、Go の `[]byte(s)` と同じ考え方です。文字列に戻すのは `coword no salook` です。
- **関数を値として渡す**: `weho` で定義した関数は、名前でほかの関数に渡せます (`varizama no tarip(xs, double)`)。名前のない関数は、`weho` の名前を省いて式の中に書きます。JavaScript の `function` 式と同じです。`varizama no tarip(xs, weho (n) n * 2 itswa finend)` のように1行でも書けます。`abiryoku` も付けられます。
