# fixmarge

<p align="center">
  <a href="https://fixmargescript.oto.im/">📕 ホームページ</a> ·
  <a href="https://fixmargescript.oto.im/docs/">📖 ドキュメント</a> ·
  <a href="https://fixmargescript.oto.im/playground/">🎮 プレイグラウンド</a>
</p>

[fiXmArge Script](https://github.com/otnc/fiXmArgeScript) の WebAssembly 版です。fiXmArge Script は、『魔法少女ノ魔女裁判』の架空言語「フィクスマージ語」の単語をキーワードにした、**非公式のファンメイド**のプログラミング言語です。

> [!IMPORTANT]
> この作品は「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。作品の権利は権利者にあります。
>   
> 完全無料で公開しており、収益は得ていません。

- **AI利用**: この言語の実装は、生成AIを使って制作しています。
- 準拠しているガイドライン: [「魔法少女ノ魔女裁判」配信・二次創作ガイドライン](https://manosaba.com/guidelines) (2026年10月21日更新版)

## 使い方

Node.js 20.16 以降とブラウザで動きます。

```
npm install fixmarge
```

```javascript
import { run, build } from "fixmarge";

const r = await run('name izu zanfession\n"weRer, " + name redogh\n', { stdin: "fxma\n" });
console.log(r.stdout); // weRer, fxma

const b = await build("x izu tadact\n");
console.log(b.output); // x Iz'u TadACt
```

| 関数 | すること |
| --- | --- |
| `run(source, { stdin, files, args, maxSteps })` | 下書きを実行し、`{ stdout, stderr, exitCode, files }` を返す |
| `build(source)` | 下書きを清書にし、`{ output, stderr }` を返す。エラーなら `output` は `null` |
| `words()` | 下書きに書ける名前の一覧を返す。エディタの色分けに使う |

- `stdin` は `zanfession` で読む入力です。ブラウザでは対話的に読めないので、実行の前にまとめて渡します。
- `files` は名前と中身の組で、`werer` で読み込むファイルにも、`:memoid` で読み書きするファイルにもなります。実行した後のファイルは結果の `files` に入ります。
- `args` は `dv no "conwa"` になる引数です。
- `maxSteps` は実行できる文と繰り返しの数の上限です (`fixmarge keeon --laikill` と同じ)。超えると `Bavellabion` エラーで止まります。ブラウザでは実行中のプログラムを外から止められないので、止まらないかもしれないプログラムを動かすときに付けます。省くと上限はありません。
- 魔女 (子プロセス)、シグナル、サーバー、ファイルの見張りは使えず、`Ca-nay` エラーになります。通信 (`gemon no fineclowt`) は、動かしている環境の `fetch` を使います。

言語の仕様は [spec/](https://github.com/otnc/fiXmArgeScript/blob/main/spec/README.md) にあります。

## ライセンス

コードのライセンスは MIT と Apache-2.0 のデュアルライセンス (`MIT OR Apache-2.0`) で、使う人がどちらかを選べます。全文は `LICENSE-MIT` と `LICENSE-APACHE` にあります。このライセンスで許諾するのはこのリポジトリのコードだけで、「魔法少女ノ魔女裁判」の権利は含みません。

`fixmarge.wasm` と `wasm_exec.js` には、Go と依存ライブラリのコードが含まれます。それらのライセンスは `THIRD_PARTY_LICENSES` にあります。
