# キーワード一覧

キーワードはすべて辞書にある語です。表の「意味」は辞書の訳で、「JavaScript」と「Python」の列は同じことをそれぞれの言語で書いた形です。「清書」はビルドした後の綴りです ([表記](notation.md))。

## 文と制御

| 書き方 | 清書 | 意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| `x izu 1` | `Iz'u` | is | `let x = 1` | `x = 1` |
| `showence x izu 1` | `showence` | 証拠 | `const x = 1` | 再代入できない定数 |
| `x redogh` | `ReDogh'` | 伝える | `console.log(x)` | `print(x)` |
| `x canosiri` | `CaNoSiri` | 罵り、叫び | `console.error(x)` | `print(x, file=sys.stderr)` |
| `x izu "名前は?" zanfession` | `zanfession` | 懺悔 | `x = prompt("名前は?")` | `x = input("名前は?")` |
| `mof 条件` | `Mof` | もし | `if` | `if` |
| `siva` | `SiVa` | しかし | `else` | `else` |
| `siva mof 条件` | `SiVa Mof` | — | `else if` | `elif` |
| `条件 azagara` | `aZaGara` | 〜しながら | `while` | `while` |
| `xs aste x azagara` | `aSte` | 〜として | `for (const x of xs)` | `for x in xs` |
| `keip` | `Kei-P` | 継続 | `continue` | `continue` |
| `ketzen` | `KetZen` | 結末 | `break` | `break` |
| `finend` | `FineNd` | 終わり | `}` | ブロックの終わり |
| `weho f(n)` | `We-Ho` | 方法 | `function f(n)` | `def f(n)` |
| `n itswa` | `itsWA` | 以上だ | `return n` | `return n` |
| `abiryoku weho f()` | `AbiRyoku` | 能力 | 対応なし | 対応なし (副作用がある関数の印) |
| `#` | `#` | — (記号) | `//` | `#` (コメント) |

Ruby の `unless` や `until` に当たる専用の語はなく、`nii` (否定) で書きます。`nii` は比較より結び付きが弱いので、括弧は要りません。

| 書き方 | JavaScript | Python |
| --- | --- | --- |
| `mof nii x senaji 1` | `if (!(x === 1))` | `if not x == 1` |
| `nii x > 10 azagara` | `while (!(x > 10))` | `while not x > 10` |

## バックトラック・テスト・例外の構文

使い方は [言語の特徴](features.md) で説明します。

| 書き方 | 清書 | 意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| `x izu 1 owk 2 mekamo` | `meKaMo` | 〜かも | 対応なし (候補を並べる) | 対応なし |
| `条件 shoubek` | `ShouBeK` | 〜するべき | `console.assert()` に近い | `assert` (ただし失敗すると巻き戻る) |
| `sentky "名前"` | `sent-ky` | 裁き | `test()` (Jest など) | テスト関数 |
| `persemon "名前"` | `Per·SemOn` | 見世物 | `bench()` (Vitest) | ベンチマーク関数 |
| `"理由" accuhearts` | `AccUhEarts` | 告発 | `throw` | `raise` |
| `teitance` | `Tei-Tance` | 抵抗 | `try` | `try` |
| `youlugive e` | `yOuLugiVe` | 許し | `catch (e)` | `except ... as e` |
| `ennote` | `En-NOte` | 成れの果て | `finally` | `finally` |
| `e no "lyuzun"` | `"lyuzun"` | 理由 | `e.message` | `str(e)` |

## 値と演算

| 書き方 | 清書 | 意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| `tadact` | `TadACt` | 正しさ | `true` | `True` |
| `ulieso` | `ULieSo` | 嘘 | `false` | `False` |
| `vanita` | `Vani·tA` | 虚無 | `null` | `None` |
| `xs no 0` | `N'o` | の | `xs[0]` | `xs[0]` |
| `x exi xs` | `eXi'` | 存在する | `xs.includes(x)` | `x in xs` |
| `x exi` | `eXi'` | 存在する | `x != null` | `x is not None` |
| `nii` | `Nii` | 否定 | `!` | `not` |
| `sond` | `S~ond'` | そして | `&&` | `and` |
| `owk` | `Owk` | または | `\|\|` | `or` |
| `senaji` | `SeNaji` | 同じ | `===` | `==` |

`exi` の優先順位は比較と同じです。右に式があれば「含まれるか」を、なければ「値があるか」を調べます。「含まれない」は `nii x exi xs` と書きます。

コメントは `#` から行末までです。清書の1行目に `#!` (シバン) も書けます。

## 真偽値

条件に書けるのは `tadact` と `ulieso` だけです。`mof`、`azagara`、`shoubek`、`nii`、`sond`、`owk` にほかの値を渡すと `imJune` エラーです。Go や Rust と同じ考え方です。

| 書き方 | fxma | JavaScript | Python |
| --- | --- | --- | --- |
| `mof 0` | `imJune` | 偽 | 偽 |
| `mof ""` | `imJune` | 偽 | 偽 |
| `mof []` | `imJune` | 真 | 偽 |

- 値があるかは `x exi`、空かどうかは `nang(xs) senaji 0` で調べます。
- `sond` と `owk` は真偽値を返します。JavaScript の `a || b` のように、左右の値をそのまま返すことはしません。
- 真偽値以外を受け付けないのは、`x % 2 shoubek` のような書き間違いを止めるためです。そのまま通すと、バックトラックが間違った候補で進んでしまいます。

## 等しさ

`senaji` は Python の `==` と同じ比べ方をします。ただし、真偽値と数は別の型として扱います。JavaScript の `===` とは、配列と辞書を中身で比べる点が違います。

| 書き方 | 結果 | 理由 |
| --- | --- | --- |
| `1 senaji 1.0` | `tadact` | 整数と小数は値で比べる |
| `"1" senaji 1` | `ulieso` | 型を変換しない。エラーにもしない |
| `tadact senaji 1` | `ulieso` | 真偽値は数ではない (Python では真になる) |
| `[1, 2] senaji [1, 2]` | `tadact` | 配列と辞書は中身を比べる |
| `vanita senaji vanita` | `tadact` | — |

- 関数、モジュール、大魔女、魔女は、同じものかどうかで比べます。
- 小数の NaN は、JavaScript や Python と同じく自分自身とも等しくありません。

## 演算子

演算子は記号で書きます。優先順位は Python と同じです。等しさと論理演算だけは、記号ではなく上の表の語を使います。次の表は優先順位の高い順です。

| 優先 | 演算子 | 意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| 1 | `**` | べき乗 | `**` | `**` |
| 2 | `+x` `-x` `~x` | 符号、ビット反転 | 同じ | 同じ |
| 3 | `*` `/` `//` `%` | 掛け算、割り算、切り捨て割り算、余り | `//` は `Math.floor(a / b)` | 同じ |
| 4 | `+` `-` | 足し算、引き算 | 同じ | 同じ |
| 5 | `<<` `>>` | ビットシフト | 同じ | 同じ |
| 6 | `&` | ビット AND | 同じ | 同じ |
| 7 | `^` | ビット XOR | 同じ | 同じ |
| 8 | `\|` | ビット OR | 同じ | 同じ |
| 9 | `<` `>` `<=` `>=` `senaji` | 比較 | `senaji` は `===` | `senaji` は `==` |
| 10 | `nii` | 否定 | `!` | `not` |
| 11 | `sond` | かつ | `&&` | `and` |
| 12 | `owk` | または | `\|\|` | `or` |

- `/` の結果はいつも小数で、`//` は切り捨てです (Python と同じ)。
- 数の書き方は [数値](#数値) にあります。
- `+` は文字列の連結に、`*` は文字列の繰り返しにも使えます。
- 比較は `a < b <= c` のようにつなげられます (Python と同じ)。
- 「等しくない」は `nii a senaji b` と書きます。`nii` は比較より優先順位が低いので、`nii (a senaji b)` と同じ意味です。`==` と `!=` は使えません。書き方を1つにそろえ、論理演算を `sond` と `owk` の語で書くのと合わせるためです。

## 複合代入

演算子を `izu` の直前に付けます。JavaScript や Python の `*=` に当たります。

| 書き方 | 清書 | JavaScript・Python |
| --- | --- | --- |
| `x +izu 1` | `x +Iz'u 1` | `x += 1` |
| `x *izu 2` | `x *Iz'u 2` | `x *= 2` |

使えるのは `+izu` `-izu` `*izu` `/izu` `//izu` `%izu` `**izu` `&izu` `|izu` `^izu` `<<izu` `>>izu` です。`izu` は式の頭に来ないので、`x * izu 2` のように空白を挟んでも意味は1つに決まります。清書では空白なしにそろえます。`showence` で作った定数に使うと `blasTOac` エラーです。

fxma:

```
x izu 10
x *izu 2
x -izu 5
x redogh  # 15
```

JavaScript:

```javascript
let x = 10;
x *= 2;
x -= 5;
console.log(x); // 15
```

Python:

```python
x = 10
x *= 2
x -= 5
print(x)  # 15
```

## 数値

数はアラビア数字で書きます。辞書の数詞は1〜13・100・1000だけで0もないので、キーワードにはしません。

| 種類 | 書き方の例 | 中身 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| 整数 | `0` `42` `1_000_000` | 64ビット整数 | `number` (または `BigInt`) | `int` |
| 小数 | `3.14` `0.5` `1e-3` `2.5E10` | 64ビット浮動小数点数 | `number` | `float` |

- `_` は桁の区切りで、値は変わりません (JavaScript や Python と同じ)。
- 負の数は単項演算子の `-` で書きます (`-5`)。
- 整数と小数を計算すると、結果は小数になります。`/` の結果はいつも小数です。
- 整数が64ビットの範囲を超えると `Bavellabion` エラーです。Python のような多倍長整数は、最初の版では扱いません。
- 小数の誤差は JavaScript や Python と同じです。`0.1 + 0.2` は `0.30000000000000004` になります。

## 文字列

文字列は `"` で囲みます。書き方は CoffeeScript や Ruby と同じです。`'` は骨格では取り除く記号なので、文字列の囲みには使いません。

| 書き方 | 意味 | JavaScript | Python |
| --- | --- | --- | --- |
| `"a\nb"` | エスケープ | `"a\nb"` | `"a\nb"` |
| `"x は #{x} です"` | 値の埋め込み | `` `x は ${x} です` `` | `f"x は {x} です"` |
| `"""` で囲む | 複数行の文字列 | `` ` `` で囲む | `"""` で囲む |

- エスケープは `\n` (改行)、`\t` (タブ)、`\r`、`\\`、`\"`、`\#`、`\u{1F52E}` (コードポイント) です。
- `#{}` の中には式を書けます。数は `coword no salook` と同じ形で文字列になります。`#{` をそのまま書きたいときは `\#{` と書きます。
- `#{}` の中はコードなので、ビルドでキーワードを清書します。それ以外の文字列の中身は、ビルドで書き換えません。
- `"""` で囲んだ文字列は、最初と最後の改行と、全部の行に共通するインデントを取り除きます (CoffeeScript と同じ)。埋め込みも使えます。
- 1文字は Unicode のコードポイント1つです。Python と同じで、`nang("魔女")` は `2` です。絵文字の組み合わせのように、見た目は1文字でもコードポイントが複数あるものは、その数だけ数えます。
- 配列と同じく、`s no 0` で1文字を、`s no 1..3` で一部を取り出せます。文字列は書き換えられないので、`s no 0 izu "x"` は `imJune` エラーです。
- `"ab" exi s` は、`s` に `ab` が含まれるかを調べます。JavaScript の `s.includes("ab")`、Python の `"ab" in s` と同じです。
