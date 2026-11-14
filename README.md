<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/logo-dark.svg" />
    <img src="assets/logo-light.svg" alt="fiXmArge Script" width="560" />
  </picture>
</h1>

<p align="center">
  <a href="https://fixmargescript.oto.im/">📕 ホームページ</a> ·
  <a href="https://fixmargescript.oto.im/docs/">📖 ドキュメント</a> ·
  <a href="https://fixmargescript.oto.im/playground/">🎮 プレイグラウンド</a>
</p>

『魔法少女ノ魔女裁判』の架空言語「フィクスマージ語」の単語をキーワードにした、**非公式のファンメイド**のプログラミング言語です。略称は fxma です。

> [!IMPORTANT]
> この作品は「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。作品の権利は権利者にあります。
>   
> 完全無料で公開しており、収益は得ていません。

- **AI利用**: この言語の実装は、生成AIを使って制作しています。
- 準拠しているガイドライン: [「魔法少女ノ魔女裁判」配信・二次創作ガイドライン](https://manosaba.com/guidelines) (2026年10月21日更新版)
- キーワードに使った語の出典は [spec/sources.md](spec/sources.md)にあります。

## 状態

開発中です。仕様は [spec/](spec/README.md)、標準パッケージは [spec/std/](spec/std/README.md) にあります。

## インストール

macOS と Linux では、次のコマンドで `~/.fixmarge/bin` に入ります。シェルの設定ファイルに `PATH` を足す行も書きます。

```
curl -fsSL https://fixmargescript.oto.im/install.sh | sh
```

Windows では、[Releases](https://github.com/otnc/fiXmArgeScript/releases) の `fixmarge_<版>_windows_setup.exe` を実行します。セットアップには署名がないので、初めて開くと「Windows によって PC が保護されました」(SmartScreen) と出ます。「詳細情報」を押し、「実行」を押すと進めます。

どの OS でも、Releases の圧縮ファイルを展開して `PATH` の通った場所に置くか、Go があれば `go install github.com/otnc/fiXmArgeScript/cmd/fixmarge@latest` でも入れられます。

新しい版が出たら `fixmarge neyar` で更新できます。`fixmarge neyar 0.2.0` のように版を指定すると、その版に切り替えます。今の版は `fixmarge gun` で出ます。

## 使い方

```
fixmarge keeon examples/fizzbuzz.fxma
fixmarge sinruits examples/fizzbuzz.fxma
```

`keeon` は下書き (`.fxma`) をそのまま実行し、`sinruits` は清書 (`.fxmac`) を書き出します。`fixmarge hak` で対話モード、`fixmarge nai` でコマンドの一覧が出ます。

外部パッケージは Git のタグから取ってきます。`fixmarge ado meif github.com/you/app` で `fixmarge.mod` を作り、`fixmarge werer github.com/alice/strs@1.2` で足し、`fixmarge unneyo github.com/alice/strs` で外します。取ってきたものは npm の `node_modules` と同じくプロジェクトの `.fixmarge` フォルダに置きます。`.fixmarge` はコミットしません (`fixmarge ado meif` が `.gitignore` に足します)。コミットするのは `fixmarge.mod` と `fixmarge.sum` で、取ってきた人は引数なしの `fixmarge werer` でそろえます。版は npm と同じく `fixmarge.mod` と `fixmarge.sum` で決め、同じパッケージの違う版をいくつも入れて、`"github.com/alice/strs@1" aste strs werer` のように読み込むときに選べます。詳しくは [spec/packages.md](spec/packages.md) にあります。

WebAssembly 版は [npm/](npm/README.md) (npm のパッケージ `fixmarge`) と [プレイグラウンド](https://fixmargescript.oto.im/playground/) (ソースは [site/](site/)) にあります。

VS Code 拡張は [editors/vscode/](editors/vscode/README.md) にあります。色分け、診断、補完、整形などは、言語サーバー `fixmarge trench` が行います。

Vim と Neovim では、このリポジトリをそのままプラグインとして入れられます (lazy.nvim なら `{ 'otnc/fiXmArgeScript' }`)。Neovim 0.11 からは `vim.lsp.enable('fixmarge')` で言語サーバーにつながります。詳しくは [doc/fixmarge.txt](doc/fixmarge.txt) (`:help fixmarge`) にあります。

## 開発

手元で作って試す手順と開発のコマンドは [CONTRIBUTING.md](CONTRIBUTING.md) にあります。

## ライセンス

コードは [MIT](LICENSE-MIT) と [Apache-2.0](LICENSE-APACHE) のどちらかを選べる二重ライセンスです。ただし、作品の権利まで許諾するものではありません。

ロゴとアイコンは [assets/build.mjs](assets/build.mjs) で図形から描いたもので、公式の素材は使っていません。ロゴの文字は [Almendra Display](https://fonts.google.com/specimen/Almendra+Display) (SIL Open Font License 1.1) の字形を線に変えたものです。

## 問い合わせ

GitHub の Issue にお願いします。
