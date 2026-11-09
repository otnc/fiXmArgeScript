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

macOS と Linux では、次のコマンドで `~/.fxma/bin` に入ります。シェルの設定ファイルに `PATH` を足す行も書きます。

```
curl -fsSL https://fixmargescript.oto.im/install.sh | sh
```

Windows では、[Releases](https://github.com/otnc/fiXmArgeScript/releases) の `fxma_<版>_windows_setup.exe` を実行します。セットアップには署名がないので、初めて開くと「Windows によって PC が保護されました」(SmartScreen) と出ます。「詳細情報」を押し、「実行」を押すと進めます。

どの OS でも、Releases の圧縮ファイルを展開して `PATH` の通った場所に置くか、Go があれば `go install github.com/otnc/fiXmArgeScript/cmd/fxma@latest` でも入れられます。

## 使い方

```
go run ./cmd/fxma keeon examples/fizzbuzz.fxma
go run ./cmd/fxma sinruits examples/fizzbuzz.fxma
```

WebAssembly 版は [npm/](npm/README.md) (npm のパッケージ `fxma`) と [playground/](playground/) (ブラウザのプレイグラウンド) にあります。

VS Code 拡張は [editors/vscode/](editors/vscode/README.md) にあります。色分け、診断、補完、整形などは、言語サーバー `fxma trench` が行います。

Vim と Neovim では、このリポジトリをそのままプラグインとして入れられます (lazy.nvim なら `{ 'otnc/fiXmArgeScript' }`)。Neovim 0.11 からは `vim.lsp.enable('fxma')` で言語サーバーにつながります。詳しくは [doc/fxma.txt](doc/fxma.txt) (`:help fxma`) にあります。

## ライセンス

コードは [MIT](LICENSE-MIT) と [Apache-2.0](LICENSE-APACHE) のどちらかを選べる二重ライセンスです。ただし、作品の権利まで許諾するものではありません。

ロゴとアイコンは [assets/build.mjs](assets/build.mjs) で図形から描いたもので、公式の素材は使っていません。ロゴの文字は [Almendra Display](https://fonts.google.com/specimen/Almendra+Display) (SIL Open Font License 1.1) の字形を線に変えたものです。作り直すときは `npm --prefix assets install` の後に `npm --prefix assets run build` を実行します (fxma.exe に埋め込むアイコンを作るのに Go も使います)。

## 問い合わせ

GitHub の Issue にお願いします。
