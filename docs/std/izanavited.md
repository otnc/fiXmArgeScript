# fxma:izanavited

パス。辞書の意味は「導かれた」。パスに導かれてファイルにたどり着く、という見立てです。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:izanavited" aste izanavited werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript（Node.js） |
| --- | --- | --- | --- | --- |
| `recanrey(a, b, ...)` | `rEcanRey` | 関係 | パスをつなぐ | `path.join` |
| `defutai(p)` | `deFuTai` | 絶対に | 今いる場所から数えた絶対パスにする | `path.resolve` |
| `kuwpace(p)` | `KuwPace` | 空間 | 入っているフォルダのパス | `path.dirname` |
| `ado(p)` | `Ado` | これ | 最後の名前 | `path.basename` |
| `chidd(p)` | `ChiDd` | 血 | 拡張子。`.` を含む | `path.extname` |
| `eatahguy(from, to)` | `eATah·Guy` | お互い | `from` から見た `to` の相対パス | `path.relative` |
| `saymal(p)` | `Say mAl` | 正常 | `.` や `..` を整理する | `path.normalize` |
| `izuDefutai(p)` | `Iz'udeFuTai` | 〜である、絶対に | 絶対パスかどうか | `path.isAbsolute` |

```
izanavited no recanrey("a", "b", "c.txt")  # "a/b/c.txt"
izanavited no ado("a/b/c.txt")             # "c.txt"
izanavited no chidd("a/b/c.txt")           # ".txt"
izanavited no eatahguy("a/b", "a/c")       # "../c"
```

- 区切り文字は、受け取るときは `/` も `\` も読み、返すときはいつも `/` にします。WindowsでもGoやNode.jsのファイル操作は `/` を受け付けるので、どのOSでも同じ結果になります。
- ファイルには触らず、文字列だけを扱います。入出力ではないので `abiryoku` は要りません。ただし `defutai` は今いる場所を読むので `abiryoku` が必要です。
- ブラウザ版でも使えます。ブラウザ版の今いる場所は `/` です。
- `chidd` は、例外の種類を取り出す `e no "chidd"` と同じ語です。拡張子はファイルの種類（血筋）を表すからです。
- 名前の多くはほかのパッケージと同じ語です。`recanrey` は `coword` の結合、`saymal` は `fxma ado saymal`（mod tidy）と同じ「整える」、`kuwpace` は `memoid no meifKuwpace` などの「フォルダ」と同じ語です。
