<h1 align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://fixmargescript.oto.im/logo-dark.png" />
    <img src="https://fixmargescript.oto.im/logo-light.png" alt="fiXmArge Script" width="560" />
  </picture>
</h1>

<p align="center">
  <a href="https://fixmargescript.oto.im/">📕 ホームページ</a> ·
  <a href="https://fixmargescript.oto.im/docs/">📖 ドキュメント</a> ·
  <a href="https://fixmargescript.oto.im/playground/">🎮 プレイグラウンド</a>
</p>

[fiXmArge Script](https://github.com/otnc/fiXmArgeScript) (fxma) の VS Code 拡張です。fiXmArge Script は、『魔法少女ノ魔女裁判』の架空言語「フィクスマージ語」の単語をキーワードにした、**非公式のファンメイド**のプログラミング言語です。

> [!IMPORTANT]
> この作品は「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。作品の権利は権利者にあります。
>   
> 完全無料で公開しており、収益は得ていません。

- **AI利用**: この言語の実装は、生成AIを使って制作しています。
- 準拠しているガイドライン: [「魔法少女ノ魔女裁判」配信・二次創作ガイドライン](https://manosaba.com/guidelines) (2026年10月21日更新版)

## できること

下書き (`.fxma`) と清書 (`.fxmac`) の両方で使えます。

- 色分け
- エラーと法の悪魔 (`fixmarge lavi`) の指摘を波線で出す
- キーワード、組み込み、標準パッケージの関数、ファイルの中の名前の補完
- ホバーで語の意味 (辞書の訳) とエラーの説明を出す
- 定義へ移動 (`werer` で読み込んだ名前は、読み込み先のファイルへ)
- 整形 (`fixmarge beaukshi`)

色分けのほかは、言語サーバー `fixmarge trench` が行います。先に `fixmarge` を入れてください (入れ方はリポジトリの README にあります)。

パッケージのファイル `fixmarge.mod` は、項目名と文字列に色を付けます。

## ファイルのアイコン

`.fxma`・`.fxmac`・`fixmarge.mod` (パッケージのファイル)・`fixmarge.sum` (ロックファイル) にアイコンが付きます。vscode-icons などのアイコンテーマを使っていても、テーマがアイコンを持っていないファイルにはこの拡張のアイコンが出ます。

アイコンテーマを使っていないなら、「ファイル アイコンのテーマ」で「fiXmArge Script のファイルアイコン (補助)」を選ぶと、fiXmArge Script のファイルにだけアイコンが付きます。

## 設定

| 設定 | 既定 | 説明 |
| --- | --- | --- |
| `fixmarge.path` | `fixmarge` | 言語サーバーとして動かす `fixmarge` の場所。名前だけなら、`PATH` とインストール先から探します |

名前だけのときは、`PATH` で見つからなければ、次のインストール先も探します。インストールする前から開いていた VS Code でも、再起動せずに見つけられます。

- Windows のセットアップ: `%LOCALAPPDATA%\Programs\fixmarge` (全員向けなら `%ProgramFiles%\fixmarge`)
- install.sh: `~/.fixmarge/bin` (`FIXMARGE_INSTALL` があればその下の `bin`)
- `go install`: `GOBIN` か `~/go/bin` (`GOPATH` があればその下の `bin`)

`Nii eXi' fixmarge` と出たら、`fixmarge` が見つかっていません。[インストール](https://fixmargescript.oto.im/docs/start/install/)してから VS Code の画面を再読み込みするか、`fixmarge.path` に場所を書いてください。

## ライセンス

コードのライセンスは MIT と Apache-2.0 のデュアルライセンス (`MIT OR Apache-2.0`) で、使う人がどちらかを選べます。全文は拡張に入っている `LICENSE-MIT` と `LICENSE-APACHE` にあります。このライセンスで許諾するのはこのリポジトリのコードだけで、「魔法少女ノ魔女裁判」の権利は含みません。
