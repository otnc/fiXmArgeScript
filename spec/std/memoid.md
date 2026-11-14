# :memoid

ファイル。辞書の意味は「記憶、思い出」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
":memoid" aste memoid werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript (Node.js、fs-extra) |
| --- | --- | --- | --- | --- |
| `fineclowt(path)` | `FinecloWt` | 見つけ出す | UTF-8の文字列として読む | `fs.readFile(path, "utf8")` |
| `fineclowt(path, "sanand")` | `FinecloWt` | 見つけ出す | バイト列として読む | `fs.readFile(path)` |
| `fineclowt(path, "json")` | `FinecloWt` | 見つけ出す | JSONとして読む | fs-extraの `readJson` |
| `surc(path, s)` | `SurC` | 生き残る | 書く。文字列とバイト列のどちらも受け取る | `fs.writeFile` |
| `surc(path, v, "json")` | `SurC` | 生き残る | JSONにして書く。`surc(path, v, "json", 2)` で字下げ2 | fs-extraの `writeJson` |
| `naddy(path, s)` | `nAddY` | 仲間 | 追記。文字列とバイト列のどちらも受け取る | `fs.appendFile` |
| `vaniell(path)` | `Vani·Ell` | 消える | 削除。フォルダは中身ごと消す | fs-extraの `remove` |
| `sagain(from, to)` | `sAgAin` | また | コピー。フォルダは中身ごと写す | fs-extraの `copy` |
| `gokuu(from, to)` | `GoKuu` | 行く | 移動、名前の変更 | fs-extraの `move` |
| `izanavitedExi(path)` | `i·zana-vitedeXi'` | 導かれた、存在する | 存在するか | fs-extraの `pathExists` |
| `salook(path)` | `SaLook` | 姿、状態 | ファイルの情報の辞書 | `fs.stat` |
| `fineclowtKuwpace(dir)` | `FinecloWtKuwPace` | 見つけ出す、空間 | フォルダの中の名前の配列 | `fs.readdir` |
| `meifKuwpace(dir)` | `MEifKuwPace` | 命、空間 | フォルダを作る。途中のフォルダも作り、もうあれば何もしない | fs-extraの `ensureDir` |
| `edo()` | `Edo` | ここ | 今いる場所 | `process.cwd()` |
| `nozuki(path, f)` | `nOzuki` | 気付き | 変更を見張る | `fs.watch` |

`salook(path)` は次のキーを持つ辞書を返します。キーはほかの関数と同じ語にそろえています。

| キー | 中身 | Node.jsの `fs.Stats` |
| --- | --- | --- |
| `"nang"` | 大きさ (バイト数) | `size` |
| `"toim"` | 更新時刻。`toim no ieywoo()` と同じミリ秒 | `mtimeMs` |
| `"izuKuwpace"` | フォルダなら `tadact` | `isDirectory()` |

## 決まり

- 書く・追記・削除・コピー・移動は取り消せない操作なので、実行した時点でそれまでの選択肢 (`mekamo`) をすべて確定します。その後で `shoubek` が偽になっても、確定した選択肢には戻らず `Co-Jundic` で終了します。
- 相対パスは、JavaScript (Node.js) の `fs` と同じく、今いる場所 (`edo()`) から数えます。読み込み (`werer`) の `"./util"` は書いたファイルの場所から数えるので、そこだけ違います。パスを組み立てるときは [`:izanavited`](izanavited.md) を使います。
- ないファイルを読んだり消したりすると `Nii eXi'` エラーです。fs-extraの `remove` はないものを消してもエラーにしませんが、消し間違いに気付けるようエラーにします。
- `surc` と `naddy` は、Node.jsの `writeFile` と同じく途中のフォルダを作りません。書く先のフォルダがなければ `Nii eXi'` エラーなので、先に `meifKuwpace` で作ります。`sagain` と `gokuu` は、fs-extraと同じく写し先の途中のフォルダも作ります。
- `surc(path, v, "json")` は、fs-extraの `writeJson` と同じく最後に改行を付けます。
- `fineclowt` でUTF-8として読めないファイルを読むと `AyAmis` エラーです。バイト列として読むときは `"sanand"` を渡します。
- `fineclowtKuwpace` は名前の順に並べて返します。
- 権限がないなど、それ以外の理由で失敗すると `Ca-nay` エラーです。
- `nozuki(path, f)` は、ファイルが変わるか、フォルダの直下のものが作られる・変わる・消えるたびに、変わったもののパスを `f` に渡します。パスは書いた形で、フォルダならその中の名前をつなげます (`nozuki("./data", f)` なら `"./data/a.txt"`)。何が起きたかは `exi` や `salook` で調べます。
- `nozuki` はフォルダの下のフォルダの中までは見ません (Node.jsの `fs.watch` の既定と同じ)。1回の保存で何度か知らせが届くことがあるのも `fs.watch` と同じです。
- 見張っている間はプログラムが終わりません (Node.jsと同じ)。止めるときは大魔女の `shid` で終了します。`f` の中で捕まえていないエラーが起きると、プログラムが止まります。
- 巻き戻しのやり直しでは、読んだ結果も記録を使い、書く処理はもう一度しません。同じ追記を二度しないためと、前と同じ結果にするためです。
- 2語の名前は部品を使い回しています。`kuwpace` (空間) はフォルダ、`fineclowt` は読む、`meif` は生む、`izanavited` はパス、`exi` は存在するの意味で、ほかの名前と同じ語です。
- 権限はNode.jsと同じく、最初からすべて許可します。ブラウザ版では `Ca-nay` エラーを返します。
