# 世界観由来の仕様

まのさばは、牢屋敷での日常 (ADVパート) →証拠品を集める捜査→発言の矛盾を暴く魔女裁判、という流れで進みます。この言語では、その要素を次の4つの機能に対応させます。

| ゲームの要素 | 言語の機能 | 普通の言語でいうと |
| --- | --- | --- |
| バッドエンドになっても直前の選択肢に戻れる | 選択肢と証明 | バックトラック (`amb` 演算子) |
| 捜査で集める証拠品 | 証拠 | 定数 |
| 少女たちの固有の魔法 | 能力 | 副作用の宣言 |
| 魔女裁判と処刑 | 裁き、告発、処刑 | テストと例外 |

13人の少女にちなんで「1スコープに変数は13個まで」という案もありましたが、使いにくいため採用しません。

## 選択肢とバッドエンド

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

## 証拠品 (定数)

`showence` で作った変数は再代入できません。証拠を書き換えようとすると `blasTOac` (冒涜的) エラーになります。証拠を書き換えるのは裁判への冒涜、という見立てです。

```
showence limit izu 15
limit izu 20  # → blasTOac
```

Pythonでいうと `limit: Final = 15` に近いですが、Pythonは実行時には止めません。JavaScriptの `const limit = 15;` と同じく、この言語は止めます。

## 魔法 (副作用の宣言)

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
- `werer` で読み込んだ関数も、実行前に調べます。`"./io" no greet werer` の `greet()` も、`"./io" aste io werer` の `io no greet()` も、標準パッケージの `memoid no fineclowt()` も同じです。
- 変数に入った関数や外部パッケージの関数のように実行前に分からないときは、呼んだ時点で調べて、`abiryoku` なしの関数の中なら `kiNbOO` エラーにします。

## 魔女裁判 (テスト) と処刑 (例外)

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

速さを測る処理も、同じファイルに `persemon` (`Per·SemOn`、見世物) のブロックで書きます。Goがベンチマークを `_test.go` に書くのと同じです。

フィクスマージ語 (ソース):

```
persemon "double"
  util no double(21)
finend
```

JavaScript (Vitest):

```javascript
bench("double", () => {
  double(21);
});
```

`fixmarge persemon` は `persemon` だけを、`fixmarge coatain` は `sentky` だけを動かします。結果は、`Per·SemOn` に続けて名前、回数、1回あたりの時間 (ナノ秒) を出します。

```
Per·SemOn double: 2109718 584ns
```

- `go test -bench` と同じく、本体を何回も動かし、かかった時間が1秒に届くまで回数を増やして測り直します。回数は一度に100倍まで、上限は10億回です。
- 本体に直接書いた `shoubek` は `sentky` と同じくテストとして扱います。落ちたら測るのをやめて `Guil-me` と落ちた場所を出し、`fixmarge persemon` は終了コード1で終わります。
- 巻き戻しでやり直すときは、前と同じ回数だけ本体を動かします。出す時間も前のものです。

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
