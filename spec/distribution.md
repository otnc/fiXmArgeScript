# 配布

CLIはGitHub ReleasesでOSごとの圧縮ファイルを配り、WebAssembly版はnpmとブラウザのプレイグラウンドで配ります。DenoやBunなど、他の言語処理系と同じ形にそろえます。リポジトリは `github.com/otnc/fiXmArgeScript` です。

## CLIの配布物

タグを打つと、GitHub Actions上のGoReleaserが次の圧縮ファイルを作ります。

| OS | CPU | ファイル |
| --- | --- | --- |
| Linux | x86\_64 / arm64 | `fixmarge_<版>_linux_<CPU>.tar.gz` |
| macOS | x86\_64 / arm64 | `fixmarge_<版>_darwin_<CPU>.tar.gz` |
| Windows | x86\_64 / arm64 | `fixmarge_<版>_windows_<CPU>.zip` |

圧縮ファイルの中身は共通です。

- `fixmarge` (Windowsでは `fixmarge.exe`)
- `README.md`: 同梱のREADME。日本語で書き、二次創作ガイドラインで求められた表示を入れる (後の章)
- `LICENSE-MIT` と `LICENSE-APACHE`: この言語のコードのライセンス
- `THIRD_PARTY_LICENSES`: Goのランタイムなど、実行ファイルに含まれる他者のコードのライセンス

一緒に `checksums.txt` (SHA-256) も公開します。

## インストール方法

| 方法 | コマンド | 時期 |
| --- | --- | --- |
| インストーラ (macOS・Linux) | `curl -fsSL https://fixmargescript.oto.im/install.sh \| sh` | 最初から |
| インストーラ (Windows) | Releasesの `fixmarge_<版>_windows_setup.exe` (Inno Setup製) | 最初から |
| 手動 | Releasesから展開してPATHの通った場所に置く | 最初から |
| Go | `go install github.com/otnc/fiXmArgeScript/cmd/fixmarge@latest` | 最初から |
| Homebrew・Scoop | 自前のtap・bucketを用意する | 後で |

版番号はセマンティックバージョニングです。`fixmarge gun` は版番号と一緒に、非公式のファンメイド作品であることを表示します。

## 版とリリース

版の元はGitのタグ (`v0.2.0`) だけにして、版をファイルに手で書きません。

- **CLI**: Go 1.24からは、ビルドのときにGitのタグとコミットから版が実行ファイルに埋め込まれます。`fixmarge gun` はそれを読んで、版とコミットとその日時を出します。タグの上なら `v0.2.0`、タグより後なら `v0.2.1-0.<日時>-<コミット>`、手元の変更が残っていれば `+dirty` が付きます。`go install ...@v0.2.0` で入れた場合も同じく `v0.2.0` になります。
- **全部同じ版**: CLI、WebAssembly版 (npm)、VS Code拡張、LSP、Neovim・Vimのプラグインは、1つのタグで同じ版にそろえます。言語の仕様とLSPと拡張は一緒に変わるからです。
- **リリースの手順**: GitHub Actionsの `release` を `workflow_dispatch` で動かし、次の版 (`patch`・`minor`・`major`・`prerelease`、または版を直接書く) と、プレリリースの識別子 (`rc`・`beta` など、空でもよい) を選びます。最新のタグから次の版を決め、テストを通してからタグを打ち、GoReleaserでReleasesを作ります。版の計算はnpmの `semver` と同じです。
  - 識別子なしの `patch`・`minor`・`major` は正式版になる (`v0.1.0` の次の `minor` は `v0.2.0`、`v0.2.0-rc.1` の次の `minor` は `v0.2.0`)。
  - 識別子ありの `patch`・`minor`・`major` は、その版のプレリリースの最初になる (`v0.1.0` の次の `minor` と `rc` は `v0.2.0-rc.0`)。
  - `prerelease` はプレリリースの番号を1つ進める (`v0.2.0-rc.0` の次は `v0.2.0-rc.1`)。識別子を書けば識別子を変えられる (`rc` から `beta` へは戻れない。下の決まりのため)。最新が正式版のときは識別子が要る (`v0.1.0` の次の `prerelease` と `rc` は `v0.1.1-rc.0`)。
  - 版を直接書くときは、識別子を空にして版に書く (`0.2.0-rc.0`)。
  - プレリリースの番号は省ける。`1.2.0-rc` と書くと、`v1.2.0-rc.<番号>` のタグのうちいちばん新しい番号の次になる (`rc.1` まであれば `1.2.0-rc.2`)。そのタグがひとつもなければ `1.2.0-rc.0`。やり直しは番号まで書く。
  - 次の版は最新のタグより新しくなければならない。ただし、最新のタグと同じ版を直接書いたときは、途中で失敗したリリースのやり直しになる (下の「リリースのやり直し」)。
  - プレリリースは `alpha`・`beta`・`rc` のどれかの識別子と、0〜999の番号1つにする (`0.2.0-rc.1`)。識別子はnpmのdist-tagに使い、段階と番号はWindowsのセットアップのファイルバージョンに入れるため。
- **リリースのやり直し**: リリースは、タグ、Releases、npm、Windowsのセットアップ、VS Code拡張の順に出していくので、途中で失敗すると一部だけ出た状態になります。そのときは `release` を、`version` に最新の版をそのまま書いて (`0.2.0`) もう一度動かします。タグは打ち直さず、そのタグの中身でテストから作り直し、どの手順も2回目に動いても壊れないようにしてあります。Releasesのファイルは作り直したもので置き換え、npm・Marketplace・Open VSXは出し済みなら飛ばします (同じ版はもう一度出せないため)。やり直せるのは最新の版だけです。`release` の後の手順 (`npm`・`installer`・`vscode`) だけが失敗したときは、GitHub Actionsの「Re-run failed jobs」でも同じことができます。`release` の手順そのものを「Re-run」すると、`patch` などをもう一度計算して版が1つ進んでしまうので使いません。
- **Windowsのセットアップのファイルバージョン**: Explorerの「ファイル バージョン」は数字4つだけなので、4つ目でプレリリースかどうかを表します。正式版は最大の `65535`、プレリリースは段階×1000+番号です (`alpha.N` は `1000+N`、`beta.N` は `2000+N`、`rc.N` は `3000+N`)。`0.2.0-rc.1` は `0.2.0.3001`、`0.2.0` は `0.2.0.65535` になり、プレリリースが正式版より古く並びます。製品バージョンの文字列は版そのもの (`0.2.0-rc.1`) です。
- **同梱するライセンス**: `THIRD_PARTY_LICENSES` は、リリースのときにGo本体のライセンスと、`go-licenses` で集めた依存のライセンスから作ります。CLIとWebAssembly版で同じファイルを使います。
- **npmのパッケージ**: リポジトリの `package.json` の `version` は `0.0.0` のままにし、リリースのときだけ書き換えて公開します (コミットはしません)。dist-tagは、プレリリースなら識別子 (`0.2.0-rc.1` なら `rc`)、正式版なら `latest` です。
- **VS Code拡張**: Marketplaceは `0.2.0-rc.1` のようなプレリリースの版を受け付けないので、Marketplaceには正式版だけを出します。VSCodiumやCursorなどのために、同じVSIXをOpen VSXにも出します。VSIXファイル (`fixmarge_<版>_vscode.vsix`) は、どの版でもReleasesに添付します。
- **Neovim・Vimのプラグイン**: 版のファイルは持たず、リポジトリのタグをそのまま使います。lazy.nvimの `version` もNeovimの `vim.pack` も、Gitのタグで版を選びます。
- **インストーラ**: Inno Setupには、ビルドのときに版を渡します (`iscc /DAppVersion=0.2.0`)。`install.sh` はReleasesの最新版を見るので、版を持ちません。

## WebAssembly版

`GOOS=js GOARCH=wasm` で `fixmarge.wasm` を作り、Go付属の `wasm_exec.js` と組み合わせて配ります。

- **Releases**: `fixmarge_<版>_wasm.zip` に `fixmarge.wasm`、`wasm_exec.js`、`README.md`、ライセンス類を入れる。
- **npm**: パッケージ `fixmarge` として公開し、JavaScriptから `run(source, { stdin, files, args, maxSteps })` と `build(source)` を呼べるようにする。`run` は `{ stdout, stderr, exitCode, files }` を返す。`maxSteps` は `--laikill` と同じ上限。`files` は名前と中身の組で、`werer` でも `:memoid` でも読め、実行した後のファイルも返す。エディタの色分けのために、下書きに書ける名前の一覧を返す `words()` も出す。Node.js (20.16以降) とブラウザの両方で動く。中身はリポジトリの `packaging/npm/` にある。
- **プレイグラウンド**: https://fixmargescript.oto.im/playground/ に置き、ブラウザだけで書いて実行できるようにする。CodeMirror 6で作り、サイト (リポジトリの `site/`。Astro と Starlight で、紹介 `/` とドキュメント `/docs/` も同じ Cloudflare Workers に置く) の1ページにしている。ファイルを足して `werer` で読み込める。止まらないプログラムでタブが固まらないよう、文と繰り返しの数に500万の上限を付ける。ボタンや欄の名前はコマンドと同じ語 (`keeon`、`sinruits`、`zanfession`、`redogh`) にし、非公式の表記とAI利用は日本語でページの頭に出す。

ブラウザでは対話的な入力ができないため、`zanfession` で読む入力は実行前にまとめて渡します。これは巻き戻しの「入力を記録して使い回す」方式とも合います。実行ファイルが大きすぎる場合は、後でTinyGoでのビルドを検討します。

# インストーラ

インストーラ本体は既存のツールで作り、画面の文言は日本語にします。フィクスマージ語の語彙ではインストーラの文言 (エラー、確認、ディスクの容量など) を表しきれないためです。中身は `packaging/installer/` にあります。

## Windows: Inno Setup

- Inno Setupで `fixmarge_<版>_windows_setup.exe` を作り、Releasesで配る。x86\_64とarm64の実行ファイルを1つのセットアップに入れ、インストールするときにCPUに合うほうを選ぶ。
- 文言はInno Setup付属の日本語の言語ファイル (`Japanese.isl`) を使う。
- 既定では管理者権限なしで自分だけに入れる (`%LOCALAPPDATA%\Programs\fixmarge`)。起動したときに、全員向け (`Program Files`、管理者権限が要る) も選べる。入れた場所を、自分だけならユーザーの、全員向けならシステムの `PATH` に足し、アンインストールで取り除く。
- Inno Setupは非営利なら無料で使える。6.5.0から営利利用には有償ライセンスを求めているが、この作品は非営利なので対象外。
- 署名していない `.exe` は、起動時にWindowsのSmartScreenの警告が出る。最初の版では署名せず、警告の出方と実行の手順をREADMEに書く。使う人が増えたら、オープンソース向けの無償の署名サービスに申し込む。
- ライセンス同意の画面には、READMEの非公式表記を出す。ガイドラインの「分かりやすい位置に明記」を満たすためです。
- リリースでは、GoReleaserが作ったWindows向けのzipからセットアップを作り、Releasesと `checksums.txt` に足す。

## macOS・Linux: install.sh

`curl -fsSL https://fixmargescript.oto.im/install.sh | sh` で入れます。`sh -s v0.2.0` のように版を渡すと、その版を入れます。渡さなければ最新の正式版です。

- OSとCPUを調べ、圧縮ファイルと `checksums.txt` を取得してSHA-256を照合し、`~/.fixmarge/bin/fixmarge` に置く。READMEとライセンス類は `~/.fixmarge` に置く。場所は環境変数 `FIXMARGE_INSTALL` で変えられる。
- 使っているシェル (bash・zsh・fish) の設定ファイルに、`PATH` に足す行を書く (Bunと同じ)。すでに `PATH` に入っていれば書かない。macOSのbashは `.bash_profile`、Linuxのbashは `.bashrc`。ほかのシェルでは、足す行を表示するだけにする。
- `install.sh` はプレイグラウンドと同じ場所 (Cloudflare Workersの静的アセット) に置く。
- HomebrewやmacOSの `.pkg` は後で検討する。
