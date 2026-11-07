// fxma は fiXmArge Script のコマンド。
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
// fxma の --help は --nai の別名で、使い方は自分で出す。
func init() { cli.HelpFlag = nil }

func main() {
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
	err := s.app().Run(context.Background(), append([]string{"fxma"}, args...))
	var code exitCode
	if errors.As(err, &code) {
		return int(code)
	}
	if err != nil {
		// urfave/cli の読み取りのエラーは英語なので出さず、使い方を出す
		s.usage(stderr)
		return 2
	}
	return 0
}

// app はコマンドの一覧。名前は辞書の語で、英語の別名も使える (RFC の実装方針)。
// 短いオプションは辞書の語と英語の両方の頭文字で、ぶつかったら辞書の語の側に付けない。
// firstArg は、ファイルより後ろをオプションとして読まず、そのままプログラムに渡すための指定。
var firstArg = 1

func (s streams) app() *cli.Command {
	app := &cli.Command{
		Name:            "fxma",
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
				Flags:   []cli.Flag{&cli.BoolFlag{Name: "karady", Aliases: []string{"compile", "k", "c"}}},
				Action:  s.sinruits,
			},
			{Name: "coatain", Aliases: []string{"test"}, Action: s.coatain},
			{Name: "neyar", Aliases: []string{"version"}, Action: s.neyar},
			{Name: "nai", Aliases: []string{"help", "explain"}, Action: s.nai},
		},
	}
	for _, c := range append([]*cli.Command{app}, app.Commands...) {
		c.OnUsageError = s.usageError
	}
	return app
}

// usageError は知らないオプションなどのとき。urfave/cli のエラー文は英語なので出さず、使い方を出す。
func (s streams) usageError(context.Context, *cli.Command, error, bool) error {
	s.usage(s.stderr)
	return exitCode(2)
}

// root はコマンドを省いたとき。ファイルだけ渡すと keeon と同じ (清書の1行目の #!/usr/bin/env fxma のため)。
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
	f, err := runner.Prepare(string(src), path)
	if f != nil {
		s.printWarnings(f.Warnings)
	}
	return f, err
}

// options はプログラムに渡す入出力。args は dv no "conwa" になる。
func (s streams) options(args []string) interp.Options {
	wd, _ := os.Getwd()
	return interp.Options{
		Stdin: s.stdin, Stdout: s.stdout, Stderr: s.stderr,
		Load: runner.Loader(os.ReadFile), WorkDir: wd, FS: afero.NewOsFs(),
		Process: osProcess(args), HTTPClient: http.DefaultClient,
	}
}

func (s streams) keeon(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 || !isSource(cmd.Args().First()) {
		s.usage(s.stderr)
		return exitCode(2)
	}
	f, err := s.prepare(cmd.Args().First())
	if err == nil {
		err = interp.Run(f, s.options(cmd.Args().Tail()))
	}
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
		// 単体の実行ファイルにするのは後の版で作る
		fmt.Fprintln(s.stderr, diag.CaNay.Spelling(), "--"+lexicon.Spell("karady"))
		return exitCode(1)
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
	fmt.Fprintln(s.stdout, "fxma", versionString())
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

// usage は使い方。urfave/cli の help は英語なので使わず、fxma nai の説明として日本語で書く。
func (s streams) usage(w io.Writer) {
	fmt.Fprint(w, `fiXmArge Script (fxma) の使い方

  fxma keeon <ファイル> [引数...]  実行する (run)。コマンドを省いてファイルだけ渡しても同じ
  fxma sinruits <下書き>     清書 .fxmac を書き出す (build)。読み込み先の下書きもまとめて清書する
  fxma coatain [ファイル...]  _coatain.fxma の sentky を実行して判決を出す (test)。省くと下のフォルダまで全部
  fxma nai <エラー名>        エラーや警告の意味を説明する (explain)
  fxma nai                   この使い方を出す (help)
  fxma neyar                 版を出す (version)
`)
}
