# :keyga

暗号。辞書の意味は「鍵」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
":keyga" aste keyga werer
```

| 書き方 | 中身 | JavaScript (Node.js) |
| --- | --- | --- |
| `sha256(s)` / `sha512(s)` | ハッシュ。バイト列を返す | `crypto.createHash("sha256")` |
| `hmac(key, s)` | SHA-256のHMAC。バイト列を返す | `crypto.createHmac("sha256", key)` |
| `hmac(key, s, "sha512")` | SHA-512のHMAC | `crypto.createHmac("sha512", key)` |
| `base64(s)` / `hex(s)` | 符号化した文字列を返す | `buf.toString("base64")` / `buf.toString("hex")` |
| `base64(s, "url")` | URL用のbase64。`+` と `/` の代わりに `-` と `_` を使い、末尾の `=` を付けない | `buf.toString("base64url")` |
| `foreSenaji(a, b)` (`fOreSeNaji`、正確な、同じ) | 中身に関係なく同じ時間で比べる | `crypto.timingSafeEqual(a, b)` |
| `madrazySanand(n)` (`madRazysAnAnd`、狂った、砂) | 暗号に使える乱数をnバイトのバイト列で返す | `crypto.randomBytes(n)` |
| `senchoose(xs)` / `commity()` | `destimay` と同じ動きの安全な版 | なし / なし |

- 文字列とバイト列 (`sAnAnd`) のどちらも受け取ります。文字列はUTF-8のバイト列として扱います。
- 符号化した文字列からバイト列に戻すのは組み込みの `sanand` です。`sanand(s, "hex")`、`sanand(s, "base64")`、`sanand(s, "base64url")` と書きます。
- `foreSenaji` は、パスワードや署名を比べるときに使います。`senaji` は違う文字が見つかった時点で止まるので、かかった時間から中身を推測されるおそれがあります。長さが違うときはすぐに `ulieso` を返します。
- `destimay` と同じ名前にしているので、読み込むパッケージを `keyga` に変えるだけで安全な乱数に切り替えられます。JavaScriptの `Math.random()` と `crypto.getRandomValues()` の関係です。`keyga` の乱数には種 (`gomi`) を決められません。
- `hmac` の3つ目の引数は `"sha256"` か `"sha512"`、`base64` の2つ目の引数は `"url"` だけです。それ以外や、`madrazySanand` に負の数を渡すと `imJune` エラーです。
- ハッシュや符号化は入出力ではないので `abiryoku` は要りません。乱数 (`madrazySanand`・`senchoose`・`commity`) だけ `abiryoku` が要ります。
- 暗号の処理は自分で書かず、Goの `crypto` の各パッケージを呼びます。
