# 非同期処理

非同期の関数は `itsver` (いつか) を付けて定義し、終わるのを待つときは `praylee` (祈る) を文の後ろに置きます。JavaScript の `async` と `await` と同じ仕組みです。

| 書き方 | 清書 | 辞書の意味 | JavaScript | Python |
| --- | --- | --- | --- | --- |
| `itsver weho f()` | `itsver` | いつか、遅かれ早かれ | `async function f()` | `async def f()` |
| `f() praylee` | `pray·lee` | 祈る | `await f()` | `await f()` |
| `[f(), g()] praylee` | `pray·lee` | 祈る | `await Promise.all([f(), g()])` | `await asyncio.gather(f(), g())` |

- `itsver` を付けた関数を呼ぶと、すぐに「いつか終わるもの」(JavaScript の Promise) が返ります。`praylee` で終わるのを待つと、結果が得られます。
- 配列に `praylee` を付けると、全部が終わるのを待ちます。
- `praylee` を使えるのは、`itsver` の関数の中と、ファイルのトップレベルです。それ以外で使うと `Ca-nay` エラーです。
- 実行は JavaScript と同じくシングルスレッドで、`praylee` のところでだけほかの処理に切り替わります。切り替わる順番を決まったものにしているので、バックトラックでやり直しても同じ結果になります。
- `praylee` で待つと、それまでの選択肢 (`mekamo`) は確定します。待っている間に外の状態が変わるためです。

細かい動きは Node.js に合わせています。

- `itsver` の関数の本体は、呼ぶとすぐに動き始めます。最初の `praylee` まではその場で実行し、そこで呼んだ側に戻ります (`async function` と同じ)。
- 終わっているものや、`1` のような Promise でない値に `praylee` を付けても、一度ほかの処理に順番を回します (`await 1` と同じ)。
- 待っている処理は、待ち始めた順に再開します。`toim no saimoku` は、終わる時刻の早い順、同じなら呼んだ順に終わります (`setTimeout` と同じ)。
- トップレベルを最後まで実行しても、終わっていない `itsver` の関数や待機、通信が残っていれば、全部終わるまで待ってから終了します。
- `praylee` されないまま、捕まえていないエラーで終わった `itsver` の関数があれば、全部終わった後でそのエラーを出して止まります (unhandled rejection と同じ)。後で `praylee` して `youlugive` で捕まえれば止まりません。
- `praylee` で待っていて、ほかに動ける処理も外から届くものもなく、二度と再開できないとき (自分自身の結果を待つなど) は、トップレベルが待っている `praylee` の場所で `Co-Jundic` エラーになります。
- 外とのやり取り (待機、通信、魔女の終了) が終わった順番は記録しておき、バックトラックでやり直すときは記録どおりの順番で再開します。
- 何度も知らせが届くもの (サーバー、ファイルの監視、シグナル、魔女の出力) には、`nozuki` (気付き) に関数を渡します。知らせが届くたびに、その関数を別の処理として呼びます。その関数の中で捕まえていないエラーが起きると、Node.js と同じくプログラムが止まります。
- サーバーとファイルの監視は、動いている間はプログラムを終わらせません。シグナルを待っているだけなら終わります (Node.js と同じ)。

fxma:

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

標準パッケージで非同期になる関数は、待機 (`toim no saimoku`)、魔女 (子プロセス) の終了待ち、通信 (`gemon no fineclowt`) です。
