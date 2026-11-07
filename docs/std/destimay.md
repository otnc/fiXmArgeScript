# fxma:destimay

乱数。辞書の意味は「運命」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:destimay" aste destimay werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `senchoose(1..6)` | `Sen-choose` | 選択肢 | 範囲や配列から1つ選ぶ | なし |
| `commity()` | `Commity` | たぶん | 0以上1未満の小数 | `Math.random()` |
| `madrazy(xs)` | `madRazy` | 狂った | かき混ぜた新しい配列を返す | なし |
| `uuid()` | `uuid` | 固有名詞 | UUID。引数なしでv4、`uuid(7)` でv7 | `crypto.randomUUID()` |
| `gomi(n)` | `GoMi` | 神 | 乱数の種（シード）を決める | なし。Pythonの `random.seed(n)` |

- `madrazy` は、配列の決まりに合わせて新しい配列を返し、渡した配列は変えません。
- `gomi` で種を決めると、その後の `senchoose`、`commity`、`madrazy` は毎回同じ順番で同じ値を出します。テスト（`fxma coatain`）やパズルの問題を毎回同じにするときに使います。`uuid` は種の影響を受けません。
- 巻き戻しのやり直しでは、種を決めたかどうかに関係なく、記録しておいた結果を使います。
- `uuid` はGoの標準の `uuid` パッケージ（Go 1.27から）で作ります。
- `senchoose` に空の範囲や配列を渡すと `Nii eXi'` エラーです。
- 種を決めないときは、実行するたびに違う値になります。
- `gomi` は、運命（`destimay`）を決めるのは神、という見立てです。同じ神なら同じ運命になります。
