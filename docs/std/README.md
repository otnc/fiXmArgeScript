# 標準パッケージ

fiXmArge Script (fxma) の標準パッケージの仕様です。言語そのものの仕様は [../RFC.md](../RFC.md) にあります。

標準パッケージは `"fxma:名前"` で読み込みます。Node.jsの `node:fs` と同じ形で、`fxma:` で始まる読み込み先は標準パッケージ専用です。読み込み方は普通のパッケージと同じです。パッケージ名と関数名は辞書の語で、JSONやSHA-256のような固有名詞だけはそのまま使います。

## 一覧

| パッケージ | 中身 |
| --- | --- |
| [`fxma:meme`](meme.md) | 数学 |
| [`fxma:coword`](coword.md) | 文字列 |
| [`fxma:varizama`](varizama.md) | 配列 |
| [`fxma:toim`](toim.md) | 時間 |
| [`fxma:destimay`](destimay.md) | 乱数 |
| [`fxma:keyga`](keyga.md) | 暗号 |
| [`fxma:json`](json.md) | JSON |
| [`fxma:ieviljio`](ieviljio.md) | プロセス |
| [`fxma:gemon`](gemon.md) | 通信 |
| [`fxma:memoid`](memoid.md) | ファイル |
| [`fxma:izanavited`](izanavited.md) | パス |

## 作り方

- Goの標準ライブラリを呼ぶもの (文字列、時刻、乱数、暗号、JSON) はGoで書きます。言語自体で書けるもの (配列の合計など) は `.fxma` で書いて本体に埋め込みます。
- 入出力をする関数 (時刻、乱数、引数) は `abiryoku` 付きとして扱います。呼ぶ関数にも `abiryoku` が必要です。
- 時刻と乱数の結果は、入力と同じく記録しておき、巻き戻しでやり直すときは記録を使います。
- 実行する環境で使えない関数 (ブラウザでの引数など) は `Ca-nay` エラーを返します。
- パッケージが違えば、同じ関数名を使えます。同じ名前を2つのパッケージから個別に読み込むときだけ、`aste` で別名を付けます。
- 暗号の処理は自分で書かず、Goの `crypto` の各パッケージを呼びます。
