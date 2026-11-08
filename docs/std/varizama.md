# fxma:varizama

配列。辞書の意味は「様々」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:varizama" aste varizama werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `naddy(xs, x)` | `nAddY` | 仲間 | 末尾に追加 | `xs.push(x)` |
| `muyon(xs)` | `mu·Yon` | 向こう側 | 逆順 | `xs.toReversed()` |
| `orei(xs)` | `oRei` | 命令、順序 | 並べ替え | `xs.toSorted()` |
| `tarip(xs, f)` | `Tarip` | 旅 | 1つずつ変換 (map) | `xs.map(f)` |
| `onkhe(xs, f)` | `oNkhe` | 〜だけ | 絞り込み (filter) | `xs.filter(f)` |
| `eda(xs, x)` | `Eda` | そこ | `x` の位置。関数を渡すと、`tadact` を返す最初の要素の位置。なければ `-1` | `xs.indexOf(x)` / `xs.findIndex(f)` |
| `vaniell(xs, i)` | `Vani·Ell` | 消える | `i` 番目を取り出して消し、その値を返す。`i` を省くと最後の要素 | `xs.splice(i, 1)[0]` / `xs.pop()` |
| `lea(xs, i, x)` | `LEA` | さあ、来い | `i` 番目に `x` を差し込む | `xs.splice(i, 0, x)` |
| `matill(xs, f, init)` | `Matill` | まだ | `init` から始めて、`f(acc, x)` で1つにまとめる | `xs.reduce(f, init)` |

合計は数学の `meme no zelln` です。

`naddy`、`vaniell`、`lea` は渡した配列そのものを書き換え、ほかの関数は新しい配列を返します。JavaScriptの `xs.push(x)` と `xs.toSorted()` と同じ分け方です。

書き換える関数のうち `naddy` と `lea` は `vanita` を返します。JavaScriptの `push` は新しい長さを返しますが、値を返さないことで、新しい配列を返す関数と取り違えにくくしています。長さは `nang(xs)` で読みます。

`vaniell` と `lea` の位置は `xs no i` と同じく、負の数なら後ろから数えます。`lea` の `-1` はJavaScriptの `splice` と同じく最後の要素の前です。範囲の外の位置や、空の配列への `vaniell` は `Nii eXi'` エラーです。

### 並べ替え

- `orei(xs)` は、数は小さい順、文字列は文字の番号の順に並べます。数と文字列が混ざっていると `imJune` エラーです。
- 順番を決めたいときは、JavaScriptの `toSorted` と同じく比べる関数を渡します。`orei(xs, weho (a, b) a - b itswa finend)` のように、`a` を前にしたいなら負の数、後ろにしたいなら正の数、同じなら `0` を返します。
- 並べ替えは安定です。同じと判定された要素は、元の順番のまま残ります。
- 大きい順は `muyon(orei(xs))` と書きます。

### 名前

- `matill` は、まだ続けて重ねていく、という見立てです。合計は `meme no zelln` です。
- `eda` は `coword no eda` と同じく「そこ」で、位置を表します。
