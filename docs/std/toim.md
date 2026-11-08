# fxma:toim

時間。辞書の意味は「時間」。標準パッケージの共通の決まりは [README.md](README.md) を見てください。

```
"fxma:toim" aste toim werer
```

| 書き方 | 清書 | 辞書の意味 | 中身 | JavaScript |
| --- | --- | --- | --- | --- |
| `ieywoo()` | `ieyWOO` | 今すぐ | 現在時刻。1970年1月1日 (UTC) からのミリ秒の整数 | `Date.now()` |
| `saimoku(ミリ秒)` | `sAiMoku` | 沈黙 | 待機。`itsver` の関数で、`praylee` で待つ | `setTimeout` |
| `confsoy(t)` | `conf-Soy` | 争い | 時刻を `[年, 月, 日, 時, 分, 秒, ミリ秒, 曜日]` の配列に分ける | `getFullYear()` など |
| `salook(t)` | `SaLook` | 姿、形 | ISO 8601の文字列にする | `toISOString()` |
| `salook(t, 書式)` | `SaLook` | 姿、形 | 書式に合わせた文字列にする | Day.jsの `format` |
| `fineclowt(s)` | `FinecloWt` | 見つけ出す | ISO 8601の文字列を読んで時刻にする | `Date.parse(s)` |

## 時刻とタイムゾーン

- 時刻は `ieywoo()` と同じミリ秒の整数で表します。専用の型は作りません。
- `confsoy`、`salook` は、そのパソコンのタイムゾーンで考えます。JavaScriptの `getHours()` と同じです。最後の引数に `"utc"` を渡すとUTCで考えます (`confsoy(t, "utc")`、`salook(t, 書式, "utc")`)。
- `confsoy` の月は1〜12です。JavaScriptの `getMonth()` は0〜11ですが、間違えやすいので人の数え方に合わせます。曜日はJavaScriptの `getDay()` と同じく、日曜が0、土曜が6です。
- `salook(t, "utc")` のように2つ目の引数が `"utc"` だけなら、書式ではなくタイムゾーンとして読み、UTCのISO 8601の文字列 (`+00:00` 付き) にします。
- `fineclowt` は、JavaScriptの `Date.parse` と同じく、時差のない日時 (`2026-11-07T21:00:00`) はそのパソコンの時刻として、日付だけ (`2026-11-07`) はUTCとして読みます。
- `salook(t)` の文字列には、`2026-11-07T21:00:00.000+09:00` のように時差を付けます。`fineclowt` は時差の付いた文字列も `Z` で終わるUTCの文字列も読めます。読めない文字列は `AyAmis` エラーです。

## 書式

書式の書き方は、JavaScriptのDay.jsと同じです。

| 書き方 | 中身 | 例 |
| --- | --- | --- |
| `YYYY` | 年 | `2026` |
| `MM` / `M` | 月 (2桁 / そのまま) | `01` / `1` |
| `DD` / `D` | 日 | `07` / `7` |
| `HH` / `H` | 時 (24時間) | `09` / `9` |
| `mm` / `m` | 分 | `05` / `5` |
| `ss` / `s` | 秒 | `03` / `3` |
| `SSS` | ミリ秒 | `042` |
| `d` | 曜日の番号 (日曜が0) | `6` |
| `Z` | 時差 | `+09:00` |
| `[文字]` | 角括弧の中はそのまま出す | `[T]` → `T` |

曜日の名前や月の名前を出す書き方は作りません。名前は言語によって違うので、`confsoy` で番号を取り出して自分で選びます。

## 名前

- `confsoy` は `coword no confsoy` (分割) と同じ語です。時刻を部品に分けるからです。
- `salook` と `fineclowt` は、`json` と同じく「文字列にする」と「文字列から読む」の組です。
- `fineclowt` は `json` や `memoid` の「読む」と同じ語です。
