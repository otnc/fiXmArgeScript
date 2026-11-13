package main

import (
	"context"
	"fmt"
	"os"
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

func loadModule(dir string) (*modload.Module, error) {
	root, err := modfetch.DefaultRoot()
	if err != nil {
		return nil, err
	}
	m, err := modload.Load(dir, &modfetch.Fetcher{Root: root})
	if err != nil {
		return nil, err
	}
	m.FxmaVersion = moduleVersion()
	return m, nil
}

// werer は fixmarge werer [<パッケージ>[@<版>]]。取ってきて fixmarge.mod と fixmarge.sum に書く。
// パッケージを省くと、fixmarge.mod に書いてあるものをすべて取ってくる (npm install)。
func (s streams) werer(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() > 1 {
		s.usage(s.stderr)
		return exitCode(2)
	}
	m, err := s.workModule()
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
	m, err := s.workModule()
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
func (s streams) adoSaymal(context.Context, *cli.Command) error {
	m, err := s.workModule()
	if err == nil {
		err = m.Tidy()
	}
	return s.report(err)
}

// workModule は今いる場所から上にたどった fixmarge.mod のパッケージ。
func (s streams) workModule() (*modload.Module, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return loadModule(wd)
}

// report はエラーがあれば出して終了コード1にする。
func (s streams) report(err error) error {
	if err == nil {
		return nil
	}
	fmt.Fprintln(s.stderr, readError(err))
	return exitCode(1)
}
