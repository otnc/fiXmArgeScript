# 言語の特徴

fxma の目立つ機能は、元のゲームの仕組みから取っています。ゲームは、日常のパート、証拠品を集める捜査、発言の矛盾を暴く魔女裁判という流れで進みます。

| ゲームの仕組み | fxma の機能 | 普通の言語でいうと |
| --- | --- | --- |
| バッドエンドになっても直前の選択肢に戻れる | 選択肢 (`mekamo`) と条件 (`shoubek`) | バックトラック (`amb` 演算子) |
| 捜査で集める証拠品 | 証拠 (`showence`) | 定数 |
| 登場人物の固有の魔法 | 能力 (`abiryoku`) | 副作用の宣言 |
| 魔女裁判と処刑 | 裁き (`sentky`)、告発 (`accuhearts`)、処刑 | テストと例外 |

## バックトラック

fxma でいちばん特徴のある機能です。

1. `a owk b owk c mekamo` は、候補のうち最初の値を返し、その場所を「選択肢」として覚えます。`mekamo` で終わる文の中では、`owk` は「または」ではなく候補の区切りです。
2. その後の `条件 shoubek` が偽になると、直前の選択肢まで巻き戻り、次の候補で続きを実行します。
3. 巻き戻すと、変数の値もその選択肢の時点に戻ります。失敗した道筋で出力したものは捨てます。
4. 候補を使い切ったら、もう1つ前の選択肢に戻ります。戻れる選択肢がないときは `Co-Jundic` (矛盾) で終了します。

例として、証言から魔女を当てるプログラムを書きます。

fxma:

```
# 容疑者は1〜3番。魔女は1人
witch izu 1 owk 2 owk 3 mekamo
# 証言A：魔女は1番ではない
nii witch senaji 1 shoubek
# 証言B：魔女は奇数番
witch % 2 senaji 1 shoubek
witch redogh
```

JavaScript:

```javascript
const witch = [1, 2, 3].find((w) => w !== 1 && w % 2 === 1);
if (witch === undefined) {
  throw new Error("Co-Jundic");
}
console.log(witch);
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

どれも `3` を出力します。fxma 版は「魔女は1か2か3かも。1ではないべき。奇数であるべき。魔女を伝える」と読めます。JavaScript や Python では絞り込みやループを自分で書きますが、fxma では候補と条件を並べるだけです。

## 定数

`showence` で作った変数には再代入できません。再代入しようとすると、実行前に `blasTOac` (冒涜的) エラーになります。

```
showence limit izu 15
limit izu 20  # → blasTOac
```

JavaScript の `const limit = 15;` と同じです。Python の `limit: Final = 15` にも近いですが、Python は実行時には止めません。

## 副作用の宣言

入出力 (`redogh`、`zanfession` など) をする関数は、`abiryoku` を付けて定義します。`abiryoku` 付きの関数を呼ぶ関数にも付けます。付け忘れると、実行前に `kiNbOO` (禁忌) エラーになります。

```
abiryoku weho greet(name)
  "weRer, " + name redogh
finend
```

JavaScript や Python には対応する仕組みがありません。Haskell の `IO` 型を小さくしたものです。どこで入出力が起きるかが分かっていると、バックトラックで出力を捨てる処理も作りやすくなります。

関数を値として渡すときは、渡した関数で決まります。`abiryoku` 付きの関数を引数に渡した呼び出しは、それ自体が入出力をする呼び出しとして扱います。

- `varizama no tarip(xs, double)` は、`double` に `abiryoku` がなければ、`abiryoku` なしの関数の中でも書けます。
- `varizama no tarip(xs, greet)` は、`greet` が `abiryoku` 付きなので、`abiryoku` 付きの関数の中でしか書けません。
- `werer` で読み込んだ関数も、実行前に調べます。`"./io" no greet werer` の `greet()` も、`"./io" aste io werer` の `io no greet()` も、標準パッケージの `memoid no fineclowt()` も同じです。
- 変数に入った関数や外部パッケージの関数のように、実行前に分からないものは呼んだときに調べます。`abiryoku` なしの関数の中で呼んでいれば `kiNbOO` エラーです。

## テスト

`sentky` ブロックの中の `shoubek` は巻き戻しをせず、テストの検査として扱います。全部通れば無罪、1つでも落ちれば有罪です。

fxma:

```
sentky "fizzbuzz"
  fizzbuzz(3) senaji "Fizz" shoubek
  fizzbuzz(15) senaji "FizzBuzz" shoubek
finend
```

JavaScript (Jest):

```javascript
test("fizzbuzz", () => {
  expect(fizzbuzz(3)).toBe("Fizz");
  expect(fizzbuzz(15)).toBe("FizzBuzz");
});
```

Python (pytest):

```python
def test_fizzbuzz():
    assert fizzbuzz(3) == "Fizz"
    assert fizzbuzz(15) == "FizzBuzz"
```

結果は、`hung gement` (判決) とテストの名前に続けて、無罪なら `Nii Guil-me`、有罪なら `Guil-me` と出します。有罪のときは、次の行にインデントして落ちた場所を出します。

```
hung gement double: Nii Guil-me
hung gement fizzbuzz: Guil-me
  ShouBeK util_coatain.fxma:4:28
hung gement parse: Guil-me
  imJune util.fxma:2:5 +
```

- `shoubek` が1つ落ちたら、Jest や Python の `assert` と同じく、その `sentky` はそこで止めて有罪にし、次の `sentky` に進みます。
- テストの検査として扱う `shoubek` は、`sentky` の本体に直接書いたものだけです。本体から呼んだ関数の中の `shoubek` は、いつもどおりバックトラックに使います。
- 落ちた `shoubek` は `youlugive` では捕まえられません。例外や実行時エラーが捕まらずに `sentky` の外に出たときも有罪で、そのエラーを出します。
- 有罪が1つでもあれば、`fixmarge coatain` は終了コード1で終わります。

テストは、名前が `_coatain.fxma` で終わるファイルに書きます。Go の `_test.go` と同じです。`Coa Tain` (法廷) は `sentky` (裁き) を開く場所なので、この名前にしました。テストする側のファイルは `werer` で読み込みます。`sentky` を書けるのはこのファイルだけで、`fixmarge keeon` はこのファイルを読み込みません。

```
# util_coatain.fxma
"./util" aste util werer

sentky "double"
  util no double(2) senaji 4 shoubek
finend
```

## ベンチマーク

処理の速さは、テストと同じファイルに `persemon` (`Per·SemOn`、見世物) のブロックで測ります。Go がベンチマークを `_test.go` に書くのと同じです。

fxma:

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
- 本体に直接書いた `shoubek` は、`sentky` と同じくテストの検査として扱います。落ちたら測るのをやめて `Guil-me` と落ちた場所を出し、`fixmarge persemon` は終了コード1で終わります。
- バックトラックでやり直すときは、前と同じ回数だけ本体を動かします。出す時間も前に測ったものです。

## 例外

例外は `accuhearts` (告発) で投げ、`teitance` (抵抗) と `youlugive` (許し) で捕まえます。JavaScript の `throw`、`try`、`catch` と同じです。どこでも捕まらなかった例外は、`KeI·tion` (処刑) としてプログラムを終了させます。

`youlugive e` と書くと、捕まえた例外が変数 `e` に入ります。例外のメッセージは `e no "lyuzun"` で取り出せます。キーは辞書の `LyuzUn` (理由) から取りました。`ennote` (`En-NOte`、成れの果て) のブロックは、例外が起きても起きなくても最後に必ず実行します (JavaScript や Python の `finally`)。`canosiri` (`CaNoSiri`、罵り) は、`redogh` と同じく文の後ろに置いて、標準エラー出力に書きます。

fxma:

```
teitance
  "bad input" accuhearts
youlugive e
  e no "lyuzun" canosiri
ennote
  "FineNd" redogh
finend
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

実行中に起きたエラーも、`accuhearts` で投げた例外と同じく `youlugive` で捕まえられます。JavaScript や Python と同じです。

- エラーの種類は、`e no "chidd"` (`ChiDd`、血) で取り出します。値は `"kYoNAihim"` のようなエラー名の清書の綴りで、`accuhearts` で投げた例外は `"AccUhEarts"` です。
- 捕まえられないのは、実行前に見つかるエラー (構文の `AyAmis`、`kiNbOO`、`showence` への再代入の `blasTOac`) と、`Co-Jundic`、`SaiVI` です。`Co-Jundic` はバックトラックの行き止まり、`SaiVI` はインタプリタのバグで、どちらもプログラムからは直せないためです。
