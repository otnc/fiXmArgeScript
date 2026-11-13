// fixmarge は fiXmArge Script のコマンド。
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/explain"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// urfave/cli は help という名前のオプションを見ると英語の help を出すので、その仕組みを外す。
// fixmarge の --help は --nai の別名で、使い方は自分で出す。
func init() { cli.HelpFlag = nil }

func main() {
	// sinruits --karady で作った実行ファイルなら、fixmarge のコマンドではなく、まとめたプログラムとして動く
	if b := embeddedBundle(); b != nil {
		os.Exit(streams{os.Stdin, os.Stdout, os.Stderr}.runBundle(b, os.Args[1:]))
	}
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// exitCode はコマンドの終了コード。urfave/cli のエラーとして返し、run で取り出す。
type exitCode int

func (c exitCode) Error() string { return fmt.Sprint(int(c)) }

// streams はコマンドが使う入出力。urfave/cli の Command には持たせず、ここから渡す。
type streams struct {
	stdin          io.Reader
	stdout, stderr io.Writer
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	s := streams{stdin, stdout, stderr}
	err := s.app().Run(context.Background(), append([]string{"fixmarge"}, args...))
	if code, ok := errors.AsType[exitCode](err); ok {
		return int(code)
	}
	if err != nil {
		// urfave/cli の読み取りのエラーは英語なので出さず、使い方を出す
		s.usage(stderr)
		return 2
	}
	return 0
}

// firstArg は、ファイルより後ろをオプションとして読まず、そのままプログラムに渡すための指定。
var firstArg = 1

// app はコマンドの一覧。名前は辞書の語で、英語の別名も使える (RFC の実装方針)。
// 短いオプションは辞書の語と英語の両方の頭文字で、ぶつかったら辞書の語の側に付けない。
func (s streams) app() *cli.Command {
	app := &cli.Command{
		Name:            "fixmarge",
		HideHelp:        true,
		HideHelpCommand: true,
		HideVersion:     true,
		Writer:          s.stdout,
		ErrWriter:       s.stderr,
		Reader:          s.stdin,
		// 既定では終了コードのエラーで os.Exit するので止める
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "neyar", Aliases: []string{"version", "v"}},
			&cli.BoolFlag{Name: "nai", Aliases: []string{"help", "n", "h"}},
		},
		Action:       s.root,
		StopOnNthArg: &firstArg,
		Commands: []*cli.Command{
			{Name: "keeon", Aliases: []string{"run"}, Action: s.keeon, StopOnNthArg: &firstArg},
			{
				Name:    "sinruits",
				Aliases: []string{"build"},
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "karady", Aliases: []string{"compile", "k", "c"}},
					&cli.StringFlag{Name: "os", Aliases: []string{"target", "o", "t"}},
				},
				Action: s.sinruits,
			},
			{Name: "coatain", Aliases: []string{"test"}, Action: s.coatain},
			{Name: "persemon", Aliases: []string{"bench"}, Action: s.persemon},
			{Name: "lavi", Aliases: []string{"lint"}, Action: s.lavi},
			{
				Name:    "beaukshi",
				Aliases: []string{"fmt"},
				Flags:   []cli.Flag{&cli.BoolFlag{Name: "curoi", Aliases: []string{"expand", "c", "e"}}},
				Action:  s.beaukshi,
			},
			{Name: "werer", Aliases: []string{"get", "install"}, Flags: []cli.Flag{widFlag()}, Action: s.werer},
			{Name: "unneyo", Aliases: []string{"uninstall", "remove"}, Flags: []cli.Flag{widFlag()}, Action: s.unneyo},
			{
				Name:    "ado",
				Aliases: []string{"mod"},
				Action:  s.ado,
				Commands: []*cli.Command{
					{Name: "meif", Aliases: []string{"init"}, Action: s.adoMeif},
					{Name: "saymal", Aliases: []string{"tidy"}, Action: s.adoSaymal},
				},
			},
			{Name: "hak", Aliases: []string{"repl"}, Action: s.hak},
			{Name: "trench", Aliases: []string{"lsp"}, Action: s.trench},
			{Name: "neyar", Aliases: []string{"version"}, Action: s.neyar},
			{Name: "nai", Aliases: []string{"help", "explain"}, Action: s.nai},
		},
	}
	for _, c := range append([]*cli.Command{app}, app.Commands...) {
		c.OnUsageError = s.usageError
		for _, sub := range c.Commands {
			sub.OnUsageError = s.usageError
		}
	}
	return app
}

// widFlag は --wid (世界)。プロジェクトではなく、グローバルに入れたパッケージを書き換える (npm の -g)。
func widFlag() cli.Flag {
	return &cli.BoolFlag{Name: "wid", Aliases: []string{"global", "g"}}
}

// usageError は知らないオプションなどのとき。urfave/cli のエラー文は英語なので出さず、使い方を出す。
func (s streams) usageError(context.Context, *cli.Command, error, bool) error {
	s.usage(s.stderr)
	return exitCode(2)
}

// root はコマンドを省いたとき。ファイルだけ渡すと keeon と同じ (清書の1行目の #!/usr/bin/env fixmarge のため)。
func (s streams) root(ctx context.Context, cmd *cli.Command) error {
	switch {
	case cmd.Bool("neyar"):
		return s.neyar(ctx, cmd)
	case cmd.Bool("nai") || cmd.Args().Len() == 0:
		s.usage(s.stdout)
		return nil
	case isSource(cmd.Args().First()):
		return s.keeon(ctx, cmd)
	}
	fmt.Fprintln(s.stderr, diag.NiiExi.Spelling(), cmd.Args().First())
	return exitCode(2)
}

func isSource(path string) bool {
	return strings.HasSuffix(path, ".fxma") || strings.HasSuffix(path, ".fxmac")
}

// readError は、ファイルを読めなかったエラーをコンソールに出せる形にする。OS のエラー文は英語なので出さない。
func readError(err error) error {
	if pe, ok := errors.AsType[*fs.PathError](err); ok {
		return errors.New(diag.NiiExi.Spelling() + " " + filepath.ToSlash(pe.Path))
	}
	return err
}

func (s streams) printWarnings(warnings []diag.Warning) {
	for _, w := range warnings {
		fmt.Fprintln(s.stderr, w)
	}
}

// prepare はファイルを読んで検査する。方言の警告はここで出す。
func (s streams) prepare(path string) (*parser.File, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, readError(err)
	}
	if err := checkPackages(filepath.Dir(path)); err != nil {
		return nil, readError(err)
	}
	f, err := runner.Prepare(string(src), path, os.ReadFile)
	if f != nil {
		s.printWarnings(f.Warnings)
	}
	return f, err
}

// options はプログラムに渡す入出力。args は dv no "conwa" になる。
// 外部パッケージは file のあるフォルダ (空なら今いる場所) から上にたどった fixmarge.mod で探す。
func (s streams) options(file string, args []string) interp.Options {
	wd, _ := os.Getwd()
	dir := wd
	if file != "" {
		dir = filepath.Dir(file)
	}
	return interp.Options{
		Stdin: s.stdin, Stdout: s.stdout, Stderr: s.stderr,
		Load: runner.Loader(os.ReadFile), External: external(dir), WorkDir: wd, FS: afero.NewOsFs(),
		Process: osProcess(args), HTTPClient: http.DefaultClient, Serve: serve, Watch: watch,
	}
}

func (s streams) keeon(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 || !isSource(cmd.Args().First()) {
		s.usage(s.stderr)
		return exitCode(2)
	}
	f, err := s.prepare(cmd.Args().First())
	if err == nil {
		err = interp.Run(f, s.options(cmd.Args().First(), cmd.Args().Tail()))
	}
	return s.finishRun(err)
}

// finishRun は実行の結果を終了コードにする。捕まらなかったエラーは出す。
func (s streams) finishRun(err error) error {
	if exit, ok := errors.AsType[*interp.Exit](err); ok {
		return exitCode(exit.Code)
	}
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return exitCode(1)
	}
	return nil
}

func (s streams) sinruits(_ context.Context, cmd *cli.Command) error {
	if cmd.Bool("karady") {
		return s.karady(cmd)
	}
	if cmd.Args().Len() != 1 || !strings.HasSuffix(cmd.Args().First(), ".fxma") {
		s.usage(s.stderr)
		return exitCode(2)
	}
	outputs, warnings, err := build.Project(cmd.Args().First(), os.ReadFile)
	s.printWarnings(warnings)
	if err != nil {
		fmt.Fprintln(s.stderr, readError(err))
		return exitCode(1)
	}
	for _, out := range outputs {
		if err := os.WriteFile(out.Name, []byte(out.Source), 0o644); err != nil {
			fmt.Fprintln(s.stderr, diag.CaNay.Spelling(), out.Name)
			return exitCode(1)
		}
		fmt.Fprintln(s.stdout, lexicon.Spell("finend"), out.Name)
	}
	return nil
}

func (s streams) neyar(context.Context, *cli.Command) error {
	fmt.Fprintln(s.stdout, "fixmarge", versionString())
	// ガイドラインが「分かりやすい位置に明記」を求めているので、ここだけ日本語で出す
	fmt.Fprintln(s.stdout, "「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。")
	fmt.Fprintln(s.stdout, "AI利用：生成AIを使って制作しています。")
	return nil
}

func (s streams) nai(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		s.usage(s.stdout)
		return nil
	}
	name := strings.Join(cmd.Args().Slice(), " ")
	spelling, text, ok := explain.Lookup(name)
	if !ok {
		fmt.Fprintln(s.stderr, diag.NiiExi.Spelling(), name)
		return exitCode(1)
	}
	fmt.Fprintf(s.stdout, "%s\n\n%s\n", spelling, text)
	return nil
}

// usage は使い方。urfave/cli の help は英語なので使わず、fixmarge nai の説明として日本語で書く。
func (s streams) usage(w io.Writer) {
	fmt.Fprint(w, `fiXmArge Script (fxma) の使い方

  fixmarge keeon <ファイル> [引数...]  実行する (run)。コマンドを省いてファイルだけ渡しても同じ
  fixmarge sinruits <下書き>     清書 .fxmac を書き出す (build)。読み込み先の下書きもまとめて清書する
  fixmarge sinruits --karady [--os <OS-CPU>] <下書き>  単体の実行ファイルを作る (--compile、--target)。
                                 OS-CPU は linux-amd64 など。省くと今の OS 向け
  fixmarge coatain [ファイル...]  _coatain.fxma の sentky を実行して判決を出す (test)。省くと下のフォルダまで全部
  fixmarge persemon [ファイル...]  _coatain.fxma の persemon で速さを測る (bench)。省くと下のフォルダまで全部
  fixmarge beaukshi [ファイル...]  下書きを整えて書き換える (fmt)。省くと下のフォルダまで全部
  fixmarge beaukshi --curoi <下書き>  呪いを展開した結果を出す (--expand)
  fixmarge lavi [ファイル...]   法の悪魔が規則違反を指摘する (lint)。省くと下のフォルダまで全部
  fixmarge ado meif <名前>       fixmarge.mod を作る (mod init)。名前はリポジトリの URL (github.com/alice/strs)
  fixmarge werer <パッケージ>[@<版>]  外部パッケージを取ってきて fixmarge.mod に足す (install)。版は 1、1.2、1.2.0、1.2.0-rc.0 の形で、
                                 合ういちばん新しい版を足す。省くといちばん新しい版
  fixmarge werer                 fixmarge.mod に書いてあるパッケージをすべて取ってくる
                                 取ってきたものは fixmarge.mod と同じ場所の .fixmarge に置く
  fixmarge werer --wid <パッケージ>  どこからでも読めるように入れる (--global、-g)。fixmarge.mod のないファイルと、
                                 fixmarge.mod に合う版がないときの補欠に使う。unneyo にも付けられる
  fixmarge unneyo <パッケージ>[@<版>]  入っている版のうち、合ういちばん新しい版を fixmarge.mod から外す (uninstall)。@* なら全部の版
  fixmarge ado saymal            werer から使うパッケージをそろえ直す (mod tidy)。足りない版を足し、使っていない版を外す
  fixmarge hak                   対話モード。1行ずつ打って実行する (repl)
  fixmarge trench                エディタのための言語サーバーを動かす (lsp)
  fixmarge nai <エラー名>        エラーや警告の意味を説明する (explain)
  fixmarge nai                   この使い方を出す (help)
  fixmarge neyar                 版を出す (version)
`)
}
