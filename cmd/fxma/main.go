// fxma は fiXmArge Script のコマンド。
package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/explain"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// version はリリースのときに -ldflags で書き換える。
var version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// command はサブコマンド1つ。名前は辞書の語で、英語の別名も使える (RFC の実装方針)。
type command struct {
	name    string
	aliases []string
	run     func(args []string, stdin io.Reader, stdout, stderr io.Writer) int
}

func commands() []command {
	return []command{
		{"keeon", []string{"run"}, cmdKeeon},
		{"sinruits", []string{"build"}, cmdSinruits},
		{"neyar", []string{"version", "--neyar", "--version", "-v"}, cmdNeyar},
		{"nai", []string{"help", "explain", "--nai", "--help", "-n", "-h"}, cmdNai},
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return cmdNai(nil, stdin, stdout, stderr)
	}
	// コマンドを省いてファイルだけ渡すと keeon と同じ。清書の1行目の #!/usr/bin/env fxma のため
	if isSource(args[0]) {
		return cmdKeeon(args, stdin, stdout, stderr)
	}
	for _, c := range commands() {
		if args[0] == c.name || contains(c.aliases, args[0]) {
			return c.run(args[1:], stdin, stdout, stderr)
		}
	}
	fmt.Fprintln(stderr, diag.NiiExi.Spelling(), args[0])
	return 2
}

func isSource(path string) bool {
	return strings.HasSuffix(path, ".fxma") || strings.HasSuffix(path, ".fxmac")
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func readSource(path string, stderr io.Writer) (string, bool) {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(stderr, diag.NiiExi.Spelling(), path)
		return "", false
	}
	return string(src), true
}

// printReadError は、ファイルを読めなかったエラーをコンソールに出せる形にする。OS のエラー文は英語なので出さない。
func printReadError(err error, stderr io.Writer) {
	var pe *fs.PathError
	if errors.As(err, &pe) {
		fmt.Fprintln(stderr, diag.NiiExi.Spelling(), filepath.ToSlash(pe.Path))
		return
	}
	fmt.Fprintln(stderr, err)
}

func printWarnings(warnings []diag.Warning, stderr io.Writer) {
	for _, w := range warnings {
		fmt.Fprintln(stderr, w)
	}
}

func cmdKeeon(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) != 1 || !isSource(args[0]) {
		return cmdNai(nil, stdin, stderr, stderr)
	}
	src, ok := readSource(args[0], stderr)
	if !ok {
		return 1
	}
	f, err := runner.Prepare(src, args[0])
	if f != nil {
		printWarnings(f.Warnings, stderr)
	}
	if err == nil {
		err = interp.Run(f, interp.Options{Stdin: stdin, Stdout: stdout, Stderr: stderr, Load: runner.Loader(os.ReadFile)})
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func cmdSinruits(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	var files []string
	for _, a := range args {
		switch a {
		case "--karady", "--compile", "-k", "-c":
			// 単体の実行ファイルにするのは後の版で作る
			fmt.Fprintln(stderr, diag.CaNay.Spelling(), a)
			return 1
		default:
			files = append(files, a)
		}
	}
	if len(files) != 1 || !strings.HasSuffix(files[0], ".fxma") {
		return cmdNai(nil, stdin, stderr, stderr)
	}
	outputs, warnings, err := build.Project(files[0], os.ReadFile)
	printWarnings(warnings, stderr)
	if err != nil {
		printReadError(err, stderr)
		return 1
	}
	for _, out := range outputs {
		if err := os.WriteFile(out.Name, []byte(out.Source), 0o644); err != nil {
			fmt.Fprintln(stderr, diag.CaNay.Spelling(), out.Name)
			return 1
		}
		fmt.Fprintln(stdout, lexicon.Spell("finend"), out.Name)
	}
	return 0
}

func cmdNeyar(_ []string, _ io.Reader, stdout, _ io.Writer) int {
	fmt.Fprintf(stdout, "fxma %s\n", version)
	// ガイドラインが「分かりやすい位置に明記」を求めているので、ここだけ日本語で出す
	fmt.Fprintln(stdout, "「魔法少女ノ魔女裁判」の非公式のファンメイド作品です。Acacia および合同会社Re,AER とは関係がなく、公認・監修を受けていません。")
	fmt.Fprintln(stdout, "AI利用：生成AIを使って制作しています。")
	return 0
}

func cmdNai(args []string, _ io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		spelling, text, ok := explain.Lookup(strings.Join(args, " "))
		if !ok {
			fmt.Fprintln(stderr, diag.NiiExi.Spelling(), strings.Join(args, " "))
			return 1
		}
		fmt.Fprintf(stdout, "%s\n\n%s\n", spelling, text)
		return 0
	}
	fmt.Fprint(stdout, `fiXmArge Script (fxma) の使い方

  fxma keeon <ファイル>      実行する (run)。コマンドを省いてファイルだけ渡しても同じ
  fxma sinruits <下書き>     清書 .fxmac を書き出す (build)
  fxma nai <エラー名>        エラーや警告の意味を説明する (explain)
  fxma nai                   この使い方を出す (help)
  fxma neyar                 版を出す (version)
`)
	return 0
}
