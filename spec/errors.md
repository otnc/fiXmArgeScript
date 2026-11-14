# エラー

エラー名も辞書の語から選んでいます。コンソールにはフィクスマージ語だけを出し、日本語の説明は付けません。出すのはエラー名、ファイルと行と列の番号、対象の名前だけです。意味はこの表か `fixmarge nai <エラー名>` で調べます。

| エラー | 辞書の意味 | 起きるとき | JavaScript でいうと | Python でいうと |
| --- | --- | --- | --- | --- |
| `AyAmis` | 過ち | 構文が正しくない。実行中に JSON・日付・数として読めない文字列を読んだ (こちらは `youlugive` で捕まえられる) | `SyntaxError` | `SyntaxError` |
| `Nii eXi'` | 存在しない | 定義されていない変数や関数を使った。配列の範囲の外を読んだ | `ReferenceError` (範囲の外は `undefined`) | `NameError`、`IndexError` |
| `imJune` | 不純 | 型が合わない。数と文字列を足した。条件に真偽値以外を書いた | `TypeError` | `TypeError` |
| `kYoNAihim` | 虚無 | 0で割った | なし (`Infinity` になる) | `ZeroDivisionError` |
| `Bavellabion` | 呪いの言葉があふれる | 上限を超えた ([深さと大きさの上限](limits.md)) | `RangeError` | `RecursionError`、`OverflowError` |
| `blasTOac` | 冒涜的 | 定数 (`showence`) に再代入した | `TypeError` (`const` への代入) | なし |
| `kiNbOO` | 禁忌 | `abiryoku` のない関数で入出力をした (ほとんどは実行前に見つかる) | なし | なし |
| `Co-Jundic` | 矛盾 | どの選択肢でも `shoubek` を満たせなかった | なし | なし |
| `KeI·tion` | 処刑 | 例外がどこでも捕まらなかった | `Uncaught Error` | 処理されなかった例外 |

そのほかのエラーと警告:

| 名前 | 辞書の意味 | 種類 | 起きるとき |
| --- | --- | --- | --- |
| `DowGuil` | 疑惑 | 警告 | 方言の綴りを使った。`fixmarge keeon`・`fixmarge sinruits`・`fixmarge lavi` で出る。`werer` で読み込んだファイル (`./` か `../` で始まるもの) の分も出す。外部パッケージの分は、使う側では直せないので出さない |
| `LaVI` | 法の悪魔 | 指摘 | `fixmarge lavi` が規則違反を見つけた。規則は [リンター](commands.md#リンター-fixmarge-lavi) にある |
| `Su·Ready` | 既に | エラー | 同じ名前を二重に定義した。関数の中から外の変数に `izu` した。読み込みが循環した |
| `CraeyTA'` | 壊れた | エラー | ダウンロードしたパッケージが `fixmarge.sum` と合わない (改ざんや破損) |
| `SaiVI` | 災い | エラー | インタプリタ自体のバグ (起きてはいけない内部エラー) |

メッセージの形は `<エラー名> <ファイル>:<行>:<列> <対象>` です。

```
Nii eXi' main.fxma:3:5 cnt
```

リンターの指摘は `LaVI` と規則の名前で始まります (`LaVI Was·Da main.fxma:4:1 tmp`)。

捕まらなかった実行時エラーは、そのエラー名で出します (`kYoNAihim main.fxma:3:7 //`)。`accuhearts` で投げた例外だけは `KeI·tion` で出し、対象にメッセージを書きます (`KeI·tion main.fxma:5:1 bad input`)。どちらも終了コード1で終わります。

## コンソールの文言

エラー以外の表示も、辞書の語だけで出します。

| 場面 | 表示 | 辞書の意味 |
| --- | --- | --- |
| ビルド完了 | `FineNd main.fxmac` | 終わり |
| 方言の警告 | `DowGuil lowhel → LowhEll` | 疑惑 |
| テストの成功 | `hung gement fizzbuzz: Nii Guil-me` | 判決、罪なし |
| テストの失敗 | `hung gement fizzbuzz: Guil-me` | 判決、罪 |

日本語を出すのは、`fixmarge gun` の非公式の表記と、`fixmarge nai` の説明と、インストーラの画面 ([配布](distribution.md)) だけです。非公式の表記は、ガイドラインが「分かりやすい位置に明記」を求めているので日本語にしています。
