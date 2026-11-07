# fxma:meme

数学。辞書の意味は「意味、価値」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:meme" aste meme werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `racdam(x)` | `RacDA-m` | 落ちる | 切り捨て | `Math.floor` |
| `hir(x)` | `Hir` | 空を飛ぶもの | 切り上げ | `Math.ceil` |
| `litto(x)` | `Litto` | ちょうど | 四捨五入。ちょうど真ん中は0から遠い方に丸める | `Math.round` に近い。`-2.5` の扱いが違う |
| `pyusense(x)` | `PyusenS'e` | 純粋な | 絶対値 | `Math.abs` |
| `mource(x)` | `mOurce` | 元、根源 | 平方根 | `Math.sqrt` |
| `tsuyong(a, b)` / `yoweak(a, b)` | `tSuYoNg` / `YowEAk` | 強い / 弱い | 最大 / 最小（配列も受け取る） | `Math.max` / `Math.min` |
| `zelln(xs)` | `Zelln` | 全員 | 合計 | なし |

- `racdam`、`hir`、`litto` は小数を返します。GoやJavaScriptと同じです。`litto(2.5)` は `3.0`、`litto(-2.5)` は `-3.0` です。
- `tsuyong` と `yoweak` に空の配列を渡すと `Nii eXi'` エラーです。`zelln` に空の配列を渡すと `0` を返します。

## 数学の名前をそのまま使う関数

三角関数などは辞書に合う語がないので、JSONやSHA-256と同じく固有名詞として扱い、数学の名前をそのまま使います。

| 書き方 | 中身 | JavaScript |
| --- | --- | --- |
| `pi` / `e` | 円周率 / 自然対数の底（定数） | `Math.PI` / `Math.E` |
| `sin(x)` / `cos(x)` / `tan(x)` | 三角関数。角度はラジアン | `Math.sin` など |
| `asin(x)` / `acos(x)` / `atan(x)` | 逆三角関数 | `Math.asin` など |
| `atan2(y, x)` | 点 (x, y) の角度 | `Math.atan2` |
| `exp(x)` | eのx乗 | `Math.exp` |
| `log(x)` / `log2(x)` / `log10(x)` | 自然対数 / 2を底 / 10を底 | `Math.log` など |

- 定数は `showence` と同じく書き換えられません。`meme no pi` で読みます。
- 0以下の数の対数、負の数の `mource`、範囲外の `asin` は `imJune` エラーです。JavaScriptなら `NaN` になるところです。
