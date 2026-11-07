# fxma:coword

文字列。辞書の意味は「言葉」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:coword" aste coword werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `salook(x)` | `SaLook` | 姿、形 | 数値と文字列を相互に変える。バイト列はUTF-8の文字列にする | `String(x)` / `Number(s)` |
| `confsoy(s, sep)` | `conf-Soy` | 争い | 分割 | `s.split(sep)` |
| `recanrey(xs, sep)` | `rEcanRey` | 関係 | 結合 | `xs.join(sep)` |
| `lainei(s, a, b)` | `Lai Nei` | 〜のように | 最初の1つだけ置換 | `s.replace(a, b)` |
| `laineiAlte(s, a, b)` | `Lai NeiAlte` | 〜のように、全て | 全部置換 | `s.replaceAll(a, b)` |
| `shojika(s)` | `ShoJikA` | 短い | 前後の空白削除 | `s.trim()` |
| `gdie(s)` / `lighluw(s)` | `gDie` / `Ligh-Luw` | 偉大な / 軽い | 大文字 / 小文字 | `s.toUpperCase()` / `s.toLowerCase()` |
| `eda(s, x)` | `Eda` | そこ | `x` が最初に現れる位置。なければ `-1` | `s.indexOf(x)` |
| `kast(s, x)` / `gun(s, x)` | `Kast` / `gu-n` | 過去 / 今 | `x` で始まるか / 終わるか | `s.startsWith(x)` / `s.endsWith(x)` |
| `sodaKast(s, n, c)` / `sodaGun(s, n, c)` | `SodaKast` / `Sodagu-n` | その程度まで、過去 / 今 | 長さが `n` になるまで、前 / 後ろを `c` で埋める。`c` を省くと空白 | `s.padStart(n, c)` / `s.padEnd(n, c)` |
| `kolsu(c)` / `batteKolsu(n)` | `KolSu` / `BatteKolSu` | 音 / 〜から、音 | 文字を番号に / 番号から文字を作る | `c.codePointAt(0)` / `String.fromCodePoint(n)` |

`salook` は、数値を渡すと文字列を、文字列を渡すと数値を返します。文字列から作る数は、JavaScriptの `Number("3")` と同じく、`"3"` でも `3.0` のように小数になります。数として読めない文字列を渡すと `AyAmis` エラーです。

- `kolsu` に空の文字列を渡すと `Nii eXi'` エラーです。
- 位置は配列の添字と同じく、文字の番号（コードポイント）で数えます。`eda` は見つからなければJavaScriptの `indexOf` と同じく `-1` を返します。`xs no -1` は最後の要素なので、そのまま添字に使わないよう気を付けます。
- `eda` の名前は、`Eda`（そこ）で「どこにあるか」を表します。`kast` と `gun` は、文字列を左から読むと頭が過去でお尻が今になる、という見立てです。
- `kolsu` は、文字の音を聞いて番号にする、という見立てです。逆向きは、JavaScriptの `fromCodePoint` と同じく「〜から」の `batte` を前に付けた `batteKolsu` です。
- 2語の名前は部品を使い回しています。`sodaKast` と `sodaGun` の `kast`（頭）と `gun`（お尻）は、`kast` と `gun`（startsWith / endsWith）と同じ語です。`laineiAlte` の `alte`（全て）は、`dv no "alte"` などの「全部」と同じ語です。
