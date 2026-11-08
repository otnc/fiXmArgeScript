# fxma

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
npm install fxma
```

```javascript
import { run, build } from "fxma";

const r = await run('name izu zanfession\n"weRer, " + name redogh\n', { stdin: "fxma\n" });
console.log(r.stdout); // weRer, fxma

const b = await build("x izu tadact\n");
console.log(b.output); // x Iz'u TadACt
```

| 関数 | すること |
| --- | --- |
| `run(source, { stdin, files, args })` | 下書きを実行し、`{ stdout, stderr, exitCode, files }` を返す |
| `build(source)` | 下書きを清書にし、`{ output, stderr }` を返す。エラーなら `output` は `null` |
| `words()` | 下書きに書ける名前の一覧を返す。エディタの色分けに使う |

- `stdin` は `zanfession` で読む入力です。ブラウザでは対話的に読めないので、実行の前にまとめて渡します。
- `files` は名前と中身の組で、`werer` で読み込むファイルにも、`fxma:memoid` で読み書きするファイルにもなります。実行した後のファイルは結果の `files` に入ります。
- `args` は `dv no "conwa"` になる引数です。
- 魔女 (子プロセス)、シグナル、サーバー、ファイルの見張りは使えず、`Ca-nay` エラーになります。通信 (`gemon no fineclowt`) は、動かしている環境の `fetch` を使います。

言語の仕様は [docs/RFC.md](https://github.com/otnc/fiXmArgeScript/blob/main/docs/RFC.md) にあります。

## ライセンス

コードは MIT と Apache-2.0 のどちらかを選べる二重ライセンスです。ただし、作品の権利まで許諾するものではありません。`wasm_exec.js` と `fxma.wasm` には Go のコード (BSD 3条項ライセンス) が含まれます。
