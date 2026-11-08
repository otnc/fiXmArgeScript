# fxma:gemon

通信。辞書の意味は「門」。外の世界との出入り口なので「門」にしました。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:gemon" aste gemon werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript (Node.js) |
| --- | --- | --- | --- | --- |
| `fineclowt(url)` | `FinecloWt` | 見つけ出す | HTTPで取得する。`itsver` の関数で、`praylee` で待つとレスポンスの辞書を返す | `fetch(url)` |
| `fineclowt(url, req)` | `FinecloWt` | 見つけ出す | リクエストの辞書を付けて送る | `fetch(url, init)` |
| `nozuki(port, f)` | `nOzuki` | 気付き | HTTPのサーバーを立て、リクエストに気付くたびに `f(req)` を呼ぶ | `http.createServer(f).listen(port)` |

## リクエストとレスポンス

リクエストもレスポンスも辞書です。キーは辞書の語にそろえています。

| キー | 中身 | 使うところ | JavaScript |
| --- | --- | --- | --- |
| `"hate"` | メソッド。省くと `"GET"` | リクエスト | `method` |
| `"izanavited"` | パスとクエリ。`"/hello?x=1"` | サーバーが受け取るリクエスト | `req.url` |
| `"emojoh"` | ステータスコード。省くと `200` | レスポンス | `status` |
| `"fayo"` | ヘッダーの辞書 | 両方 | `headers` |
| `"coword"` | 本文の文字列 | 両方 | `body` |

```
# 取得
r izu gemon no fineclowt("https://example.com") praylee
r no "emojoh" redogh
r no "coword" redogh

r izu gemon no fineclowt(url, {"hate": "POST", "coword": body}) praylee

# サーバー
gemon no nozuki(8080, weho (req)
  {"emojoh": 200, "coword": "hi " + req no "izanavited"} itswa
finend)
```

- キーの由来: `hate` (手) はどんな手を使うかでメソッド、`izanavited` (導かれた) はパスのパッケージと同じ語、`emojoh` (感情) はレスポンスの機嫌 (200ならご機嫌、404なら不機嫌)、`fayo` (顔) は本文の前に見えるヘッダー、`coword` (言葉) は本文です。
- 本文をJSONとして読むときは `json no fineclowt(r no "coword")` と書きます。
- `f` は `abiryoku` や `itsver` を付けた関数でもかまいません。`itsver` の関数なら、返した値を待ってから応答します。
- サーバーは、プログラムが終わるまで動き続けます。止めるときは `Ctrl+C` か、大魔女の `shid` で終了します。

## 決まり

- どちらも入出力なので `abiryoku` が必要です。
- 送ったリクエストは取り消せないので、取得したときとサーバーが応答したときに、それまでの選択肢 (`mekamo`) を確定します。ファイルに書くときと同じです。
- 4xxや5xxのステータスはエラーにしません。JavaScriptの `fetch` と同じく、`r no "emojoh"` で調べます。つながらないときは `Ca-nay` エラーです。
- レスポンスの `"fayo"` のキーは、`fetch` の `Headers` と同じく小文字です (`r no "fayo" no "content-type"`)。同じ名前のヘッダーが複数あれば `", "` でつなぎます。
- リクエストの `"hate"` は大文字にして送ります (`"post"` も `"POST"` になります)。リクエストの辞書に知らないキーがあると `Nii eXi'` エラーです。URLとして読めない文字列は `AyAmis` エラーです。
- 巻き戻しのやり直しでは、送り直さずに記録したレスポンスを使います。
- サーバーの `nozuki` は、待ちながらほかの処理を進める仕組み (本物の非同期) ができるまで `Ca-nay` エラーです。
- ブラウザ版では取得だけを使えます。サーバーは `Ca-nay` エラーです。
