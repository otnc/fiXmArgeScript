# 非同期処理

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
":toim" aste toim werer

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
