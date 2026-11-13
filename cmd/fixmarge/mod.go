package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/modload"
)

// external は werer の外部パッケージを、dir から上にたどった fixmarge.mod で解決する。fixmarge.mod は初めて使うときに一度だけ読む。
func external(dir string) func(from, path string) (string, error) {
	load := sync.OnceValues(func() (*modload.Module, error) { return checkedModule(dir) })
	return func(from, path string) (string, error) {
		m, err := load()
		if err != nil {
			return "", err
		}
		return m.Resolve(from, path)
	}
}

// checkedModule は実行に使うパッケージ。fixmarge.mod と fixmarge.sum が合っていなければエラー。
func checkedModule(dir string) (*modload.Module, error) {
	m, err := loadModule(dir)
	if err != nil {
		return nil, err
	}
	return m, m.Check()
}

// checkPackages は、dir から上に fixmarge.mod があれば、fixmarge.sum と合っているかを確かめる。実行を始める前に呼ぶ。
// 合っていないまま動かすと、外部パッケージを読み込む所まで進んでから止まることになる。
func checkPackages(dir string) error {
	if modfile.Find(dir) == "" {
		return nil
	}
	_, err := checkedModule(dir)
	return err
}

// loadModule は dir から上にたどった fixmarge.mod のパッケージ。fixmarge.mod に合う版のない読み込み先は、
// グローバルに入れたもの (fixmarge werer --wid) を補欠にする。
// fixmarge.mod がなければ、グローバルに入れたものだけを読める (使い捨てのファイルと fixmarge hak のため)。
func loadModule(dir string) (*modload.Module, error) {
	m, err := modload.Load(dir, &modfetch.Fetcher{})
	if err != nil {
		if g := globalModule(); g != nil && modfile.Find(dir) == "" {
			return g, nil
		}
		return nil, err
	}
	m.FxmaVersion = moduleVersion()
	m.Global = globalModule
	return m, nil
}

// globalDir はグローバルに入れたパッケージの置き場所。fixmarge.mod と fixmarge.sum と .fixmarge を持つ、プロジェクトと同じ形のフォルダ。
// 環境変数 FIXMARGE_WID があればそこ、なければ OS の設定の場所の fixmarge/wid。
func globalDir() (string, error) {
	if dir := os.Getenv("FIXMARGE_WID"); dir != "" {
		return dir, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "fixmarge", "wid"), nil
}

// globalModule はグローバルに入れたパッケージ。まだ何も入れていないか、読めなければ nil。
var globalModule = sync.OnceValue(func() *modload.Module {
	dir, err := globalDir()
	if err != nil {
		return nil
	}
	if _, err := os.Stat(filepath.Join(dir, modfile.Name)); err != nil {
		return nil
	}
	m, err := modload.Load(dir, &modfetch.Fetcher{})
	if err != nil {
		return nil
	}
	m.FxmaVersion = moduleVersion()
	return m
})

// newGlobalModule は fixmarge werer --wid と fixmarge unneyo --wid が書き換えるグローバルのパッケージ。なければ空の一覧を作る。
func newGlobalModule() (*modload.Module, error) {
	dir, err := globalDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	name := filepath.Join(dir, modfile.Name)
	if _, err := os.Stat(name); err != nil {
		if err := os.WriteFile(name, nil, 0o644); err != nil {
			return nil, err
		}
	}
	m, err := modload.Load(dir, &modfetch.Fetcher{})
	if err != nil {
		return nil, err
	}
	m.FxmaVersion = moduleVersion()
	return m, nil
}

// fromGlobal は、ファイル path の werer がグローバルに入れたものを補欠として読むかを調べる関数を返す (fixmarge lavi の WiD)。
// fixmarge.mod のないファイルは、グローバルから読むのがふつうなので調べない。
func fromGlobal(path string) func(importPath string) bool {
	dir := filepath.Dir(path)
	if modfile.Find(dir) == "" {
		return nil
	}
	m, err := loadModule(dir)
	if err != nil {
		return nil
	}
	return func(importPath string) bool { return m.FromGlobal(path, importPath) }
}

// werer は fixmarge werer [<パッケージ>[@<版>]]。取ってきて fixmarge.mod と fixmarge.sum に書く。
// パッケージを省くと、fixmarge.mod に書いてあるものをすべて取ってくる (npm install)。
func (s streams) werer(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() > 1 {
		s.usage(s.stderr)
		return exitCode(2)
	}
	m, err := s.workModule(cmd)
	switch {
	case err != nil:
	case cmd.Args().Len() == 0:
		err = m.Install()
	default:
		_, err = m.Get(modfile.SplitSpec(cmd.Args().First()))
	}
	return s.report(err)
}

// unneyo は fixmarge unneyo <パッケージ>[@<版>]。fixmarge.mod から外し、fixmarge.sum をそろえ直す。
func (s streams) unneyo(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		s.usage(s.stderr)
		return exitCode(2)
	}
	m, err := s.workModule(cmd)
	if err == nil {
		err = m.Remove(modfile.SplitSpec(cmd.Args().First()))
	}
	return s.report(err)
}

// ado はサブコマンドを省いたとき。
func (s streams) ado(context.Context, *cli.Command) error {
	s.usage(s.stderr)
	return exitCode(2)
}

// adoMeif は fixmarge ado meif <名前>。今いる場所に fixmarge.mod を作る。
func (s streams) adoMeif(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		s.usage(s.stderr)
		return exitCode(2)
	}
	wd, err := os.Getwd()
	if err == nil {
		err = modload.Init(wd, cmd.Args().First(), moduleVersion())
	}
	return s.report(err)
}

// adoSaymal は fixmarge ado saymal。ソースの werer に合わせて fixmarge.mod と fixmarge.sum をそろえ直す。
func (s streams) adoSaymal(_ context.Context, cmd *cli.Command) error {
	m, err := s.workModule(cmd)
	if err == nil {
		err = m.Tidy()
	}
	return s.report(err)
}

// workModule はコマンドが書き換えるパッケージ。--wid ならグローバル、なければ今いる場所から上にたどった fixmarge.mod のもの。
func (s streams) workModule(cmd *cli.Command) (*modload.Module, error) {
	if cmd.Bool("wid") {
		return newGlobalModule()
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	m, err := modload.Load(wd, &modfetch.Fetcher{})
	if err != nil {
		return nil, err
	}
	m.FxmaVersion = moduleVersion()
	return m, nil
}

// report はエラーがあれば出して終了コード1にする。
func (s streams) report(err error) error {
	if err == nil {
		return nil
	}
	fmt.Fprintln(s.stderr, readError(err))
	return exitCode(1)
}
