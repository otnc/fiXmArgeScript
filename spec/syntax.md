# 基本構文

1行に1文を書き、ブロックは `finend` で閉じます。インデントに意味はありません。関数の呼び出しは `f(a, b)` の形です。

キーワードを置く位置は、その語が日本語で文のどこに来る語から作られたかで決めます。

| 位置 | 決め方 | キーワード |
| --- | --- | --- |
| 文の後ろ | 動詞、語尾 (〜かも、〜べき)、「〜ながら」 | `mekamo` `shoubek` `itswa` `azagara` `redogh` `canosiri` `zanfession` `accuhearts` `werer` `praylee` |
| 文の頭 | 文の頭に来る語 (もし、しかし)、ブロックや宣言の見出し | `mof` `siva` `weho` `abiryoku` `itsver` `oqun` `curoi` `sentky` `persemon` `teitance` `youlugive` `ennote` `showence` `nii` |
| 2つのものの間 | 2つをつなぐ語 | `izu` `no` `exi` `senaji` `sond` `owk` `aste` `batte` |
| 1語で文になる | 位置は問題にならない | `keip` `ketzen` `finend` |

`nii` (否定) は、日本語の「ない」なら文の後ろになります。ただし辞書に語源がなく、`mof nii` でunlessを書くため、文の頭に置きます。

## 後ろに置く語のかかる範囲

文の後ろに置く語は2種類で、どちらも自分より左側をまるごと受け取ります。日本語の文末の動詞が文全体にかかるのと同じです。

| 種類 | キーワード | かかる範囲 | 例 |
| --- | --- | --- | --- |
| 式の終わりの語 | `mekamo` `praylee` `zanfession` | 式全体。`izu` があればその右側全体 | `y izu f() praylee` は `y = await f()` |
| 文の終わりの語 | `redogh` `canosiri` `itswa` `accuhearts` `shoubek` `azagara` `werer` | 文全体 | `x + 1 redogh` は `print(x + 1)` |

- 1つの文に、式の終わりの語と文の終わりの語を続けて書けます。`f() praylee redogh` は `print(await f())` です。
- 式の途中で使うときは括弧で囲みます。`(f() praylee) + 1` のように書きます。
- `zanfession` は左側を入力のプロンプトとして出し、読んだ1行を返します。`name izu "名前は?" zanfession` はPythonの `name = input("名前は?")` です。プロンプトがいらなければ `name izu zanfession` と書きます。

## 変数と出力

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

## 変数のスコープ

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

## 条件分岐

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

## 繰り返しと関数 (FizzBuzz)

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

## 配列

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

## 範囲

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

## no で取り出す

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
