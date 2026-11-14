<h1>
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/logo-dark.svg" />
    <img src="assets/logo-light.svg" alt="fiXmArge Script" width="560" />
  </picture>
</h1>

『魔法少女ノ魔女裁判』の架空言語「フィクスマージ語」の単語をキーワードにした、**非公式のファンメイド**のプログラミング言語です。略称は fxma です。

> [!IMPORTANT]
> この作品は「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。作品の権利は権利者にあります。
>   
> 完全無料で公開しており、収益は得ていません。

- **AI利用**: この言語の実装は、生成AIを使って制作しています。
- 準拠しているガイドライン: [「魔法少女ノ魔女裁判」配信・二次創作ガイドライン](https://manosaba.com/guidelines) (2026年10月21日更新版)
- キーワードに使った語の出典は [docs/RFC.md](docs/RFC.md) の「出典」にあります。

## 状態

開発中です。仕様は [docs/RFC.md](docs/RFC.md)、標準パッケージは [docs/std/](docs/std/README.md) にあります。

## インストール

macOS と Linux では、次のコマンドで `~/.fixmarge/bin` に入ります。シェルの設定ファイルに `PATH` を足す行も書きます。

```
curl -fsSL https://fixmargescript.oto.im/install.sh | sh
```

Windows では、[Releases](https://github.com/otnc/fiXmArgeScript/releases) の `fixmarge_<版>_windows_setup.exe` を実行します。セットアップには署名がないので、初めて開くと「Windows によって PC が保護されました」(SmartScreen) と出ます。「詳細情報」を押し、「実行」を押すと進めます。

どの OS でも、Releases の圧縮ファイルを展開して `PATH` の通った場所に置くか、Go があれば `go install github.com/otnc/fiXmArgeScript/cmd/fixmarge@latest` でも入れられます。

## 使い方

```
fixmarge keeon examples/fizzbuzz.fxma
fixmarge sinruits examples/fizzbuzz.fxma
```

`keeon` は下書き (`.fxma`) をそのまま実行し、`sinruits` は清書 (`.fxmac`) を書き出します。`fixmarge hak` で対話モード、`fixmarge nai` でコマンドの一覧が出ます。

外部パッケージは Git のタグから取ってきます。`fixmarge ado meif github.com/you/app` で `fixmarge.mod` を作り、`fixmarge werer github.com/alice/strs@1.2` で足し、`fixmarge unneyo github.com/alice/strs` で外します。取ってきたものは npm の `node_modules` と同じくプロジェクトの `.fixmarge` フォルダに置きます。`.fixmarge` はコミットしません (`fixmarge ado meif` が `.gitignore` に足します)。コミットするのは `fixmarge.mod` と `fixmarge.sum` で、取ってきた人は引数なしの `fixmarge werer` でそろえます。版は npm と同じく `fixmarge.mod` と `fixmarge.sum` で決め、同じパッケージの違う版をいくつも入れて、`"github.com/alice/strs@1" aste strs werer` のように読み込むときに選べます。詳しくは [docs/RFC.md](docs/RFC.md) の「パッケージとモジュール」にあります。

WebAssembly 版は [npm/](npm/README.md) (npm のパッケージ `fixmarge`) と [playground/](playground/) (ブラウザのプレイグラウンド) にあります。

VS Code 拡張は [editors/vscode/](editors/vscode/README.md) にあります。色分け、診断、補完、整形などは、言語サーバー `fixmarge trench` が行います。

Vim と Neovim では、このリポジトリをそのままプラグインとして入れられます (lazy.nvim なら `{ 'otnc/fiXmArgeScript' }`)。Neovim 0.11 からは `vim.lsp.enable('fixmarge')` で言語サーバーにつながります。詳しくは [doc/fixmarge.txt](doc/fixmarge.txt) (`:help fixmarge`) にあります。

## 開発

[Go](https://go.dev/) と [pnpm](https://pnpm.io/) (と Node.js) を使います。開発のコマンドは、ルートの [package.json](package.json) の `scripts` にまとめてあります。最初に `pnpm install` を実行します。

| コマンド | すること |
| --- | --- |
| `pnpm fixmarge keeon examples/fizzbuzz.fxma` | 手元のソースの fixmarge を動かす。`fixmarge` の後ろはそのまま渡る |
| `pnpm install:cli` | 手元のソースの fixmarge を `go install` で入れ、`fixmarge` だけで呼べるようにする |
| `pnpm check` | 整形の確認、`go vet`、テスト、サンプルの検査 (`lavi` と `coatain`) をまとめて行う |
| `pnpm test` | Go のテストだけ |
| `pnpm test:wasm` | WebAssembly 版を作り、npm のパッケージをテストする |
| `pnpm test:dict` | キーワードの綴りを辞書と照らし合わせる (辞書を GitHub から取ってくる) |
| `pnpm fmt` | Go のコードを整形する |
| `pnpm dev` | プレイグラウンドを手元で開く |
| `pnpm build` | fixmarge を `dist/` に作る |
| `pnpm build:wasm` | WebAssembly 版 (`npm/fixmarge.wasm`) を作る |
| `pnpm build:playground` | プレイグラウンドを Cloudflare に送る形まで作る |
| `pnpm build:vscode` | VS Code 拡張の VSIX を作る |
| `pnpm build:assets` | ロゴとアイコンを作り直す |
| `pnpm update:grammar` | 語の表を変えた後に、エディタの色分けのファイルを作り直す |
| `pnpm setup:env` | `.env.example` から `.env` を作る (リリースとデプロイに使うトークンを書く) |

`main` に push すると、プレイグラウンドは GitHub Actions が出し直します。リリースは Actions の release を手で動かします。

## ライセンス

コードは [MIT](LICENSE-MIT) と [Apache-2.0](LICENSE-APACHE) のどちらかを選べる二重ライセンスです。ただし、作品の権利まで許諾するものではありません。

ロゴとアイコンは [assets/build.mjs](assets/build.mjs) で図形から描いたもので、公式の素材は使っていません。ロゴの文字は [Almendra Display](https://fonts.google.com/specimen/Almendra+Display) (SIL Open Font License 1.1) の字形を線に変えたものです。作り直すときは `pnpm build:assets` を実行します (fixmarge.exe に埋め込むアイコンを作るのに Go も使います)。

## 問い合わせ

GitHub の Issue にお願いします。
