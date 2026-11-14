# キーワード一覧

キーワードはすべて辞書にある語です。「意味」は辞書の訳で、「Python」「JavaScript」はそれに当たる書き方です。

## 文と制御

| 書き方 | 清書 | 意味 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| `x izu 1` | `Iz'u` | is | `x = 1` | `let x = 1` |
| `showence x izu 1` | `showence` | 証拠 | 再代入できない定数 | `const x = 1` |
| `x redogh` | `ReDogh'` | 伝える | `print(x)` | `console.log(x)` |
| `x canosiri` | `CaNoSiri` | 罵り、叫び | `print(x, file=sys.stderr)` | `console.error(x)` |
| `x izu "名前は?" zanfession` | `zanfession` | 懺悔 | `x = input("名前は?")` | `x = prompt("名前は?")` |
| `mof 条件` | `Mof` | もし | `if` | `if` |
| `siva` | `SiVa` | しかし | `else` | `else` |
| `siva mof 条件` | `SiVa Mof` | — | `elif` | `else if` |
| `条件 azagara` | `aZaGara` | 〜しながら | `while` | `while` |
| `xs aste x azagara` | `aSte` | 〜として | `for x in xs` | `for (const x of xs)` |
| `keip` | `Kei-P` | 継続 | `continue` | `continue` |
| `ketzen` | `KetZen` | 結末 | `break` | `break` |
| `finend` | `FineNd` | 終わり | ブロックの終わり | `}` |
| `weho f(n)` | `We-Ho` | 方法 | `def f(n)` | `function f(n)` |
| `n itswa` | `itsWA` | 以上だ | `return n` | `return n` |
| `abiryoku weho f()` | `AbiRyoku` | 能力 | 副作用がある関数の印 | 対応なし |
| `#` | `#` | — (記号) | `#` (コメント) | `//` |

unlessやuntilに当たる専用の語は作らず、`nii` (否定) で書きます。`nii` は比較より結び付きが弱いので、括弧は要りません。

| 書き方 | Python | JavaScript |
| --- | --- | --- |
| `mof nii x senaji 1` | `if not x == 1` (Rubyなどの `unless`) | `if (!(x === 1))` |
| `nii x > 10 azagara` | `while not x > 10` (Rubyなどの `until`) | `while (!(x > 10))` |

## 世界観由来の構文

| 書き方 | 清書 | 意味 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| `x izu 1 owk 2 mekamo` | `meKaMo` | 〜かも | 対応なし (候補を並べる) | 対応なし |
| `条件 shoubek` | `ShouBeK` | 〜するべき | `assert` (ただし失敗すると巻き戻る) | `console.assert()` に近い |
| `sentky "名前"` | `sent-ky` | 裁き | テスト関数 | `test()` (Jestなど) |
| `persemon "名前"` | `Per·SemOn` | 見世物 | ベンチマーク関数 | `bench()` (Vitest) |
| `"理由" accuhearts` | `AccUhEarts` | 告発 | `raise` | `throw` |
| `teitance` | `Tei-Tance` | 抵抗 | `try` | `try` |
| `youlugive e` | `yOuLugiVe` | 許し | `except ... as e` | `catch (e)` |
| `ennote` | `En-NOte` | 成れの果て | `finally` | `finally` |
| `e no "lyuzun"` | `"lyuzun"` | 理由 | `str(e)` | `e.message` |

## 値と演算

| 書き方 | 清書 | 意味 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| `tadact` | `TadACt` | 正しさ | `True` | `true` |
| `ulieso` | `ULieSo` | 嘘 | `False` | `false` |
| `vanita` | `Vani·tA` | 虚無 | `None` | `null` |
| `xs no 0` | `N'o` | の | `xs[0]` | `xs[0]` |
| `x exi xs` | `eXi'` | 存在する | `x in xs` | `xs.includes(x)` |
| `x exi` | `eXi'` | 存在する | `x is not None` | `x != null` |
| `nii` | `Nii` | 否定 | `not` | `!` |
| `sond` | `S~ond'` | そして | `and` | `&&` |
| `owk` | `Owk` | または | `or` | `\|\|` |
| `senaji` | `SeNaji` | 同じ | `==` | `===` |

`exi` は比較と同じ優先順位です。右に式があれば「含まれるか」、なければ「値があるか」を調べます。「含まれない」は `nii x exi xs` と書きます。

コメントは `#` から行末までです。Pythonと同じで、清書の1行目の `#!` もそのまま使えます。

## 真偽値

条件に書けるのは `tadact` と `ulieso` だけです。`mof`、`azagara`、`shoubek`、`nii`、`sond`、`owk` に他の値を渡すと `imJune` エラーです。GoやRustと同じ考え方です。

| 書き方 | fxma | Python | JavaScript |
| --- | --- | --- | --- |
| `mof 0` | `imJune` | 偽 | 偽 |
| `mof ""` | `imJune` | 偽 | 偽 |
| `mof []` | `imJune` | 偽 | 真 |

- 値があるかは `x exi`、空かどうかは `nang(xs) senaji 0` で調べます。
- `sond` と `owk` は真偽値を返します。Pythonの `a or b` のように、左右の値をそのまま返すことはしません。
- `x % 2 shoubek` のような書き間違いをそのまま通すと、巻き戻しが間違った候補で進んでしまいます。止めて知らせるために、真偽値以外は受け付けません。

## 等しさ

`senaji` はPythonの `==` と同じ比べ方をします。ただし真偽値と数は別の型として扱います。

| 書き方 | 結果 | 理由 |
| --- | --- | --- |
| `1 senaji 1.0` | `tadact` | 整数と小数は値で比べる |
| `"1" senaji 1` | `ulieso` | 型を変換しない。エラーにもしない |
| `tadact senaji 1` | `ulieso` | 真偽値は数ではない。Pythonでは真になる |
| `[1, 2] senaji [1, 2]` | `tadact` | 配列と辞書は中身を比べる |
| `vanita senaji vanita` | `tadact` | — |

- 関数、モジュール、大魔女、魔女は、同じものかどうかで比べます。
- 小数のNaNは、Pythonと同じく自分自身とも等しくありません。

## 演算子

演算子は記号で書き、優先順位はPythonと同じです。等しいかどうかと論理演算だけは、上の表の語を使います。表は優先順位の高い順です。

| 優先 | 演算子 | 意味 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| 1 | `**` | べき乗 | `**` | `**` |
| 2 | `+x` `-x` `~x` | 符号、ビット反転 | 同じ | 同じ |
| 3 | `*` `/` `//` `%` | 掛け算、割り算、切り捨て割り算、余り | 同じ | `//` は `Math.floor(a / b)` |
| 4 | `+` `-` | 足し算、引き算 | 同じ | 同じ |
| 5 | `<<` `>>` | ビットシフト | 同じ | 同じ |
| 6 | `&` | ビットAND | 同じ | 同じ |
| 7 | `^` | ビットXOR | 同じ | 同じ |
| 8 | `\|` | ビットOR | 同じ | 同じ |
| 9 | `<` `>` `<=` `>=` `senaji` | 比較 | `senaji` は `==` | `senaji` は `===` |
| 10 | `nii` | 否定 | `not` | `!` |
| 11 | `sond` | かつ | `and` | `&&` |
| 12 | `owk` | または | `or` | `\|\|` |

- `/` の結果は常に小数、`//` は切り捨てです。Pythonと同じです。
- 数値の書き方は、下の「数値」を見てください。
- `+` は文字列の連結、`*` は文字列の繰り返しにも使えます。
- 比較は `a < b <= c` のようにつなげられます。Pythonと同じです。
- 「等しくない」は `nii a senaji b` と書きます。`nii` は比較より優先順位が低いので、`nii (a senaji b)` と同じ意味になります。記号の `==` と `!=` は使えません。書き方を1つにそろえ、論理演算を `sond` と `owk` の語で書くのと合わせるためです。

## 複合代入

演算子を `izu` の直前に付けます。PythonやJavaScriptの `*=` に当たります。

| 書き方 | 清書 | Python・JavaScript |
| --- | --- | --- |
| `x +izu 1` | `x +Iz'u 1` | `x += 1` |
| `x *izu 2` | `x *Iz'u 2` | `x *= 2` |

使えるのは `+izu` `-izu` `*izu` `/izu` `//izu` `%izu` `**izu` `&izu` `|izu` `^izu` `<<izu` `>>izu` です。`izu` は式の頭に来ないので、`x * izu 2` のように空白を挟んでも曖昧になりません。清書では空白なしにそろえます。`showence` で作った定数に使うと `blasTOac` エラーです。

フィクスマージ語 (ソース):

```
x izu 10
x *izu 2
x -izu 5
x redogh  # 15
```

Python:

```python
x = 10
x *= 2
x -= 5
print(x)  # 15
```

JavaScript:

```javascript
let x = 10;
x *= 2;
x -= 5;
console.log(x); // 15
```

## 数値

数はアラビア数字で書きます。辞書の数詞は1〜13・100・1000しかなく、0もないため、キーワードにはしません。

| 種類 | 書き方の例 | 中身 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| 整数 | `0` `42` `1_000_000` | 64ビット整数 | `int` | `number` (または `BigInt`) |
| 小数 | `3.14` `0.5` `1e-3` `2.5E10` | 64ビット浮動小数点数 | `float` | `number` |

- `_` は桁の区切りで、値には影響しません。Pythonと同じです。
- 負の数は、単項演算子の `-` で書きます (`-5`)。
- 整数と小数を計算すると、結果は小数になります。`/` の結果は常に小数です。
- 整数が64ビットの範囲を超えたら `Bavellabion` エラーにします。Pythonのような多倍長整数は、最初の版では扱いません。
- 小数の誤差はPythonやJavaScriptと同じです。`0.1 + 0.2` は `0.30000000000000004` になります。

## 文字列

文字列は `"` で囲みます。書き方はCoffeeScriptやRubyと同じです。`'` は骨格では取り除く記号なので、文字列の囲みには使いません。

| 書き方 | 意味 | Python | JavaScript |
| --- | --- | --- | --- |
| `"a\nb"` | エスケープ | `"a\nb"` | `"a\nb"` |
| `"x は #{x} です"` | 値の埋め込み | `f"x は {x} です"` | `` `x は ${x} です` `` |
| `"""` で囲む | 複数行の文字列 | `"""` で囲む | `` ` `` で囲む |

- エスケープは `\n` (改行)、`\t` (タブ)、`\r`、`\\`、`\"`、`\#`、`\u{1F52E}` (文字の番号) です。
- `#{}` の中には式を書けます。数は `coword no salook` と同じ形で文字列になります。`#{` をそのまま書きたいときは `\#{` と書きます。
- `#{}` の中はコードなので、ビルドでキーワードを清書します。その外の文字列はビルドで書き換えません。
- `"""` で囲んだ文字列は、最初と最後の改行と、全部の行に共通する字下げを取り除きます。CoffeeScriptと同じです。埋め込みも使えます。
- 1文字は、Unicodeの文字の番号 (コードポイント) 1つです。Pythonと同じで、`nang("魔女")` は `2` です。絵文字の組み合わせのように、見た目は1文字でも番号が複数あるものは、その数だけ数えます。
- 配列と同じく、`s no 0` で1文字を、`s no 1..3` で一部を取り出せます。文字列は書き換えられないので、`s no 0 izu "x"` は `imJune` エラーです。
- `"ab" exi s` は、`s` に `ab` が含まれるかを調べます。Pythonの `"ab" in s` と同じです。
