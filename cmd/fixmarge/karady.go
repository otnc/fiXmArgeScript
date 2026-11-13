package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/karady"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modload"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// karady は sinruits --karady。下書きと読み込み先を fixmarge の実行ファイルの後ろに付けて、単体の実行ファイルにする。
// --os で今と違う OS を選ぶと、同じ版の fixmarge をリリースから取ってきて使う。
func (s streams) karady(cmd *cli.Command) error {
	if cmd.Args().Len() != 1 || !strings.HasSuffix(cmd.Args().First(), ".fxma") {
		s.usage(s.stderr)
		return exitCode(2)
	}
	entry := cmd.Args().First()
	target := karady.Target{OS: runtime.GOOS, Arch: runtime.GOARCH}
	if t := cmd.String("os"); t != "" {
		var err error
		if target, err = karady.ParseTarget(t); err != nil {
			return s.report(err)
		}
	}
	if err := checkPackages(filepath.Dir(entry)); err != nil {
		return s.report(err)
	}
	b, warnings, err := karady.Collect(entry, os.ReadFile, bundledExternal(filepath.Dir(entry)))
	s.printWarnings(warnings)
	if err != nil {
		return s.report(err)
	}
	base, err := runtimeFor(target)
	if err != nil {
		return s.report(err)
	}
	out := strings.TrimSuffix(entry, ".fxma") + target.Exe()
	if err := writeExecutable(out, base, b); err != nil {
		fmt.Fprintln(s.stderr, diag.CaNay.Spelling(), out)
		return exitCode(1)
	}
	fmt.Fprintln(s.stdout, lexicon.Spell("finend"), out)
	return nil
}

// bundledExternal は外部パッケージを、fixmarge keeon と同じく dir から上にたどった fixmarge.mod で探す。
func bundledExternal(dir string) karady.External {
	load := sync.OnceValues(func() (*modload.Module, error) { return checkedModule(dir) })
	return func(from, path string) (string, string, error) {
		m, err := load()
		if err != nil {
			return "", "", err
		}
		file, err := m.Resolve(from, path)
		if err != nil {
			return "", "", err
		}
		name, err := m.Bundled(from, path)
		return file, name, err
	}
}

// runtimeFor は target 向けの fixmarge の実行ファイル。今の OS なら自分自身を使う。
// ほかの OS は同じ版をリリースから取ってくるので、タグから作った版でなければ Ca-nay。
func runtimeFor(target karady.Target) (string, error) {
	if target.OS == runtime.GOOS && target.Arch == runtime.GOARCH {
		return os.Executable()
	}
	version := moduleVersion()
	if !modload.Released(version) {
		return "", diag.NewUnplaced(diag.CaNay, version)
	}
	cache, err := modfetch.CacheDir()
	if err != nil {
		return "", err
	}
	r := &karady.Runtimes{Root: filepath.Join(cache, "karady")}
	return r.Get(version, target)
}

func writeExecutable(out, base string, b *karady.Bundle) error {
	src, err := os.Open(base)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return err
	}
	if err := karady.Append(dst, b); err != nil {
		dst.Close()
		return err
	}
	return dst.Close()
}

// embeddedBundle は自分の実行ファイルの後ろに付いたソース。sinruits --karady で作ったものでなければ nil。
func embeddedBundle() *karady.Bundle {
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	f, err := os.Open(exe)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil
	}
	b, err := karady.Read(f, info.Size())
	if err != nil {
		return nil
	}
	return b
}

// runBundle は実行ファイルにまとめたプログラムを実行する。引数はオプションも含めてすべてプログラムに渡す。
func (s streams) runBundle(b *karady.Bundle, args []string) int {
	opts := s.options("", args)
	readFile := func(name string) ([]byte, error) {
		src, ok := b.Files[name]
		if !ok {
			return nil, os.ErrNotExist
		}
		return []byte(src), nil
	}
	opts.Load = runner.Loader(readFile)
	opts.External = func(from, path string) (string, error) {
		name, ok := b.External[from][path]
		if !ok {
			return "", diag.NewUnplaced(diag.NiiExi, path)
		}
		return name, nil
	}
	f, err := runner.Prepare(b.Files[b.Main], b.Main, readFile)
	if err == nil {
		err = interp.Run(f, opts)
	}
	code, _ := s.finishRun(err).(exitCode)
	return int(code)
}
