// Package runner は、ソースを読んで検査し、実行するまでの流れをまとめる。CLI と WebAssembly 版で共有する。
package runner

import (
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/check"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// ReadFile は werer で読み込むファイルを読む。なければ fs.ErrNotExist を包んだエラーを返す。
// nil なら、読み込めるのは標準パッケージだけ。
type ReadFile = func(name string) ([]byte, error)

// Prepare はソースを読んで実行前の検査をする。
// 下書きは清書と同じ字句の並びになるので、清書を書き出さずにそのまま読む。
func Prepare(src, file string, readFile ReadFile) (*parser.File, error) {
	f, err := parser.Parse(src, file)
	if err != nil {
		return nil, err
	}
	return f, check.File(f, Modules(file, readFile))
}

// Modules は file に書いた werer の読み込み先について、abiryoku 付きの関数を調べる (check.Modules)。
// 分かるのは標準パッケージと、./ か ../ で始まるファイル。外部パッケージの付け忘れは実行したときに見つかる。
func Modules(file string, readFile ReadFile) check.Modules {
	return func(path string) map[string]bool {
		if pkg, ok := strings.CutPrefix(path, modpath.StdPrefix); ok {
			return stdAbiryoku(pkg)
		}
		name, ok := modpath.Local(file, path)
		if !ok || readFile == nil {
			return nil
		}
		src, err := readFile(name)
		if err != nil {
			return nil
		}
		// 読めないファイルのエラーは、実行したときに werer の場所で出す
		f, err := parser.Parse(string(src), name)
		if err != nil {
			return nil
		}
		return check.AbiryokuExports(f)
	}
}

func stdAbiryoku(pkg string) map[string]bool {
	out := map[string]bool{}
	for _, name := range interp.StdMembers(pkg) {
		if info, ok := interp.StdFunc(pkg, name); ok && info.Abiryoku {
			out[name] = true
		}
	}
	return out
}

// Loader は、ファイルを読む関数を werer の読み込みに使う形 (interp.Options の Load) にする。
// 読み込まれるファイルも、実行するファイルと同じく読んで検査する。
func Loader(readFile ReadFile) func(name string) (*parser.File, error) {
	if readFile == nil {
		return nil
	}
	return func(name string) (*parser.File, error) {
		src, err := readFile(name)
		if err != nil {
			return nil, err
		}
		f, err := Prepare(string(src), name, readFile)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

// Run はソースを実行する。werer のファイルは readFile で読む。方言の警告は、エラーがあっても返す。
func Run(src, file string, readFile ReadFile, opts interp.Options) ([]diag.Warning, error) {
	f, err := Prepare(src, file, readFile)
	if f == nil {
		return nil, err
	}
	if err != nil {
		return f.Warnings, err
	}
	opts.Load = Loader(readFile)
	return f.Warnings, interp.Run(f, opts)
}

// RunTests は _coatain.fxma の sentky を実行して判決を出す。有罪が1つでもあれば failed が true。
func RunTests(src, file string, readFile ReadFile, opts interp.Options) (failed bool, err error) {
	f, err := Prepare(src, file, readFile)
	if err != nil {
		return false, err
	}
	opts.Load = Loader(readFile)
	return interp.RunTests(f, opts)
}

// RunBenches は _coatain.fxma の persemon を実行して速さを出す。本体の shoubek が落ちたものがあれば failed が true。
func RunBenches(src, file string, readFile ReadFile, opts interp.Options) (failed bool, err error) {
	f, err := Prepare(src, file, readFile)
	if err != nil {
		return false, err
	}
	opts.Load = Loader(readFile)
	return interp.RunBenches(f, opts)
}
