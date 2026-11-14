# 開発の手引き

fiXmArge Script を手元で作って試すための手順です。言語の仕様は [spec/](spec/README.md) に章ごとに、標準パッケージは [spec/std/](spec/std/README.md) にあります。

## 用意するもの

- [Go](https://go.dev/) (版は [go.mod](go.mod) の `go` の行)
- [Node.js](https://nodejs.org/) 24 と [pnpm](https://pnpm.io/) (版は [package.json](package.json) の `packageManager`)

[mise](https://mise.jdx.dev/) を使っているなら、`mise install` で Go の版がそろいます。

最初に一度、リポジトリの直下で次を実行します。

```
pnpm install
```

## コマンド

開発のコマンドは、ルートの [package.json](package.json) の `scripts` にまとめてあります。どれもリポジトリの直下で動かします。

| コマンド | すること |
| --- | --- |
| `pnpm fixmarge keeon examples/fizzbuzz.fxma` | 手元のソースの fixmarge を動かす。`fixmarge` の後ろはそのまま渡る |
| `pnpm install:cli` | 手元のソースの fixmarge を `go install` で入れ、`fixmarge` だけで呼べるようにする |
| `pnpm check` | 整形の確認、`go vet`、テスト、サンプルの検査 (`lavi` と `coatain`) をまとめて行う |
| `pnpm test` | Go のテストだけ |
| `pnpm test:wasm` | WebAssembly 版を作り、npm のパッケージをテストする |
| `pnpm test:dict` | キーワードの綴りを辞書と照らし合わせる (辞書を GitHub から取ってくる) |
| `pnpm fmt` | Go のコードを整形する |
| `pnpm dev` | サイト (紹介・ドキュメント・プレイグラウンド) を手元で開く |
| `pnpm build` | fixmarge を `dist/` に作る |
| `pnpm build:wasm` | WebAssembly 版 (`packaging/npm/fixmarge.wasm`) を作る |
| `pnpm build:site` | サイトを Cloudflare に送る形まで作る |
| `pnpm build:vscode` | VS Code 拡張の VSIX を作る |
| `pnpm build:assets` | ロゴとアイコンを作り直す (fixmarge.exe に埋め込むアイコンを作るのに Go も使う) |
| `pnpm update:grammar` | 語の表を変えた後に、エディタの色分けのファイルを作り直す |
| `pnpm setup:env` | `.env.example` から `.env` を作る (リリースとデプロイに使うトークンを書く) |
| `pnpm check:env` | `.env` が `.env.example` とずれていないか確かめる |

変更を送る前に `pnpm check` が通ることを確かめてください。

## リポジトリの形

JavaScript のパッケージは pnpm の workspace にまとめてあります ([pnpm-workspace.yaml](pnpm-workspace.yaml))。

| 場所 | 中身 |
| --- | --- |
| `cmd/fixmarge/` | コマンド (CLI) |
| `cmd/fixmarge-wasm/` | WebAssembly 版の入口 |
| `internal/` | 字句解析・構文解析・インタプリタ・標準パッケージ・言語サーバーなど |
| `site/` | サイト (Astro + Starlight)。紹介 `/`、ドキュメント `/docs/`、プレイグラウンド `/playground/` |
| `editors/vscode/` | VS Code 拡張 |
| `packaging/npm/` | npm のパッケージ `fixmarge` (WebAssembly 版を包む) |
| `packaging/installer/` | Windows のセットアップと `install.sh` |
| `packaging/assets/` | ロゴとアイコンを描くスクリプト |
| `examples/` | サンプル |
| `spec/` | 言語の仕様 (章ごとのファイル。目次は README.md) と標準パッケージ (std/)。サイトの /docs/ にも出す |
| `doc/` | Vim と Neovim のヘルプ (`:help fixmarge`)。Vim はこの名前のフォルダしか読まないので変えられない |
| `scripts/build/` | 配るものを作る手順 (Windows 向けの exe、THIRD_PARTY_LICENSES) |
| `scripts/check/` | CI で確かめる手順 (セットアップを入れて消す、ファイル バージョン) |
| `scripts/release/` | リリースの手順 (次の版を決める、npm・セットアップ・VS Code 拡張を出す) |

VS Code 拡張だけは workspace に入れず、npm で扱います。VSIX を作る vsce が pnpm の `node_modules` を読めないのと、拡張の名前が npm のパッケージと同じ `fixmarge` でぶつかるためです。`pnpm build:vscode` から呼べます。

workflow (`.github/workflows/`) には、1行で書ける手順だけを書きます。何行もある手順は `scripts/` の用途のフォルダ (`build/`・`check/`・`release/`) にファイルとして置き (Linux では `.sh`、Windows では `.ps1`)、workflow からはそれを呼びます。

エディタの色分けのファイル (`editors/vscode/syntaxes/`、`syntax/` など) は語の表から作っています。手で書き換えず、`pnpm update:grammar` で作り直します。ずれていると `pnpm test` が落ちます。

## 公開

- **サイト** (https://fixmargescript.oto.im/) と `install.sh`: `main` に push すると、関係するファイルが変わったときだけ GitHub Actions が出し直します ([.github/workflows/deploy.yml](.github/workflows/deploy.yml))。
- **リリース**: GitHub Actions の release を手で動かします ([.github/workflows/release.yml](.github/workflows/release.yml))。次の版 (`patch`・`minor`・`major`・`prerelease` か、版そのもの) を選ぶと、タグを打ち、Releases、npm、VS Code Marketplace、Open VSX に出します。途中で失敗したときは、`version` に同じ版をそのまま書いて (`0.2.0`) もう一度動かすと、出せていないところだけ出し直せます。版の決まりとやり直しの細かい決まりは [spec/distribution.md](spec/distribution.md) の「版とリリース」にあります。
- 使うトークンは [.env.example](.env.example) に並べてあります。同じ名前で GitHub の Secrets にも入れます。

## ブランチとプルリクエスト

`main` はルールセットで守っています。変更はブランチを切ってプルリクエストにし、CI (`test`・`wasm`・`installer`・`vscode`) が通ってから squash でマージします。squash のコミットの題はプルリクエストの題になるので、題はコミットメッセージと同じ決まりで書きます。マージしたブランチは自動で消えます。

`v*` のタグは、作るのも動かすのも消すのもルールセットで止めています。作れるのは release の workflow だけで、otoneko-release App のトークンを使います。そのため GitHub の Variables に `RELEASE_APP_CLIENT_ID`、Secrets に `RELEASE_APP_PRIVATE_KEY` (App の秘密鍵) を入れ、App をこのリポジトリに入れておきます。
