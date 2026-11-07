# fxma:json

JSON。名前は固有名詞なのでそのまま使います。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:json" aste json werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `fineclowt(s)` | `FinecloWt` | 見つけ出す | JSONを読む | `JSON.parse(s)` |
| `salook(v)` | `SaLook` | 姿、形 | 1行のJSONにする | `JSON.stringify(v)` |
| `salook(v, n)` | `SaLook` | 姿、形 | 字下げ `n` で整形したJSONにする | `JSON.stringify(v, null, n)` |

| JSON | fxma |
| --- | --- |
| 数 | 小数。JavaScriptの `JSON.parse` と同じで、`1` も `1.0` になる |
| 文字列 | 文字列 |
| `true` / `false` | `tadact` / `ulieso` |
| `null` | `vanita` |
| 配列 | 配列 |
| オブジェクト | 辞書。キーはJSONに書かれた順に並ぶ |

- 数はいつも小数なので、2の53乗を超える整数は、JavaScriptと同じく正確に読めません。
- 正しくないJSONを読むと `AyAmis` エラーです。JavaScriptの `JSON.parse` が `SyntaxError` を出すのと同じです。
- `salook` は、整数も小数も数として書き出します。関数やバイト列のようにJSONにできない値は `imJune` エラーです。
- 字下げはJavaScriptと同じく10までで、1より小さければ1行にします。空の配列と辞書は `[]` と `{}` です。
- 文字列でエスケープするのは、JavaScriptと同じく `"` と `` と制御文字だけです。日本語はそのまま書き出します。
- 値の後ろに別の値やゴミが続くJSON（`1 2`）も `AyAmis` エラーです。
- JSONでは文字列そのものもJSONにできるため、`coword` のように1つの名前で両方向にはしません。
