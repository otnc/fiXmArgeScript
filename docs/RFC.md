# fiXmArge Script 仕様書 (草案)

## 概要

fiXmArge Script (略称 fxma) は、『魔法少女ノ魔女裁判』 (まのさば) の架空言語「フィクスマージ語」の単語だけでキーワードを作った、小さなインタプリタ言語です。目標は、FizzBuzzと簡単な推理パズルが書けることです。

名前は、元にした架空言語の名前「fiXmArge (フィクスマージ語)」から取りました。略称の fxma は、**f**i**XmA**rge の f・X・m・A です。コマンドは `fixmarge` (npm のパッケージと VS Code 拡張も同じ名前)、ファイルの拡張子は下書きが `.fxma`、清書が `.fxmac`、パッケージのファイルは `fixmarge.mod` と `fixmarge.sum` です。

架空言語そのものや公式の作品と混同されないよう、「非公式のファンメイド作品」であることを名前の近くに必ず書きます。また、「Script」を省いて「fiXmArge」だけでこの言語を呼ばないようにします。

設計方針は次の4つです。

- **語彙は公式に寄せる**: キーワードは、設定資料集の翻訳表を写したとされるファン辞書の語から選ぶ。辞書にない語は作らない。
- **書くときは記号なし**: `Iz'u` は `izu` と書く。ビルドすると正式な綴りに清書される。
- **ゲームの流れを言語にする**: 選択肢とバッドエンド、証拠品、魔法、魔女裁判を、それぞれバックトラック、定数、副作用の宣言、テストに対応させる。
- **普通の言語として読める**: 各機能にPythonとJavaScriptの書き方を並べて、フィクスマージ語を知らなくても意味がわかるようにする。

語彙の出典と信頼度は、最後の「出典」にまとめています。

## 表記

ソースコードは英数字だけで書き、ビルドで正式な綴りに清書します。辞書の264語は、記号と空白を抜いて小文字にしても綴りがぶつからないため、元の綴りに一意に戻せます。この記号なしの形を「骨格」と呼びます。

| 正式な綴り | 骨格 (ソースに書く形) | 取り除くもの |
| --- | --- | --- |
| `Iz'u` | `izu` | `'` |
| `S~ond'` | `sond` | `~` と `'` |
| `Syu Tain` | `syutain` | 空白 |
| `Kei-P` | `keip` | `-` |
| `Vani·tA` | `vanita` | `·` |

### 規則

1. ソースに使えるのは、英数字、`_`、演算子の記号、コメントの #、`"` で囲む文字列です。文字列とコメントの中は、日本語も含めて自由に書けます。
2. 下書きのキーワードは小文字だけで書きます。`IZU` のように大文字を含むと `AyAmis` (構文エラー) です。清書では、辞書どおりの大文字小文字 (`Iz'u`、`AccUhEarts` など) に戻します。
3. 変数名や関数名に、キーワードの骨格と同じ名前は使えません。大文字を混ぜただけの名前 (Izu など) も、紛らわしいので使えません。
4. 下書き (ソース) と清書は拡張子で分けます。下書きが `.fxma`、清書が `.fxmac` です。インタプリタが実行するのは清書です。下書きは必ずビルドを通してから実行されます (`fixmarge keeon` はメモリ上でビルドして実行する)。清書の読み込みは厳密で、正式な綴り以外 (方言や骨格) は受け付けません。
5. ビルドは何度かけても同じ結果になります。
6. 標準パッケージの関数名や辞書のキーで、JavaScriptでは複数の語からなる名前 (`isAbsolute` など) は、辞書の語をJavaScriptと同じ順に並べて組みます。2語目から頭を大文字にしてつなぎます (`izuDefutai`)。部品にはキーワードの語も使えます。清書では、各部品を正式な綴りに戻してそのままつなぎます (`Iz'udeFuTai`)。部品に分けられるので、同じ語を別の名前でも使い回せます。

記号を書かないので、`'` と文字列、`-` と引き算がぶつかりません。空白入りの語もソースでは1語になるため、字句解析が単純になります。

### 方言と表記ブレ

作詞者のSLAVE.V-V-R氏は、発音の違いを「方言の分岐の原初」のようなものだと説明しています。そこで表記ブレはエラーにせず、方言として受け入れます。

- **標準語**: 辞書CSVの綴りを正式な綴りとします。
- **方言表**: 異なる綴りの骨格を、標準語の骨格への別名として登録します。例: `lowhel` → `lowhell`、`prizujeen` → `prisjeen`。
- **警告**: 方言でも実行できますが、`DowGuil lowhel → LowhEll` のような警告が出ます。
- **ビルド**: 清書は常に標準語になります。
- **やらないこと**: 編集距離によるあいまい一致。打ち間違いが別の語に化けるのを防ぐためです。

語を足したときに骨格がぶつからないかは、CIで検査します。

## キーワード一覧

キーワードはすべて辞書にある語です。「意味」は辞書の訳で、「Python」「JavaScript」はそれに当たる書き方です。

### 文と制御

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

### 世界観由来の構文

| 書き方 | 清書 | 意味 | Python | JavaScript |
| --- | --- | --- | --- | --- |
| `x izu 1 owk 2 mekamo` | `meKaMo` | 〜かも | 対応なし (候補を並べる) | 対応なし |
| `条件 shoubek` | `ShouBeK` | 〜するべき | `assert` (ただし失敗すると巻き戻る) | `console.assert()` に近い |
| `sentky "名前"` | `sent-ky` | 裁き | テスト関数 | `test()` (Jestなど) |
| `"理由" accuhearts` | `AccUhEarts` | 告発 | `raise` | `throw` |
| `teitance` | `Tei-Tance` | 抵抗 | `try` | `try` |
| `youlugive e` | `yOuLugiVe` | 許し | `except ... as e` | `catch (e)` |
| `ennote` | `En-NOte` | 成れの果て | `finally` | `finally` |
| `e no "lyuzun"` | `"lyuzun"` | 理由 | `str(e)` | `e.message` |

### 値と演算

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

### 真偽値

条件に書けるのは `tadact` と `ulieso` だけです。`mof`、`azagara`、`shoubek`、`nii`、`sond`、`owk` に他の値を渡すと `imJune` エラーです。GoやRustと同じ考え方です。

| 書き方 | fxma | Python | JavaScript |
| --- | --- | --- | --- |
| `mof 0` | `imJune` | 偽 | 偽 |
| `mof ""` | `imJune` | 偽 | 偽 |
| `mof []` | `imJune` | 偽 | 真 |

- 値があるかは `x exi`、空かどうかは `nang(xs) senaji 0` で調べます。
- `sond` と `owk` は真偽値を返します。Pythonの `a or b` のように、左右の値をそのまま返すことはしません。
- `x % 2 shoubek` のような書き間違いをそのまま通すと、巻き戻しが間違った候補で進んでしまいます。止めて知らせるために、真偽値以外は受け付けません。

### 等しさ

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

### 演算子

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

### 複合代入

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

### 数値

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

### 文字列

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

## 基本構文

1行に1文を書き、ブロックは `finend` で閉じます。インデントに意味はありません。関数の呼び出しは `f(a, b)` の形です。

キーワードを置く位置は、その語が日本語で文のどこに来る語から作られたかで決めます。

| 位置 | 決め方 | キーワード |
| --- | --- | --- |
| 文の後ろ | 動詞、語尾 (〜かも、〜べき)、「〜ながら」 | `mekamo` `shoubek` `itswa` `azagara` `redogh` `canosiri` `zanfession` `accuhearts` `werer` `praylee` |
| 文の頭 | 文の頭に来る語 (もし、しかし)、ブロックや宣言の見出し | `mof` `siva` `weho` `abiryoku` `itsver` `oqun` `curoi` `sentky` `teitance` `youlugive` `ennote` `showence` `nii` |
| 2つのものの間 | 2つをつなぐ語 | `izu` `no` `exi` `senaji` `sond` `owk` `aste` `batte` |
| 1語で文になる | 位置は問題にならない | `keip` `ketzen` `finend` |

`nii` (否定) は、日本語の「ない」なら文の後ろになります。ただし辞書に語源がなく、`mof nii` でunlessを書くため、文の頭に置きます。

### 後ろに置く語のかかる範囲

文の後ろに置く語は2種類で、どちらも自分より左側をまるごと受け取ります。日本語の文末の動詞が文全体にかかるのと同じです。

| 種類 | キーワード | かかる範囲 | 例 |
| --- | --- | --- | --- |
| 式の終わりの語 | `mekamo` `praylee` `zanfession` | 式全体。`izu` があればその右側全体 | `y izu f() praylee` は `y = await f()` |
| 文の終わりの語 | `redogh` `canosiri` `itswa` `accuhearts` `shoubek` `azagara` `werer` | 文全体 | `x + 1 redogh` は `print(x + 1)` |

- 1つの文に、式の終わりの語と文の終わりの語を続けて書けます。`f() praylee redogh` は `print(await f())` です。
- 式の途中で使うときは括弧で囲みます。`(f() praylee) + 1` のように書きます。
- `zanfession` は左側を入力のプロンプトとして出し、読んだ1行を返します。`name izu "名前は?" zanfession` はPythonの `name = input("名前は?")` です。プロンプトがいらなければ `name izu zanfession` と書きます。

### 変数と出力

フィクスマージ語 (ソース):

```
x izu 20
x + 1 redogh
```

Python:

```python
x = 20
print(x + 1)
```

JavaScript:

```javascript
let x = 20;
console.log(x + 1);
```

### 変数のスコープ

変数はブロックごとに分かれます。JavaScriptの `let` と同じです。ただし `izu` は宣言と代入を兼ねるので、次の決まりで見分けます。

- `izu` の左の名前が外のブロックにあれば、その変数に代入します。なければ、今のブロックに新しく作ります。
- ブロックの中で作った変数は、`finend` の後では使えません。使うと `Nii eXi'` エラーです。
- 関数の中からは、外の変数を読めますが `izu` はできません。`Su·Ready` エラーです。関数の中で外と同じ名前の変数を作ることもできません。外の値を変えたいときは、戻り値か、配列や辞書を使います。
- 引数だけは、外と同じ名前を付けられます。関数の中では、その名前は引数を指します。
- `weho` で定義した関数は、JavaScriptの `function` 宣言と同じく、同じブロックの中なら定義より前の行から呼べます (巻き上げ)。
- `aste` で受け取るループの変数は、1周ごとに新しく作ります。JavaScriptの `for (const x of xs)` と同じです。
- `showence` で作る定数も、同じ決まりに従います。

```
count izu 0

weho f()
  count redogh  # 読むだけなら使える
  count +izu 1  # → Su·Ready
finend

mof count senaji 0
  count +izu 1  # 外の count に代入する
  y izu 1       # このブロックだけの変数
finend
y redogh        # → Nii eXi'
```

関数の中で外の変数に代入できると、外にたまたま同じ名前があったときに黙って書き換えてしまいます。宣言の語を持たないCoffeeScriptで問題になった点なので、関数の境目でだけ止めます。

### 条件分岐

フィクスマージ語 (ソース):

```
mof x > 10
  "big" redogh
siva mof x senaji 10
  "ten" redogh
siva
  "small" redogh
finend
```

Python:

```python
if x > 10:
    print("big")
elif x == 10:
    print("ten")
else:
    print("small")
```

JavaScript:

```javascript
if (x > 10) {
  console.log("big");
} else if (x === 10) {
  console.log("ten");
} else {
  console.log("small");
}
```

### 繰り返しと関数 (FizzBuzz)

フィクスマージ語 (ソース):

```
weho fizzbuzz(n)
  mof n % 15 senaji 0
    "FizzBuzz" itswa
  siva mof n % 3 senaji 0
    "Fizz" itswa
  siva mof n % 5 senaji 0
    "Buzz" itswa
  finend
  n itswa
finend

x izu 1
x <= 15 azagara
  fizzbuzz(x) redogh
  x +izu 1
finend
```

Python:

```python
def fizzbuzz(n):
    if n % 15 == 0:
        return "FizzBuzz"
    elif n % 3 == 0:
        return "Fizz"
    elif n % 5 == 0:
        return "Buzz"
    return n

x = 1
while x <= 15:
    print(fizzbuzz(x))
    x = x + 1
```

JavaScript:

```javascript
function fizzbuzz(n) {
  if (n % 15 === 0) {
    return "FizzBuzz";
  } else if (n % 3 === 0) {
    return "Fizz";
  } else if (n % 5 === 0) {
    return "Buzz";
  }
  return n;
}

let x = 1;
while (x <= 15) {
  console.log(fizzbuzz(x));
  x = x + 1;
}
```

ビルド後の清書 (関数の冒頭だけ):

```
We-Ho fizzbuzz(n)
  Mof n % 15 SeNaji 0
    "FizzBuzz" itsWA
```

`redogh` で出す形は次のとおりです。文字列はそのまま出し、数はJavaScriptの `String(x)` と同じ形 (`3.0` は `3`) にします。真偽値と `vanita` は清書の綴り (`TadACt`、`ULieSo`、`Vani·tA`) で出します。配列と辞書はJSONに近い形で、中の文字列は `"` で囲みます (`[1, 2.5, "a"]`、`{"k": TadACt}`)。`#{}` で埋め込むときも同じ形です。

この例では `redogh` を関数の外で呼んでいるので、`fizzbuzz` に `abiryoku` は必要ありません (次の章)。

### 配列

`[3, 1, 4]` で作り、要素は `xs no 0` (xsの0) で取り出します。`xs[0]` のように `[]` で添字を書く形はありません。

- 添字は0からで、`xs no -1` は最後の要素です。
- `no` は他の演算子より強く結び付きます。式を添字にするときは `xs no (i + 1)` のように括弧で囲みます。
- `+` で配列をつなげられます。
- 配列と辞書は、JavaScriptやPythonと同じく参照で扱います。別の変数に入れたり関数に渡したりしても中身は共有され、片方で書き換えるともう片方にも見えます。
- `showence` は変数への再代入を止めるだけで、配列や辞書の中身の書き換えは止めません。JavaScriptの `const` と同じです。
- 範囲の外を読むと `Nii eXi'` エラーです。
- 長さは組み込みの `nang(xs)` で取ります。
- `no` の語順 (A no B ＝ AのB) は、辞書の訳「の」に合わせたものです。公式の資料で語順はまだ確認できていません。

フィクスマージ語 (ソース):

```
xs izu [3, 1, 4]
xs no 0 izu 9
xs no -1 redogh  # 4
```

Python:

```python
xs = [3, 1, 4]
xs[0] = 9
print(xs[-1])  # 4
```

JavaScript:

```javascript
let xs = [3, 1, 4];
xs[0] = 9;
console.log(xs.at(-1)); // 4
```

### 範囲

`..` と `...` で数の範囲を書けます。CoffeeScriptと同じく、`..` は終わりを含み、`...` は含みません。刻み幅は、辞書で「〜によって、〜を通して」を意味する `batte` (清書では `Batte`) で指定します。

| 書き方 | 意味 | Python | JavaScript |
| --- | --- | --- | --- |
| `1..5` | 1から5まで | `range(1, 6)` | なし |
| `1...5` | 1から4まで | `range(1, 5)` | なし |
| `0..10 batte 2` | 0、2、…、10 | `range(0, 11, 2)` | なし |
| `10..1 batte -1` | 10、9、…、1 | `range(10, 0, -1)` | なし |
| `xs no 2..4` | 2番目から4番目 | `xs[2:5]` | `xs.slice(2, 5)` |
| `xs no 2...4` | 2番目から3番目 | `xs[2:4]` | `xs.slice(2, 4)` |
| `xs no 2..` | 2番目から最後まで | `xs[2:]` | `xs.slice(2)` |
| `xs no ..2` | 最初から2番目 | `xs[:3]` | `xs.slice(0, 3)` |
| `xs no 1..2 izu [7, 8]` | 1番目から2番目を置き換える | `xs[1:3] = [7, 8]` | `xs.splice(1, 2, 7, 8)` |

- `no` の右には、数・変数・括弧で囲んだ式と、それを `..` か `...` でつないだ範囲を書けます。
- 範囲は `aste ... azagara` の繰り返しにそのまま使えます。

フィクスマージ語 (ソース):

```
# 0から10までの偶数
0..10 batte 2 aste x azagara
  x redogh
finend
```

Python:

```python
for x in range(0, 11, 2):
    print(x)
```

JavaScript:

```javascript
for (let x = 0; x <= 10; x += 2) {
  console.log(x);
}
```

### no で取り出す

`A no B` は「AのB」で、Aの中からBを取り出します。JavaScriptの `A[B]` と同じく、右側は値として読みます。

| 書き方 | 取り出すもの | JavaScript |
| --- | --- | --- |
| `xs no 0` | 配列の要素 | `xs[0]` |
| `xs no i` | 変数 `i` の値を添字にした要素 | `xs[i]` |
| `xs no 2..4` | 配列の一部 | `xs.slice(2, 5)` |
| `s no 0` | 文字列の1文字 | `[...s][0]` |
| `d no "lang"` | 辞書の値 | `d["lang"]` |
| `e no "lyuzun"` | 例外の理由 | `e["message"]` |
| `dv no "murm"(w)` | 大魔女の関数を呼ぶ | `dv["murm"](w)` |
| `util no double(21)` | 読み込んだモジュールの関数を呼ぶ | `util.double(21)` |

- 右に識別子を書くと、変数として読みます。`xs no i` は `i` の値で取り出します。
- 大魔女、魔女、例外の中身は、辞書と同じく文字列のキーで取り出します。JavaScriptの `obj.name` が `obj["name"]` の省略形なのと同じで、中身と辞書のキーを区別しません。キーは文字列なのでビルドで書き換えず、骨格のまま書きます。
- 左がモジュールのときだけ、右の識別子をモジュールの中の名前として読みます。モジュールの中の名前は読み込んだ時点で決まっているので、実行する前に分かります。`"./util" no double werer` の `double` も同じです。
- 関数を取り出すと、取り出し元がくっついたままになります。JavaScriptの `obj.f.bind(obj)` と同じです。`dv no "murm"` を他の関数に渡すと、受け取った側はその大魔女の力で魔女を殺せます。
- 配列に文字列のキーを使ったり、ないキーを読んだりすると `Nii eXi'` エラーです。
- 添字と範囲には、`2.0` のように端数のない小数も使えます。JavaScriptの `xs[2.0]` と同じです。`meme no racdam` や `coword no salook` は小数を返すので、その結果をそのまま添字にできます。`2.5` のように端数があると `imJune` エラーです。

## 世界観由来の仕様

まのさばは、牢屋敷での日常 (ADVパート) →証拠品を集める捜査→発言の矛盾を暴く魔女裁判、という流れで進みます。この言語では、その要素を次の4つの機能に対応させます。

| ゲームの要素 | 言語の機能 | 普通の言語でいうと |
| --- | --- | --- |
| バッドエンドになっても直前の選択肢に戻れる | 選択肢と証明 | バックトラック (`amb` 演算子) |
| 捜査で集める証拠品 | 証拠 | 定数 |
| 少女たちの固有の魔法 | 能力 | 副作用の宣言 |
| 魔女裁判と処刑 | 裁き、告発、処刑 | テストと例外 |

13人の少女にちなんで「1スコープに変数は13個まで」という案もありましたが、使いにくいため採用しません。

### 選択肢とバッドエンド

この言語のいちばんの特徴です。

1. `a owk b owk c mekamo` は、候補のうち最初の値を返し、その場所を「選択肢」として覚えます。`mekamo` で終わる文の中では、`owk` は「または」ではなく候補の区切りになります。
2. その後の `条件 shoubek` が偽になると、バッドエンドです。直前の選択肢まで巻き戻り、次の候補で続きを実行します。
3. 巻き戻すと、変数の値もその選択肢の時点に戻ります。失敗した道筋での出力は捨てられます。
4. 候補を使い切ったら、もうひとつ前の選択肢に戻ります。戻れる選択肢がないときは `Co-Jundic` (矛盾) で終了します。

例: 証言から魔女を当てる。

フィクスマージ語 (ソース):

```
# 容疑者は1〜3番。魔女は1人
witch izu 1 owk 2 owk 3 mekamo
# 証言A：魔女は1番ではない
nii witch senaji 1 shoubek
# 証言B：魔女は奇数番
witch % 2 senaji 1 shoubek
witch redogh
```

Python:

```python
for witch in [1, 2, 3]:
    if witch == 1:
        continue
    if witch % 2 != 1:
        continue
    print(witch)
    break
else:
    raise Exception("Co-Jundic")
```

JavaScript:

```javascript
const witch = [1, 2, 3].find((w) => w !== 1 && w % 2 === 1);
if (witch === undefined) {
  throw new Error("Co-Jundic");
}
console.log(witch);
```

どれも `3` を出力します。フィクスマージ語版は「魔女は1か2か3かも。1ではないべき。奇数であるべき。魔女を伝える」と読めます。PythonやJavaScriptではループや絞り込みを自分で書きますが、この言語では候補と条件を並べるだけで済みます。

### 証拠品 (定数)

`showence` で作った変数は再代入できません。証拠を書き換えようとすると `blasTOac` (冒涜的) エラーになります。証拠を書き換えるのは裁判への冒涜、という見立てです。

```
showence limit izu 15
limit izu 20  # → blasTOac
```

Pythonでいうと `limit: Final = 15` に近いですが、Pythonは実行時には止めません。JavaScriptの `const limit = 15;` と同じく、この言語は止めます。

### 魔法 (副作用の宣言)

入出力 (`redogh`、`zanfession`) をする関数は、`abiryoku` を付けて定義します。`abiryoku` 付きの関数を呼ぶ関数も同じです。付け忘れると、実行前に `kiNbOO` (禁忌) エラーになります。

```
abiryoku weho greet(name)
  "weRer, " + name redogh
finend
```

Pythonには対応する仕組みがありません。Haskellの `IO` 型を小さくしたものです。巻き戻しで出力を捨てるためにも、どこで入出力が起きるかがわかっていると実装が楽になります。

関数を値として渡すときは、渡した関数で決まります。`abiryoku` 付きの関数を引数に渡した呼び出しは、それ自体が入出力をする呼び出しとして扱います。

- `varizama no tarip(xs, double)` は、`double` に `abiryoku` がなければ、`abiryoku` なしの関数の中でも書けます。
- `varizama no tarip(xs, greet)` は、`greet` が `abiryoku` 付きなので、`abiryoku` 付きの関数の中でしか書けません。
- 変数に入った関数のように実行前に分からないときは、呼んだ時点で調べて、`abiryoku` なしの関数の中なら `kiNbOO` エラーにします。

### 魔女裁判 (テスト) と処刑 (例外)

`sentky` ブロックの中の `shoubek` は巻き戻しをせず、テストとして扱われます。全部通れば無罪、ひとつでも落ちれば有罪です。

フィクスマージ語 (ソース):

```
sentky "fizzbuzz"
  fizzbuzz(3) senaji "Fizz" shoubek
  fizzbuzz(15) senaji "FizzBuzz" shoubek
finend
```

Python (pytest):

```python
def test_fizzbuzz():
    assert fizzbuzz(3) == "Fizz"
    assert fizzbuzz(15) == "FizzBuzz"
```

JavaScript (Jest):

```javascript
test("fizzbuzz", () => {
  expect(fizzbuzz(3)).toBe("Fizz");
  expect(fizzbuzz(15)).toBe("FizzBuzz");
});
```

結果の表示は、判決 (`hung gement`) に続けて、無罪なら `Nii Guil-me`、有罪なら `Guil-me` を出します。有罪のときは、次の行に字下げして落ちた場所を出します。

```
hung gement double: Nii Guil-me
hung gement fizzbuzz: Guil-me
  ShouBeK util_coatain.fxma:4:28
hung gement parse: Guil-me
  imJune util.fxma:2:5 +
```

- `shoubek` が1つ落ちたら、JestやPythonの `assert` と同じく、その `sentky` はそこで止めて有罪にし、次の `sentky` に進みます。
- テストとして扱う `shoubek` は、`sentky` の本体に直接書いたものだけです。本体から呼んだ関数の中の `shoubek` は、いつもどおり巻き戻しに使います。
- 落ちた `shoubek` は `youlugive` では捕まえられません。例外や実行時エラーが捕まらずに `sentky` の外に出たときも有罪で、そのエラーを出します。
- 有罪が1つでもあれば、`fixmarge coatain` は終了コード1で終わります。

テストは、名前が `_coatain.fxma` で終わるファイルに書きます。Goの `_test.go` と同じです。`Coa Tain` (法廷) は裁き (`sentky`) を開く場所なので、この名前にしました。テストするファイルは `werer` で読み込みます。`sentky` を書けるのはこのファイルだけで、`fixmarge keeon` はこのファイルを読み込みません。

```
# util_coatain.fxma
"./util" aste util werer

sentky "double"
  util no double(2) senaji 4 shoubek
finend
```

例外は `accuhearts` (告発) で投げ、`teitance` (抵抗) と `youlugive` (許し) で捕まえます。Pythonの `raise`、`try`、`except` と同じです。誰にも許されなかった例外は、`KeI·tion` (処刑) としてプログラムを終了させます。

`youlugive e` と書くと、捕まえた例外が変数 `e` に入ります。例外の理由 (メッセージ) は `e no "lyuzun"` で取り出せます。辞書の `LyuzUn` (理由) から取ったキーです。`ennote` (`En-NOte`、成れの果て) のブロックは、例外が起きても起きなくても最後に必ず実行されます (PythonやJavaScriptの `finally`)。`canosiri` (`CaNoSiri`、罵り) は、`redogh` と同じく文の後ろに置いて、標準エラー出力に書きます。

フィクスマージ語 (ソース):

```
teitance
  "bad input" accuhearts
youlugive e
  e no "lyuzun" canosiri
ennote
  "FineNd" redogh
finend
```

Python:

```python
import sys

try:
    raise Exception("bad input")
except Exception as e:
    print(e, file=sys.stderr)
finally:
    print("FineNd")
```

JavaScript:

```javascript
try {
  throw new Error("bad input");
} catch (e) {
  console.error(e.message);
} finally {
  console.log("FineNd");
}
```

実行中に起きたエラーも、`accuhearts` で投げた例外と同じく `youlugive` で捕まえられます。PythonやJavaScriptと同じです。

- エラーの種類は、理由と同じく文字列のキーで取り出します。値は `"kYoNAihim"` のようなエラー名の清書の綴りで、`accuhearts` で投げた例外は `"AccUhEarts"` です。キーは `"chidd"` (`ChiDd`、血) で、`e no "chidd"` のように書きます。
- 捕まえられないのは、実行前に見つかるエラー (構文の `AyAmis`、`kiNbOO`、`showence` への再代入の `blasTOac`) と、`Co-Jundic`、`SaiVI` です。`Co-Jundic` は巻き戻しの行き止まりで、`SaiVI` はインタプリタのバグなので、プログラムからは直せないためです。

## 非同期処理

非同期の関数は `itsver` (いつか) を付けて定義し、終わるのを待つときは `praylee` (祈る) を文の後ろに置きます。JavaScriptの `async` と `await` と同じ仕組みです。

| 書き方 | 清書 | 辞書の意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| `itsver weho f()` | `itsver` | いつか、遅かれ早かれ | `async function f()` | `async def f()` |
| `f() praylee` | `pray·lee` | 祈る | `await f()` | `await f()` |
| `[f(), g()] praylee` | `pray·lee` | 祈る | `await Promise.all([f(), g()])` | `await asyncio.gather(f(), g())` |

- `itsver` を付けた関数を呼ぶと、すぐに「いつか終わるもの」が返ります。`praylee` で終わるのを待つと、結果が得られます。
- 配列に `praylee` を付けると、全部が終わるのを待ちます。
- `praylee` は `itsver` の関数の中と、ファイルの一番外側で使えます。それ以外で使うと `Ca-nay` エラーです。
- 実行はJavaScriptと同じく1本の流れで、`praylee` のところでだけ別の処理に切り替わります。切り替わる順番を決まったものにすることで、巻き戻しのやり直しでも同じ結果になります。
- `praylee` で待つとき、それまでの選択肢 (`mekamo`) は確定します。待っている間に外の世界が変わるためです。

細かい決まりは、JavaScriptのNode.jsに合わせています。

- `itsver` の関数の本体は、呼んだらすぐに動き始めます。最初の `praylee` まではその場で実行し、そこで呼んだ側に戻ります (`async function` と同じ)。
- `praylee` は、もう終わっているものや、`1` のような「いつか終わるもの」でない値に付けても、一度ほかの処理に順番を回します (`await 1` と同じ)。
- 順番を待っている処理は、待ち始めた順に再開します。`toim no saimoku` は、終わる時刻の早い順に、同じなら呼んだ順に終わります (`setTimeout` と同じ)。
- ファイルの一番外側を最後まで実行しても、まだ終わっていない `itsver` の関数や待機、通信が残っていれば、全部終わるまで待ってから終了します。
- `praylee` しないまま、捕まえていないエラーで終わった `itsver` の関数があれば、全部終わった後でそのエラーを出して止まります (unhandled rejection と同じ)。後で `praylee` して `youlugive` で捕まえれば止まりません。
- `praylee` で待っているのに、ほかに動ける処理も外から届くものもなく、二度と再開できないとき (自分自身の結果を待つなど) は、一番外側が待っている `praylee` の場所で `Co-Jundic` エラーになります。
- 外の世界とのやり取りが終わった順番 (待機、通信、魔女の終わり) は記録しておき、巻き戻しのやり直しでは記録どおりの順番で再開します。
- 何度も知らせが届くもの (サーバー、ファイルの見張り、シグナル、魔女の出力) は `nozuki` (気付き) に関数を渡します。知らせが届くたびに、その関数を別の処理として呼びます。関数の中で捕まえていないエラーが起きると、Node.jsと同じくプログラムが止まります。
- サーバーとファイルの見張りは、動いている限りプログラムを終わらせません。シグナルを待っているだけなら終わります (Node.jsと同じ)。

フィクスマージ語 (ソース):

```
"fxma:toim" aste toim werer

itsver abiryoku weho greet(name)
  toim no saimoku(1000) praylee
  "weRer, " + name redogh
finend

[greet("a"), greet("b")] praylee
```

JavaScript:

```javascript
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function greet(name) {
  await sleep(1000);
  console.log("weRer, " + name);
}

await Promise.all([greet("a"), greet("b")]);
```

非同期になる標準パッケージの関数は、待機 (`toim no saimoku`)、魔女 (プロセス) の終了待ち、通信の取得 (`gemon no fineclowt`) です。

## マクロ

マクロは `curoi` (`cUrOi`、呪い) で定義します。呪いは、実行される前にプログラムそのものを書き換えます。最初の版では、引数として受け取ったコードを当てはめるだけの、型にはめる方式にします (Rustの `macro_rules!` の簡単な形)。

```
curoi twice(body)
  body
  body
finend

twice("Mof" redogh)  # Mof を2回出力する
```

- 呪いの引数には値ではなく、コードのかたまりが渡ります。展開は実行の前に1度だけ行います。
- 呪いの中で作った変数は自動で別の名前になり、呼び出し側の変数とぶつかりません。
- 呪いが呪いを呼ぶのは、展開の深さの上限 (64段) までです。超えると `Bavellabion` (あふれる) エラーです。
- 展開は字句の並びの段階で、構文解析の前に行います。そのため引数には式だけでなく、`"Mof" redogh` のような文も渡せます。
- 引数が `1 + 2` のような2語以上の式なら、括弧で包んでから当てはめます。`square(1 + 2)` が `1 + 2 * 1 + 2` になることはありません。文を渡したときは包みません。
- 「呪いの中で作った変数」は、`izu` の左、`aste` の後ろ、`youlugive` の後ろに書いた名前です。展開のたびに別の名前に付け替えます。
- 呪いはファイルのどこで定義しても、そのファイル全体で使えます (`weho` の巻き上げと同じ)。同じ名前の呪いを2つ定義すると `Su·Ready` エラーです。ほかのファイルからは使えません (`oqun` は付けられません)。
- 呼び出しの引数の数が合わないと `AyAmis` エラーです。`util no twice(...)` のようにモジュールの中の名前として書いたものは、呪いとして扱いません。
- 最初の版では、呪いの本体の文字列の `#{}` の中に書いた引数の名前は置き換えません。`#{}` の中で呪いを呼ぶことはできます。
- 清書 (`.fxmac`) には展開する前の形が残ります。展開した結果は `fixmarge beaukshi --curoi` で確かめられます。
- 構文そのものを組み立てる本格的なマクロ (LispやRustの手続きマクロ) は、後の版で検討します。

## パッケージとモジュール

1つのファイルが1つのパッケージです。読み込みは「名前を付ける」か「個別に読み込む」の2通りだけで、名前を付けずに読み込む形はありません。外部パッケージの管理はGoと同じ仕組みです。

### 読み込みと公開

| やること | 書き方 | 辞書の意味 | JavaScript |
| --- | --- | --- | --- |
| 名前を付けて読み込む | `"./util" aste util werer` | として、迎える | `import * as util from "./util.js"` |
| 個別に読み込む | `"./util" no double, triple werer` | の、迎える | `import { double, triple } from "./util.js"` |
| 個別に別名で読み込む | `"./util" no double aste d werer` | の、として、迎える | `import { double as d } from "./util.js"` |
| 公開する | `oqun weho double(n)` | 開く | `export function double(n)` |
| 使う | `util no double(21)` / `double(21)` | の | `util.double(21)` / `double(21)` |

- `werer` (迎える) は動詞なので文の後ろ、`oqun` (開く) は見出しに付く語なので前に置きます。
- `"./util"` は `./util.fxma` を読みます。`"./util" werer` だけの書き方は `AyAmis` (構文エラー) です。
- 公開できるのは `oqun` を付けた関数と定数 (`showence`) だけです。付いていないものを外から使うと `Ca-nay` (できない) エラーです。
- 読み込まれるファイルには、定義 (`weho` と `showence`) と `werer` だけを書けます。処理を実行する文は、`fixmarge keeon` で直接動かすファイルにだけ書けます。それ以外の文があると `Ca-nay` エラーです。
- 読み込まれるファイルの一番外側は `abiryoku` の外として扱います。`showence` の右で入出力をする関数を呼ぶと `kiNbOO` エラーです。
- 読み込んだ名前は、JavaScriptの `import` と同じく書き換えられません。`double izu 1` や `util no limit izu 1` は `blasTOac` エラーです。
- 同じ名前を2回読み込むと `Su·Ready` エラーです。`aste` で別名を付けます。
- 清書 (`.fxmac`) を実行しているときは、読み込み先も清書を読みます。`fixmarge sinruits` は読み込み先の下書きもまとめて清書するので、そのまま動きます。
- 読み込みが循環すると `Su·Ready` (既に)、読み込み先がないと `Nii eXi'` エラーです。

フィクスマージ語 (ソース):

```
# util.fxma
oqun weho double(n)
  n * 2 itswa
finend

# main.fxma
"./util" aste util werer
util no double(21) redogh  # 42
```

JavaScript:

```javascript
// util.js
export function double(n) {
  return n * 2;
}

// main.js
import * as util from "./util.js";
console.log(util.double(21)); // 42
```

### fixmarge.mod と fixmarge.sum

Goの `go.mod` と `go.sum` と同じ役割のファイルです。`fixmarge.mod` の項目名は辞書の語で、下書きと同じく記号なしの小文字で書きます。

```
ado github.com/alice/strs
mokpose "文字列の便利関数"
neyo fxma >=0.1.0
foryulvu MIT OR Apache-2.0

neyo github.com/bob/util v1.2.0
unneyo github.com/bob/util v1.3.0
lainei github.com/bob/util ../util
```

| 項目 | Go | 辞書の意味 | 中身 |
| --- | --- | --- | --- |
| `ado` | `module` | これ | このモジュールの名前 (リポジトリのURL) |
| `mokpose` | なし | 目的 | パッケージの説明 (npmの `description`) |
| `neyo fxma` | `go 1.22` | 必要 | 対応する fxma の版。npmの `engines` と同じ範囲の書き方 (`*`、`>=0.1.0`、`^1.2`) |
| `foryulvu` | なし | 許可する | ライセンス。SPDXのライセンス式をそのまま書く |
| `neyo` | `require` | 必要 | 使うパッケージと版 |
| `unneyo` | `exclude` | 要らない | 使わせない版 |
| `lainei` | `replace` | 〜のように | パッケージを別の場所のもので置き換える |
| `youlugive` | なし (ESLintの設定の `rules`) | 許し | `fixmarge lavi` で指摘しない規則 |

- `foryulvu` の `MIT OR Apache-2.0` は、`OR`・`AND`・`WITH` まで含めてSPDXの1つの表記なので、フィクスマージ語に置き換えません。識別子がSPDXの一覧になければエラーです。
- `foryulvu` と `neyo fxma` は必須です。どちらかがないパッケージは `fixmarge werer` で取り込めません (`Ca-nay`)。
- 手元の fxma が `neyo fxma` の範囲に入っていなければ、`fixmarge werer` の時点でエラーです。
- `fixmarge.sum` は人が書かないファイルなので、語は入れません。`go.sum` と同じく「パッケージ名、版、ハッシュ」を並べます。

### 外部パッケージの取得

| コマンド | 英語の別名 | Go | すること |
| --- | --- | --- | --- |
| `fixmarge ado meif <名前>` | `mod init` | `go mod init` | `fixmarge.mod` を作る |
| `fixmarge werer <パッケージ>@<版>` | `get` | `go get` | パッケージを追加する |
| `fixmarge ado saymal` | `mod tidy` | `go mod tidy` | 使っていないものを消し、足りないものを足す |

- パッケージの名前はGitリポジトリのURLから `https://` を除いたもの (`github.com/owner/repo`) で、版はGitのタグ (`v1.2.0`) です。タグは `v` で始まるsemverだけを版として使います。
- 取ってくるのはfixmargeに組み込んだGitの実装 (go-git) なので、`git` コマンドは要りません。タグが指すコミットのファイルを取ってきます。
- `fixmarge werer <パッケージ>` のように版を省くと、いちばん新しい正式版を使います。正式版がなければいちばん新しいプレリリースです (Goの `latest` と同じ)。
- ダウンロードするのは `fixmarge werer` と `fixmarge ado saymal` のときだけで、`fixmarge keeon` の途中ではダウンロードしません。取ってきていないパッケージを読み込むと `Nii eXi'` です。
- ダウンロードしたものはOSのキャッシュの場所 (Linuxは `~/.cache/fixmarge/mod`、Windowsは `%LocalAppData%\fixmarge\mod`、macOSは `~/Library/Caches/fixmarge/mod`) に置き、プロジェクトの間で使い回します。環境変数 `FIXMARGE_CACHE` を設定すると、その下の `mod` に置きます。
- 同じパッケージを別の版で要求されたら、要求の中でいちばん新しい版を1つだけ使います (Goの最小バージョン選択)。`unneyo` で使わせない版が選ばれたら、その次に新しい版に上げます。
- `fixmarge.mod` の `neyo` には、直接読み込むものだけでなく、間接的に使うものも含めて使うパッケージをすべて並べます (Go 1.17からの `go.mod` と同じ)。`fixmarge keeon` は `fixmarge.mod` だけを見て読み込み先を決めます。
- `fixmarge.sum` のハッシュはGoの `go.sum` と同じ `h1:` の形です (ファイルごとのSHA-256をまとめたもの)。`fixmarge keeon` は読み込む前にキャッシュの中身をこのハッシュと比べ、合わないか `fixmarge.sum` に書かれていなければ `CraeyTA'` です。版を選ぶためにたどった、選ばれなかった版のハッシュも残ります。
- `lainei` で手元のフォルダに置き換えたパッケージは、ダウンロードもハッシュの確認もしません。
- `fixmarge ado saymal` は、パッケージの下書き (`.fxma`) の `werer` を集めて使うパッケージを決め直します。清書 (`.fxmac`)、`.` で始まるフォルダ、別の `fixmarge.mod` があるフォルダは見ません。`neyo` にない読み込み先は、名前の長い方からリポジトリを探して、いちばん新しい版を足します。読み込み先のファイルがなければ、その `werer` の位置で `Nii eXi'` です。
- WebAssembly版のプレイグラウンドはGitを使えないので、標準ライブラリと同じ画面で書いたファイルだけを読み込めます。

外部パッケージの読み込みは次のとおりです。

| 書き方 | 読むファイル | JavaScript |
| --- | --- | --- |
| `"github.com/alice/strs" aste strs werer` | パッケージの直下の `strs.fxma` (名前の最後と同じ) | `import * as strs from "strs"` |
| `"github.com/alice/strs/case" aste c werer` | パッケージの直下の `case.fxma` | `import * as c from "strs/case"` |

- 外部パッケージは、清書を実行しているときも下書き (`.fxma`) を読みます。キャッシュの中身は書き換えないためです。
- 外部パッケージの中の `"./util"` は、そのパッケージの中のファイルを読みます。外部パッケージの中の外部パッケージは、実行しているパッケージの `fixmarge.mod` で決めた版を読みます。

取ってくるときのエラーは次のとおりです。

| エラー | 起きるとき |
| --- | --- |
| `Nii eXi'` | `fixmarge.mod` が見つからない。リポジトリやタグがない |
| `AyAmis` | `fixmarge.mod` の書き方や版の書き方が正しくない。`foryulvu` がSPDXのライセンス式として読めない |
| `Ca-nay` | 取ってきたパッケージに `fixmarge.mod` がない。`ado` が取ってきた名前と違う。`foryulvu` か `neyo fxma` がない。手元の fxma が `neyo fxma` の範囲に入らない |
| `Su·Ready` | `fixmarge ado meif` で、もう `fixmarge.mod` がある。`fixmarge.mod` に同じ項目や同じパッケージを二度書いた |
| `CraeyTA'` | 取ってきたものが `fixmarge.sum` と合わない |

開発中の fxma (タグから作っていない版) では、`neyo fxma` の範囲を確かめません。`fixmarge ado meif` は手元の fxma の版から `neyo fxma >=0.1.0` のように書き、開発中の版なら `neyo fxma *` と書きます。

## 標準パッケージ

標準パッケージは `"fxma:名前"` で読み込みます (`"fxma:meme" aste meme werer`)。各パッケージの仕様は [std/README.md](std/README.md) から、パッケージごとのファイルに分けて書いています。

### 組み込み (読み込み不要)

| 書き方 | 清書 | 辞書の意味 | 中身 | Python |
| --- | --- | --- | --- | --- |
| `nang(xs)` | `nAng` | 長い | 配列・文字列・辞書の長さ | `len(xs)` |
| `sanand(s)` | `sAnAnd` | 砂 | 文字列をUTF-8のバイト列にする。`sanand(s, "hex")`、`sanand(s, "base64")` で16進・base64の文字列から戻す | `s.encode()`、`bytes.fromhex(s)`、`base64.b64decode(s)` |

### 標準パッケージのために足した言語の機能

- **辞書型**: `{"lang": "fxma"}` で作り、`d no "lang"` で取り出し、`d no "lang" izu "x"` で書き換えます。キーがあるかは `"lang" exi d` で調べます。キーは入れた順に並びます。JavaScriptの `Map` やPythonの `dict` と同じです。`d aste k azagara` で回すと、`k` にはキーが入った順に入ります。JavaScriptの `for (const k in obj)` と同じで、値は `d no k` で取り出します。JSONのオブジェクトを受け取るために必要です。
- **バイト列型**: 型の名前は `sAnAnd` (砂) で、1バイトを砂の1粒に見立てます。PythonのbytesやGoの `[]byte` に当たる型です。`madrazySanand`、`sha256`、`sha512`、`hmac` はバイト列を返し、`base64` と `hex` はバイト列か文字列を受け取って文字列を返します。文字列はUTF-8のバイト列として扱います。バイト列を直接書く書き方 (リテラル) はなく、組み込みの `sanand` で文字列から作ります。型の名前で変換するのは、Goの `[]byte(s)` と同じ発想です。文字列に戻すのは `coword no salook` です。
- **関数を値として渡す**: `weho` で定義した関数は、名前で他の関数に渡せます (`varizama no tarip(xs, double)`)。名前のない関数は、`weho` の名前を省いて式の中に書きます。JavaScriptの `function` 式と同じです。`varizama no tarip(xs, weho (n) n * 2 itswa finend)` のように1行でも書けます。`abiryoku` も付けられます。

## エラー一覧

エラー名も辞書の語から選んでいます。世界観を守るため、コンソールの出力には日本語を併記しません。エラー名、行と列の番号、対象の名前だけを出し、意味はこの表とREADMEで調べられるようにします。

| エラー | 意味 | 起きるとき | Pythonでいうと |
| --- | --- | --- | --- |
| `AyAmis` | 過ち | 構文が正しくない。実行中にJSON・日付・数として読めない文字列を読んだ (こちらは `youlugive` で捕まえられる) | `SyntaxError` |
| `Nii eXi'` | 存在しない | 定義されていない変数や関数を使った。配列の範囲の外を読んだ | `NameError` |
| `imJune` | 不純 | 型が合わない。数と文字列の足し算、条件に真偽値以外を書いた | `TypeError` |
| `kYoNAihim` | 虚無 | 0で割った | `ZeroDivisionError` |
| `Bavellabion` | 呪いの言葉があふれる | 整数が64ビットの範囲を超えた。呪い (マクロ) の展開が深すぎる | なし。Pythonの整数はあふれない |
| `blasTOac` | 冒涜的 | 証拠品 (`showence`) に再代入した | なし |
| `kiNbOO` | 禁忌 | `abiryoku` なしの関数で入出力した (ほとんどは実行前に検出) | なし |
| `Co-Jundic` | 矛盾 | どの選択肢でも `shoubek` を満たせなかった | なし |
| `KeI·tion` | 処刑 | 例外が誰にも捕まらなかった | 未処理の例外 |

そのほかのエラーと警告:

| 名前 | 辞書の意味 | 種類 | 起きるとき |
| --- | --- | --- | --- |
| `DowGuil` | 疑惑 | 警告 | 方言の綴りを使った。`fixmarge keeon`・`fixmarge sinruits`・`fixmarge lavi` で出る |
| `LaVI` | 法の悪魔 | 指摘 | `fixmarge lavi` が規則違反を見つけた。規則は「法の悪魔」の節 |
| `Su·Ready` | 既に | エラー | 同じ名前を二重に定義した。関数の中から外の変数に `izu` した。読み込みが循環した |
| `CraeyTA'` | 壊れた | エラー | ダウンロードしたパッケージが `fixmarge.sum` と合わない (改ざん・破損) |
| `SaiVI` | 災い | エラー | インタプリタ自体のバグ (起きてはいけない内部エラー) |

指摘の表示例: `LaVI Was·Da main.fxma:4:1 tmp`

メッセージの例:

```
Nii eXi' main.fxma:3:5 cnt
```

形は `<エラー名> <ファイル>:<行>:<列> <対象>` です。

捕まらなかった実行時エラーは、そのエラー名で出します (`kYoNAihim main.fxma:3:7 //`)。`accuhearts` で投げた例外だけは `KeI·tion` で出し、対象に理由を書きます (`KeI·tion main.fxma:5:1 bad input`)。どちらも終了コード1で終わります。

### コンソールの文言

エラー以外の表示も、辞書の語だけで出します。

| 場面 | 表示 | 辞書の意味 |
| --- | --- | --- |
| ビルド完了 | `FineNd main.fxmac` | 終わり |
| 方言の警告 | `DowGuil lowhel → LowhEll` | 疑惑 |
| テストの結果 | `hung gement fizzbuzz: Nii Guil-me` | 判決、罪なし |
| テストの失敗 | `hung gement fizzbuzz: Guil-me` | 判決、罪 |

例外として、`fixmarge neyar` が出す非公式表記だけは日本語にします。ガイドラインが「分かりやすい位置に明記」を求めているためです。インストーラの画面も日本語です (「インストーラ」の章)。

## 実装方針

実装はGoで行います。同じコードから、CLIの実行ファイルとWebAssembly版の両方を作れるからです。TypeScript、Rust、Pythonも比べましたが、配布のしやすさでGoを選びました。

### コードの書き方

コードは、人が読んで分かる構造にします。コメントを厚くするのではなく、分け方と共通化で読みやすくします。

- **1つの部品に1つの役割**: 字句解析、構文解析、検査、実行、ビルド、LSPを別々のパッケージに分け、依存の向きを一方向にします。
- **元データは1か所**: キーワード、エラー名、標準パッケージの関数名は1つの表にまとめ、色分け定義や補完候補はそこから生成します。
- **コアは入出力を持たない**: 画面やファイルに触る処理は外から渡し、コアは渡されたものだけを使います。テストもしやすくなります。
- **同じことは同じ名前で**: 言語の関数名と同じく、実装の名前も意味が同じならそろえます。
- **凝った書き方をしない**: 短さより、読んだ人が迷わないことを優先します。
- **コメントは「なぜ」だけ**: 何をしているかはコードで分かるようにし、コメントにはそうした理由だけを書きます。

### 構成

1. **正規化**: 方言表で標準語の骨格に寄せる。
2. **字句解析**: キーワード、数値、識別子、文字列、記号に分ける。
3. **構文解析**: 再帰下降法で構文木を作る。
4. **検査**: 実行前に `abiryoku` の付け忘れと、`showence` への再代入を見つける。
5. **実行**: 構文木をたどるインタプリタ。
6. **ビルド**: 字句の並びを保ったまま、辞書の語だけを正式な綴りに置き換え、インタプリタが実行する形である `.fxmac` を出す。置き換えるのは、語の表にある語 (キーワード、組み込み、標準パッケージの名前と関数名、2語の名前とその部品、辞書のキーの語など) を識別子として書いたものすべてで、文脈は見ない。変数に `meme` と名付けても `MeMe` になる。文字列の中のキー (`"lyuzun"`) は書き換えない。

コマンドとオプションは辞書の語で書きます。慣れた人のために、英語の別名も使えます。どちらで書いても動きは同じで、ヘルプや `fixmarge nai` の説明には辞書の語を出します。

| コマンド | 辞書の意味 | 英語の別名 | すること | 他の言語でいうと |
| --- | --- | --- | --- | --- |
| `fixmarge keeon main.fxma` | 続ける | `run` | 下書きをメモリ上でビルドし、その清書を実行する | `go run` |
| `fixmarge keeon main.fxmac` | 続ける | `run` | 清書をそのまま実行する | ビルド済みのファイルを実行 |
| `fixmarge sinruits main.fxma` | 現実 | `build` | 清書 `main.fxmac` を出力する。`werer` で読み込んでいる下書き (`./` か `../` で始まるもの) も、たどってまとめて清書する | `go build`、`tsc` |
| `fixmarge sinruits --karady main.fxma` | 体 | `--compile` | インタプリタと清書と読み込んだパッケージをまとめた単体の実行ファイルを出力する | `deno compile`、`bun build --compile` |
| `fixmarge coatain` | 法廷 | `test` | 名前が `_coatain.fxma` で終わるファイルの `sentky` を実行して判決を出す。ファイルやフォルダを渡せて、省くと今いる場所より下のフォルダまで全部探す (隠しフォルダと `node_modules` は飛ばす) | `go test`、vitest |
| `fixmarge beaukshi main.fxma` | 美しい | `fmt` | 下書きを整形してファイルを書き換える。下の「整形」 | `go fmt`、`deno fmt` |
| `fixmarge beaukshi --curoi main.fxma` | 呪い | `--expand` | 呪い (マクロ) を展開した結果を表示する | `cargo expand` |
| `fixmarge hak` | 話す | `repl` | 対話モード。1行ずつ打って実行する | 引数なしの `python` |
| `fixmarge nai <エラー名>` | なぜ | `explain` | エラーや警告の意味と直し方を表示する | `rustc --explain` |
| `fixmarge nai`、`fixmarge --nai` | なぜ | `help`、`--help` | 使い方を表示する | `go help` |
| `fixmarge persemon main.fxma` | 見世物 | `bench` | 関数の速さを測る | `go test -bench` |
| `fixmarge werer <パッケージ>@<版>` | 迎える | `get` | パッケージを追加する | `go get` |
| `fixmarge ado meif <名前>` | これ、命 | `mod init` | `fixmarge.mod` を作る | `go mod init` |
| `fixmarge ado saymal` | これ、正常 | `mod tidy` | 使っていないものを消し、足りないものを足す | `go mod tidy` |
| `fixmarge neyar`、`fixmarge --neyar` | 年 | `version`、`--version` | 版と、非公式のファンメイド作品であることを表示する | `go version` |
| `fixmarge trench` | 友達 | `lsp` | エディタのための言語サーバーを動かす | `gopls` |
| `fixmarge lavi` | 法の悪魔 | `lint` | 規則違反を指摘する。下の「法の悪魔」 | ESLint、`go vet` |

- 短いオプションは、辞書の語の頭文字と英語の別名の頭文字の両方を使えます。`--karady` は `-k` でも `-c` でも書けます。
- 同じコマンドの中で頭文字がぶつかったら、ぶつかった側には短い形を付けません。たとえば `--nai` は `-n` ですが、`--neyar` には `-n` を付けず、英語の `-v` だけを使えます。`--help` は `-h` です。
- `fixmarge keeon main.fxma a b` のように、ファイルより後ろに書いた引数は、オプションも含めてそのままプログラムに渡します (`dv no "conwa"`)。
- コマンドを省いてファイルだけを渡すと、`keeon` と同じく実行します。そのため清書の1行目には `#!/usr/bin/env fixmarge` を書けます。macOSとLinuxでは、実行権限を付ければ `./main.fxmac` で動きます。

コンソールにはフィクスマージ語しか出さないため、`fixmarge nai` だけは意味を日本語で説明します。

`fixmarge hak` では、1回の入力を実行し終えるたびに、その中の選択肢 (`mekamo`) を確定して出力を書き出します。`praylee` で待つときと同じ考え方です。後の入力で `shoubek` が偽になっても前の入力には戻らず、`Co-Jundic` になります。

`fixmarge hak` の細かい決まりは次のとおりです。JavaScriptのnodeのREPLやChromeのコンソールに合わせています。

- 前の入力で作った変数や関数は、後の入力でもそのまま使えます。関数の中から、後の入力で作った関数も呼べます。
- 入力が式の文で終わっていたら、その値を出します (`1 + 2` なら `3`)。`vanita` は出しません。文字列は配列の中と同じく `"ab"` と引用符で囲み、`redogh` の出力と見分けられるようにします。
- 前の入力で作った名前は、`weho`、`showence`、`werer` で作り直せます。同じ入力の中で二度作ると、ファイルと同じく `Su·Ready` です。`showence` で作った名前に `izu` で代入すると、前の入力のものでも `blasTOac` です。
- 入力を実行し終えたら、ファイルの一番外側と同じく、まだ終わっていない `itsver` の関数や待機が全部終わるまで待ってから次の入力を受け付けます。
- ブロックや括弧が閉じていないなど、入力が途中で終わっていたら `... ` を出して続きの行を待ちます。続きを待っているときに空行を打つと、そこまでで読んでエラーを出します (Pythonと同じ)。
- エラーになった入力は、なかったことにします。その入力で代入した変数は前の値に戻ります。ただし、もう出した出力と、ファイルを書く、入力を読むなどの外の世界とのやり取りは戻りません。
- エラーの場所は `hak:<行>:<列>` です。行は始めてからの通し番号で、なかったことにした入力の行も数えます。
- 入力の終わり (macOSとLinuxでは Ctrl+D、Windowsでは Ctrl+Z の後に Enter) か、`fxma:ieviljio` の `shid` で終わります。

### 整形 (fixmarge beaukshi)

`fixmarge beaukshi` は下書きを決まった形に整えて、ファイルを書き換えます。Goの `go fmt` やDenoの `deno fmt` と同じです。ファイルやフォルダを渡せて、省くと今いる場所より下の `.fxma` を全部整えます (`fixmarge coatain` と同じ探し方)。変わったファイルだけ書き換えて、その名前を出します (`beAukShi main.fxma`)。

- 字下げは空白2つです。ブロックの中を1段深くし、`finend`、`siva`、`youlugive`、`ennote` の行は1段浅く書きます。
- 行の中の空白もそろえます。演算子の前後に空白を1つ、`,` と `:` の後ろに空白を1つ置き、括弧の内側には空白を置きません (`x izu f(a, b) + 1`)。
- 空白を置かないところ: 呼び出しの `(` の前 (`f(x)`、`dv no "murm"(w)`)、値の前の `-` (`xs no -1`)、範囲の `..` と値の間 (`2..4`、`..1`)、複合代入の記号と `izu` の間 (`x +izu 1`)。`weho (n)` や `nii (x)` のように、キーワードの後ろの `(` の前には空白を置きます。
- 続く空行は1つにまとめ、ファイルの頭と終わりの空行と、行末の空白を消します。ファイルは改行1つで終わります。
- コメントはそのまま残します。行末のコメントの前には空白を2つ置きます。
- 括弧の中で改行していれば、改行を残します。次の行の字下げは、開いたままの括弧やブロックを開いた行1つにつき1段です (Prettierと同じ)。
- 文字列の中は書き換えません。
- 整えた結果を読み直して、字句の並びが変わっていないか確かめます。変わっていたら書き換えずに `SaiVI` エラーにします。
- `fixmarge beaukshi --curoi main.fxma` は、呪いを展開した結果を画面に出すだけで、書き換えません。呪いの中で作った名前は、付け替えた名前 (`tmp·1`) のまま出します。

### 法の悪魔 (fixmarge lavi)

`fixmarge lavi` は、エラーにはならないけれど間違いのもとになる書き方を、法の悪魔 (`LaVI`) が指摘するコマンドです。JavaScriptのESLintやGoの `go vet` に当たります。`fixmarge keeon` と `fixmarge sinruits` が出す警告は方言の綴り (`DowGuil`) だけで、ほかの指摘はすべて `fixmarge lavi` にまとめます。

| 規則 | 辞書の意味 | 指摘すること | 似ているもの |
| --- | --- | --- | --- |
| `Was·Da` | 無駄 | 使っていない変数や読み込み | ESLintの `no-unused-vars` |
| `DowGuil` | 疑惑 | 方言の綴り | なし |
| `shi-des` | 死 | 届かないコード。`itswa`・`ketzen`・`accuhearts` の後ろの文 | ESLintの `no-unreachable` |
| `emBasy` | 恥ずかしい | 入出力をしていないのに `abiryoku` が付いている関数。使わない能力をひけらかしている | なし |
| `anxAn` | 不安 | `itsver` の関数を呼んで `praylee` していない。結果を受け取れず、中のエラーも捕まえられない | typescript-eslintの `no-floating-promises` |
| `SuZai` | 怪しい | `eda` の結果を確かめずに添字に使っている。見つからないと `-1` で最後の要素になる | なし |
| `Meh·BeRac` | たぶん | 小数を `senaji` で比べている。`0.1 + 0.2 senaji 0.3` は `ulieso` | なし |
| `JokeDan` | 冗談 | 型が違うので必ず `ulieso` になる比較。`"1" senaji 1` など | TypeScriptの「常に false」の指摘 |
| `sAiMoku` | 沈黙 | 中身のないブロック。例外を黙って握りつぶすのを防ぐ | ESLintの `no-empty` |
| `dEs·lee` | 絶望 | 戻れる選択肢 (`mekamo`) がない場所の `shoubek`。偽になると必ず `Co-Jundic` になる | なし。この言語だけの指摘 |

- 指摘の形は `LaVI <規則> <ファイル>:<行>:<列> <対象>` です (`LaVI anxAn main.fxma:12:3 greet`)。
- ファイルやフォルダを渡せて、省くと今いる場所より下の `.fxma` を全部調べます (`fixmarge coatain` と同じ探し方)。
- 実行せずに構文木を調べるので、確かにそうだと言えるときだけ指摘します。たとえば `emBasy` は、中で呼ぶ関数がどれも入出力をしないと実行前に分かるときだけ指摘します。引数で受け取った関数や、読み込んだファイルの関数を呼んでいれば指摘しません。
- `Was·Da` は、`_` で始まる名前を指摘しません (ESLintと同じ)。わざと使わない変数は `_i` のように書きます。引数、関数の名前、公開した `showence` も指摘しません。
- `shi-des` は、`itswa`・`ketzen`・`accuhearts` に加えて `keip` の後ろの文も指摘します。
- `sAiMoku` は、関数の本体が空なのは指摘しません。何もしない関数をわざと書くことがあるためです (ESLintの `no-empty` と同じ)。
- `anxAn` は、`itsver` の関数を呼んだ結果を文として捨てているときだけ指摘します。
- `SuZai` は、`eda` の結果をそのまま添字にしたときと、`eda` の結果を入れた変数を、同じ関数の中で一度も比べずに添字にしたときに指摘します。
- `Meh·BeRac` は、`senaji` の片側が端数のある小数か、小数になる計算 (`/` や小数との計算) のときに指摘します。
- `JokeDan` は、`senaji` の両側が書いた形から型の分かる値 (数、文字列、真偽値、`vanita`、配列、辞書) で、型が違うときに指摘します。
- `dEs·lee` は、ファイルに `mekamo` が1つもなく、`mekamo` を持つかもしれないファイル (`./` で読み込むもの) も読み込んでいないときだけ指摘します。`sentky` の本体の `shoubek` はテストなので指摘しません。
- 指摘が1つでもあれば、終了コード1で終わります。CIに入れると、指摘が残っているときに落ちます。
- わざとそう書いているときは、行の終わりに `# lavi: mehberac` のように規則の骨格を書くと、その行のその規則だけ止まります。悪魔に「これは見逃して」と頼むという見立てです。
- パッケージ全体で止めるときは、`fixmarge.mod` に `youlugive mehberac suzai` のように許す規則を並べます。

### 巻き戻しの実現方法

最初は「やり直し方式」で作ります。

1. 各 `mekamo` で何番目の候補を選んだかを、番号の列として記録する。
2. `shoubek` が偽になったら、最後の番号を次の候補に進め、プログラムを頭から実行し直す。
3. 出力は、戻れる選択肢が残っている間だけためておき、選択肢がなくなったら書き出す。選択肢がなければすぐ書き出すので、普通のプログラムは他の言語と同じように動く。入力は記録しておき、やり直しのときは記録を使う。新しく1行読むときは、`praylee` と同じくそれまでの選択肢を確定する。プロンプトを画面に出してから入力を待つためで、読んだ後の `shoubek` が偽になっても、読む前の選択肢には戻らない。

実装は単純ですが、候補が多いと実行回数が掛け算で増えます。パズル程度なら問題ありません。遅ければ、後から継続渡し方式 (CPS) に切り替えます。

### 辞書データの扱い

参考にした辞書リポジトリにはライセンスがありません。そのCSVは同梱せず、この言語で使う語だけを自前のキーワード表に書きます。歌詞の文もリポジトリには入れません。

### 対応する環境

使う人のGoのインストールは前提にしません。Goで作ったプログラムはそれ1つで動く実行ファイルになるためです。`go install` はGoを使う人向けの入れ方の1つです。

幅広い環境で動かすため、字句解析・構文解析・実行のコアは、画面やファイルを直接触りません。入出力は、使う側から関数として渡してもらいます。

| 環境 | 作り方 | 時期 |
| --- | --- | --- |
| CLI (Windows・macOS・Linux) | 通常のGoのビルド | 最初から |
| ブラウザ | WebAssembly (`GOOS=js`) | 最初から |
| WASIランタイム (wasmtimeなど) | WebAssembly (`GOOS=wasip1`) | 後で |
| 他のGoプログラムへの組み込み | コアをGoのライブラリとして公開 | 後で |
| C・Pythonなどへの組み込み | `-buildmode=c-shared` で共有ライブラリにする | 後で |
| マイコン | TinyGoでビルド。Gitとファイル操作は使えない | 後で |

メモリの少ない環境のために、巻き戻しの記録の上限を設定で決められるようにします。

### エディタ

エディタの機能は、言語サーバー `fixmarge trench` (LSP、標準入出力で話す) にまとめます。goplsと同じく、CLIと同じ実行ファイルに入れます。

| 機能 | 中身 |
| --- | --- |
| 診断 | 字句・構文・実行前の検査のエラーと、法の悪魔 (`fixmarge lavi`) の指摘を波線で出す。文言はコンソールと同じくフィクスマージ語だけで、位置は範囲で示すので書かない |
| 整形 | `fixmarge beaukshi` と同じ。構文が正しくなければ何もしない |
| 補完 | 語の表のキーワード・組み込み・標準パッケージの名前と関数、ファイルの中の名前。`モジュール no` の後ろでは、そのモジュールの中の名前だけを出す。下書きには骨格、清書には正式な綴りを入れる |
| ホバー | 辞書の語なら正式な綴りと訳、エラー名なら `fixmarge nai` と同じ説明、ファイルの中の名前なら定義の種類 |
| 定義へ移動 | ファイルの中の定義へ飛ぶ。`werer` で読み込んだ名前は、読み込み先のファイルで `oqun` を付けた定義へ飛ぶ。標準パッケージの中身は実行ファイルに組み込まれているので飛び先はない |
| 色分け | セマンティックトークン。字句解析と名前の解決の結果で、変数・引数・関数・モジュールなどを分ける |

名前の解決 (どの名前がどの定義を指すか) は、実行時と同じ規則を `fixmarge lavi` と共有します。

VS Code拡張 (`editors/vscode`) は、色分けの土台にTextMate文法を持ち、ほかは言語サーバーに任せます。TextMate文法は語の表から生成するので、`fixmarge` を入れていなくても色は付きます。下書きと清書は語の書き方が違うので、別の言語 (`fxma` と `fxmac`) と別の文法にします。言語サーバーは `PATH` の `fixmarge` を動かし、設定 `fxma.path` で場所を変えられます。見つからなければ `Nii eXi' fixmarge` と出します。

Neovim・Vimのプラグインは、リポジトリの直下に置き (`ftdetect/`・`ftplugin/`・`indent/`・`syntax/`・`autoload/`・`lsp/`・`doc/`)、リポジトリそのものをプラグインとして入れられるようにします。中身はVS Code拡張と同じく、ファイルタイプ (`fxma` と `fxmac`)、語の表から生成する色分け、`fixmarge beaukshi` と同じ字下げです。言語サーバーは、Neovim 0.11からは `lsp/fxma.lua` を置くので `vim.lsp.enable('fxma')` でつながります。勝手にはつなぎません (nvim-lspconfigと同じ)。Vimには言語サーバーをつなぐ機能が組み込まれていないので、yegappan/lspとvim-lspの設定の例をヘルプ (`:help fxma`) に書きます。

## 配布

CLIはGitHub ReleasesでOSごとの圧縮ファイルを配り、WebAssembly版はnpmとブラウザのプレイグラウンドで配ります。DenoやBunなど、他の言語処理系と同じ形にそろえます。リポジトリは `github.com/otnc/fiXmArgeScript` です。

### CLIの配布物

タグを打つと、GitHub Actions上のGoReleaserが次の圧縮ファイルを作ります。

| OS | CPU | ファイル |
| --- | --- | --- |
| Linux | x86\_64 / arm64 | `fixmarge_<版>_linux_<CPU>.tar.gz` |
| macOS | x86\_64 / arm64 | `fixmarge_<版>_darwin_<CPU>.tar.gz` |
| Windows | x86\_64 / arm64 | `fixmarge_<版>_windows_<CPU>.zip` |

圧縮ファイルの中身は共通です。

- `fixmarge` (Windowsでは `fixmarge.exe`)
- `README.md`: 同梱のREADME。日本語で書き、二次創作ガイドラインで求められた表示を入れる (後の章)
- `LICENSE-MIT` と `LICENSE-APACHE`: この言語のコードのライセンス
- `THIRD_PARTY_LICENSES`: Goのランタイムなど、実行ファイルに含まれる他者のコードのライセンス

一緒に `checksums.txt` (SHA-256) も公開します。

### インストール方法

| 方法 | コマンド | 時期 |
| --- | --- | --- |
| インストーラ (macOS・Linux) | `curl -fsSL https://fixmargescript.oto.im/install.sh \| sh` | 最初から |
| インストーラ (Windows) | Releasesの `fixmarge_<版>_windows_setup.exe` (Inno Setup製) | 最初から |
| 手動 | Releasesから展開してPATHの通った場所に置く | 最初から |
| Go | `go install github.com/otnc/fiXmArgeScript/cmd/fixmarge@latest` | 最初から |
| Homebrew・Scoop | 自前のtap・bucketを用意する | 後で |

版番号はセマンティックバージョニングです。`fixmarge neyar` は版番号と一緒に、非公式のファンメイド作品であることを表示します。

### 版とリリース

版の元はGitのタグ (`v0.2.0`) だけにして、版をファイルに手で書きません。

- **CLI**: Go 1.24からは、ビルドのときにGitのタグとコミットから版が実行ファイルに埋め込まれます。`fixmarge neyar` はそれを読んで、版とコミットとその日時を出します。タグの上なら `v0.2.0`、タグより後なら `v0.2.1-0.<日時>-<コミット>`、手元の変更が残っていれば `+dirty` が付きます。`go install ...@v0.2.0` で入れた場合も同じく `v0.2.0` になります。
- **全部同じ版**: CLI、WebAssembly版 (npm)、VS Code拡張、LSP、Neovim・Vimのプラグインは、1つのタグで同じ版にそろえます。言語の仕様とLSPと拡張は一緒に変わるからです。
- **リリースの手順**: GitHub Actionsの `release` を `workflow_dispatch` で動かし、次の版 (`patch`・`minor`・`major`・`prerelease`、または版を直接書く) と、プレリリースの識別子 (`rc`・`beta` など、空でもよい) を選びます。最新のタグから次の版を決め、テストを通してからタグを打ち、GoReleaserでReleasesを作ります。版の計算はnpmの `semver` と同じです。
  - 識別子なしの `patch`・`minor`・`major` は正式版になる (`v0.1.0` の次の `minor` は `v0.2.0`、`v0.2.0-rc.1` の次の `minor` は `v0.2.0`)。
  - 識別子ありの `patch`・`minor`・`major` は、その版のプレリリースの最初になる (`v0.1.0` の次の `minor` と `rc` は `v0.2.0-rc.0`)。
  - `prerelease` はプレリリースの番号を1つ進める (`v0.2.0-rc.0` の次は `v0.2.0-rc.1`)。識別子を書けば識別子を変えられる (`rc` から `beta` へは戻れない。下の決まりのため)。
  - 版を直接書くときは、識別子を空にして版に書く (`0.2.0-rc.0`)。
  - 次の版は最新のタグより新しくなければならず、プレリリースの識別子は英小文字だけにする (`0.1.1-0` のような数字だけの識別子はdist-tagにできないため)。
- **同梱するライセンス**: `THIRD_PARTY_LICENSES` は、リリースのときにGo本体のライセンスと、`go-licenses` で集めた依存のライセンスから作ります。CLIとWebAssembly版で同じファイルを使います。
- **npmのパッケージ**: リポジトリの `package.json` の `version` は `0.0.0` のままにし、リリースのときだけ書き換えて公開します (コミットはしません)。dist-tagは、プレリリースなら識別子 (`0.2.0-rc.1` なら `rc`)、正式版なら `latest` です。
- **VS Code拡張**: Marketplaceは `0.2.0-rc.1` のようなプレリリースの版を受け付けないので、Marketplaceには正式版だけを出します。VSIXファイル (`fixmarge_<版>_vscode.vsix`) は、どの版でもReleasesに添付します。
- **Neovim・Vimのプラグイン**: 版のファイルは持たず、リポジトリのタグをそのまま使います。lazy.nvimの `version` もNeovimの `vim.pack` も、Gitのタグで版を選びます。
- **インストーラ**: Inno Setupには、ビルドのときに版を渡します (`iscc /DAppVersion=0.2.0`)。`install.sh` はReleasesの最新版を見るので、版を持ちません。

### WebAssembly版

`GOOS=js GOARCH=wasm` で `fixmarge.wasm` を作り、Go付属の `wasm_exec.js` と組み合わせて配ります。

- **Releases**: `fixmarge_<版>_wasm.zip` に `fixmarge.wasm`、`wasm_exec.js`、`README.md`、ライセンス類を入れる。
- **npm**: パッケージ `fxma` として公開し、JavaScriptから `run(source, { stdin, files, args })` と `build(source)` を呼べるようにする。`run` は `{ stdout, stderr, exitCode, files }` を返す。`files` は名前と中身の組で、`werer` でも `fxma:memoid` でも読め、実行した後のファイルも返す。エディタの色分けのために、下書きに書ける名前の一覧を返す `words()` も出す。Node.js (20.16以降) とブラウザの両方で動く。中身はリポジトリの `npm/` にある。
- **プレイグラウンド**: Cloudflare Workersの静的アセットとして https://fixmargescript.oto.im/ に置き、ブラウザだけで書いて実行できるようにする。ViteとCodeMirror 6で作り、リポジトリの `playground/` にある。ファイルを足して `werer` で読み込める。ボタンや欄の名前はコマンドと同じ語 (`keeon`、`sinruits`、`zanfession`、`redogh`) にし、非公式の表記とAI利用は日本語でページの頭に出す。

ブラウザでは対話的な入力ができないため、`zanfession` で読む入力は実行前にまとめて渡します。これは巻き戻しの「入力を記録して使い回す」方式とも合います。実行ファイルが大きすぎる場合は、後でTinyGoでのビルドを検討します。

## インストーラ

インストーラ本体は既存のツールで作り、画面の文言は日本語にします。フィクスマージ語の語彙ではインストーラの文言 (エラー、確認、ディスクの容量など) を表しきれないためです。中身は `installer/` にあります。

### Windows: Inno Setup

- Inno Setupで `fixmarge_<版>_windows_setup.exe` を作り、Releasesで配る。x86\_64とarm64の実行ファイルを1つのセットアップに入れ、インストールするときにCPUに合うほうを選ぶ。
- 文言はInno Setup付属の日本語の言語ファイル (`Japanese.isl`) を使う。
- 既定では管理者権限なしで自分だけに入れる (`%LOCALAPPDATA%\Programs\fixmarge`)。起動したときに、全員向け (`Program Files`、管理者権限が要る) も選べる。入れた場所を、自分だけならユーザーの、全員向けならシステムの `PATH` に足し、アンインストールで取り除く。
- Inno Setupは非営利なら無料で使える。6.5.0から営利利用には有償ライセンスを求めているが、この作品は非営利なので対象外。
- 署名していない `.exe` は、起動時にWindowsのSmartScreenの警告が出る。最初の版では署名せず、警告の出方と実行の手順をREADMEに書く。使う人が増えたら、オープンソース向けの無償の署名サービスに申し込む。
- ライセンス同意の画面には、READMEの非公式表記を出す。ガイドラインの「分かりやすい位置に明記」を満たすためです。
- リリースでは、GoReleaserが作ったWindows向けのzipからセットアップを作り、Releasesと `checksums.txt` に足す。

### macOS・Linux: install.sh

`curl -fsSL https://fixmargescript.oto.im/install.sh | sh` で入れます。`sh -s v0.2.0` のように版を渡すと、その版を入れます。渡さなければ最新の正式版です。

- OSとCPUを調べ、圧縮ファイルと `checksums.txt` を取得してSHA-256を照合し、`~/.fixmarge/bin/fixmarge` に置く。READMEとライセンス類は `~/.fixmarge` に置く。場所は環境変数 `FIXMARGE_INSTALL` で変えられる。
- 使っているシェル (bash・zsh・fish) の設定ファイルに、`PATH` に足す行を書く (Bunと同じ)。すでに `PATH` に入っていれば書かない。macOSのbashは `.bash_profile`、Linuxのbashは `.bashrc`。ほかのシェルでは、足す行を表示するだけにする。
- `install.sh` はプレイグラウンドと同じ場所 (Cloudflare Workersの静的アセット) に置く。
- HomebrewやmacOSの `.pkg` は後で検討する。

## 二次創作ガイドラインと法的な対応

この言語は、公式の配信・二次創作ガイドライン (2026年10月21日更新版) の「9. 二次創作ゲーム・アプリ」に当たるものとして扱います。公開のたびに、最新のガイドラインを確認し直します。

### ガイドラインの条件と対応

| 条件 | 項目 | この言語での対応 |
| --- | --- | --- |
| 完全無料で提供する | 9 | CLI、WebAssembly版、プレイグラウンドをすべて無料にする |
| 広告、寄付、投げ銭などで収益を得ない | 9 | GitHub Sponsorsなどの寄付リンクと、サイトの広告を置かない |
| 公式のイラスト、ロゴ、楽曲、シナリオなどを直接組み込まない | 7・9 | ロゴや画像は使わない。歌詞の文は入れず、キーワードには単語だけを使う |
| 素材は自作か、権利処理済みのものだけ | 9 | アイコンなどは自作する。ライセンスのない辞書CSVは同梱しない |
| 「非公式のファンメイド作品」と配信ページと作品内に明記する | 9 | 配信ページ: リポジトリREADME、Releases、npm、プレイグラウンド。作品内: `fixmarge neyar`、インストーラの同意画面、同梱README |
| 公式の続編や公認作品と誤認させない | 9・10 | 名前や説明に「公式」「公認」を使わない |
| 生成AIを使った作品は「AI利用」と明記する | 6 | 制作にClaudeを使うため、READMEに「AI利用」と書く |
| ゲームのデータを抽出して再配布しない | 7 | ゲームのファイルは解析しない |
| 重大なネタバレを避けにくい場所に書かない | 1 | README、サンプル、エラー文に犯人や結末を書かない |
| 特定の思想の主張や中傷に使わない | 3 | 該当する内容を含めない |
| 発信には用途別のハッシュタグを使う | 4 | 告知には `#まのさば二次創作` を使う |

### そのほかの法的な対応

- **コードのライセンス**: `MIT OR Apache-2.0` の二重ライセンスにする。Rustと同じで、使う側がどちらかを選べる。ただし、作品の権利まで許諾するものではないとREADMEに書く。
- **Goのライセンス**: 実行ファイルと `wasm_exec.js` にはGoのコード (BSD 3条項ライセンス) が含まれる。その著作権表示を `THIRD_PARTY_LICENSES` に入れて同梱する。
- **依存ライブラリ**: 自前で書くより保守しやすい場合は、Goのライブラリ (標準・外部とも) を使う。実行ファイルに入るもの (今は urfave/cli (MIT)、spf13/afero (Apache-2.0)、golang.org/x/text (BSD 3条項)) は、同じファイルにライセンスを追記する。

### 同梱README (README.md) に書くこと

ガイドラインが日本語なので、文書はすべて日本語で書きます。リポジトリと配布物は同じ `README.md` を使い、仕様書 (`RFC.md`) も含めて英語版は作りません。

1. 「魔法少女ノ魔女裁判」の非公式のファンメイド作品であること。Acaciaおよび合同会社Re,AERとは関係がなく、公認・監修を受けていないこと
2. 作品の権利は権利者にあること
3. 完全無料で、収益を得ていないこと
4. 「AI利用」 (生成AIを使って制作したこと)
5. 準拠したガイドラインのURLと最終更新日
6. キーワードに使った語の出典 (この仕様書の出典と同じもの)
7. ライセンスと、同梱した他者のライセンスの場所
8. この作品への問い合わせ先 (公式の窓口ではないことを添える)

キーワードに作品の単語を使うことが「公式素材の組み込み」に当たるかは、ガイドラインには書かれていません。公式は個別の事前確認には原則応じないため、文章ではなく単語だけを使うという範囲にとどめます。この章は法律の助言ではありません。

## 未決事項

- [ ] 使うキーワードを、設定資料集の翻訳表と照らし合わせる
- [ ] 構文そのものを組み立てる本格的なマクロ。最初の版は型はめ式だけにして、後の版で考える
- [ ] `no` の語順を、歌詞ブックレットなどの公式の資料で確かめる
- [ ] 公開前に、fiXmArge Script という名前が他のソフトや、GitHubの fiXmArge アカウント (ファンの辞書) と紛らわしくないか確認する

## 出典

### フィクスマージ語

公式の一次資料は、クラウドファンディング特典の設定資料集 (単語と日本語訳) と、サウンドトラックの歌詞ブックレットです。下の資料はどれも二次資料です。

二次創作の条件は [「魔法少女ノ魔女裁判」配信・二次創作ガイドライン](https://manosaba.com/guidelines) (2026年10月21日更新版) に従います。

| 資料 | 種類 | 使ったところ |
| --- | --- | --- |
| [Asankilp/fiXmArge-dictionary](https://github.com/Asankilp/fiXmArge-dictionary) | ファン辞書 (CSV、264語。設定資料集の表を元にしているようです) | キーワードとエラー名のすべて。骨格がぶつからないことの確認 |
| [bilibili 考察記事1](https://www.bilibili.com/opus/1119453646289895475) | 上の辞書の参考資料 | 未確認 |
| [bilibili 考察記事2](https://www.bilibili.com/opus/1116880763320533001) | 上の辞書の参考資料 | 未確認 |
| [魔法少女ノ魔女裁判 攻略・考察Wiki: 考察](https://w.atwiki.jp/manosaba/pages/21.html) | 有志の解読 | 造語のパターン (英語＋日本語)、初期の単語表 |
| [魔法少女ノ魔女裁判 攻略・考察Wiki: 小ネタ](https://w.atwiki.jp/manosaba/pages/22.html) | 有志の解読 | 審問カットインの語 (検索結果の抜粋のみ) |
| [特別企画: フィクスマージ語を解読しよう (たくんまのゆめにっき)](https://scrapbox.io/takunnma-box/%E7%89%B9%E5%88%A5%E4%BC%81%E7%94%BB:%E3%83%95%E3%82%A3%E3%82%AF%E3%82%B9%E3%83%9E%E3%83%BC%E3%82%B8%E8%AA%9E%E3%82%92%E8%A7%A3%E8%AA%AD%E3%81%97%E3%82%88%E3%81%86) | 有志の解読 | 一人称・二人称、`LowhEl` などの表記ブレ |
| [フィクスマージ語 (ピクシブ百科事典、英語版)](https://dic.pixiv.net/en/a/The%20fictional%20language%20fiXmArge) | 百科事典 | 一次資料 (設定資料集と歌詞ブックレット) の所在、研究書の一覧 |
| [フィクスマージ語 (ピクシブ百科事典)](https://dic.pixiv.net/a/%E3%83%95%E3%82%A3%E3%82%AF%E3%82%B9%E3%83%9E%E3%83%BC%E3%82%B8%E8%AA%9E) | 百科事典 | 未確認 (ページを開けなかった) |
| [【正誤表】宝生マーゴと学ぶフィクスマージ語の世界](https://savakan-fish.fanbox.cc/posts/11279936) | 同人誌の正誤表 | `Iz'u` と `N'o` の語順の解釈 |
| [svea氏の解読ポスト (X)](https://x.com/sveanium_248945/status/1786991851477856500) | 有志の解読 | `Iz'u` と `N'o`、`·` を含む語の整理 |
| [SLAVE.V-V-R氏のポスト: 方言の分岐 (X)](https://x.com/SLAVE_V_V_R/status/1785483559513575718) | 作詞者の発言 | 方言を仕様として受け入れる根拠 |
| [SLAVE.V-V-R氏のポスト: 新しい語 (X)](https://x.com/SLAVE_V_V_R/status/1880818427272773859) | 作詞者の発言 | 新しい語がYouTubeのコメントで公開されていること |
| [魔法少女ノ魔女裁判 (Wikipedia)](https://ja.wikipedia.org/wiki/%E9%AD%94%E6%B3%95%E5%B0%91%E5%A5%B3%E3%83%8E%E9%AD%94%E5%A5%B3%E8%A3%81%E5%88%A4) | 百科事典 | 主題歌のために作られた言語であること |
| [ニコニコ公式のポスト (X)](https://x.com/nico_nico_info/status/2075075628395876450) | ニュース | 設定資料集に翻訳表があるという返信 |
| [かろてら氏の単語帳 (X)](https://x.com/karoterra2008/status/1956686114452115474) | ファンの単語帳 | 未確認 (画像を開けなかった) |
| [フィクスマージ語研究ノート第2版 (ISDN)](https://isdn.jp/2784874501013) | 同人誌の書誌情報 | 語源研究の資料があること |
| [fiXmArge/fiXmArge-language-dict](https://github.com/fiXmArge/fiXmArge-language-dict) | ファン辞書 | 中身はまだREADMEだけ |
| [Magical Girl's Witch Trial/Settings (NamuWiki)](https://en.namu.wiki/w/%EB%A7%88%EB%B2%95%EC%86%8C%EB%85%80%EC%9D%98%20%EB%A7%88%EB%85%80%EC%9E%AC%ED%8C%90/%EC%84%A4%EC%A0%95) | 百科事典 | 語の作られ方の説明 (検索結果の抜粋のみ) |

### ゲームの仕組みと物語

- [魔法少女ノ魔女裁判 (Steamストアページ)](https://store.steampowered.com/app/3101040/Magical_Girl_Witch_Trials?l=japanese): 矛盾を暴いて魔女を探す議論パート
- [gamebiz: Acacia設立とゲーム説明](https://gamebiz.jp/news/384172): ADVパートと魔女裁判パート、固有の魔法
- [game8: 第1話の攻略チャート](https://game8.jp/manosaba/706345): 捜査パートでの証拠品集め
- [game8: バッドエンド一覧](https://game8.jp/manosaba/707055): バッドエンドから直前の選択肢に戻れること
- [電ファミ: 発表記事](https://news.denfaminicogamer.jp/news/240404k): 13人の魔女候補と牢屋敷の設定
- [電ファミ: Switch版発売](https://news.denfaminicogamer.jp/news/2607093v): 絶海の孤島の牢屋敷という舞台
