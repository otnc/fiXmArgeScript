package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/modload"
	"github.com/otnc/fiXmArgeScript/internal/release"
)

// neyar は fixmarge 自身をリリースの版に入れ替える (upgrade)。
// 版の指定は werer の @ の後ろと同じ形 (1、1.2、1.2.0、1.2.0-rc.0) で、省くと最新の正式版。
func (s streams) neyar(_ context.Context, cmd *cli.Command) error {
	spec := cmd.Args().First()
	if cmd.Args().Len() > 1 || !modfile.ValidSpec(spec) {
		s.usage(s.stderr)
		return exitCode(2)
	}
	// go run や手元のビルドは置き場所が一時フォルダなどで、入れ替えても意味がない
	current := moduleVersion()
	if !modload.Released(current) {
		return s.report(diag.NewUnplaced(diag.CaNay, current))
	}
	versions, err := (&modfetch.Fetcher{}).Versions(release.Repository)
	if err != nil {
		return s.report(err)
	}
	next := modfile.Pick(versions, spec)
	if next == "" {
		return s.report(diag.NewUnplaced(diag.NiiExi, spec))
	}
	if next == current {
		fmt.Fprintln(s.stdout, diag.SuReady.Spelling(), next)
		return nil
	}
	exe, err := os.Executable()
	if err == nil {
		exe, err = filepath.EvalSymlinks(exe)
	}
	if err != nil {
		return s.report(err)
	}
	tmp, err := os.MkdirTemp("", "fixmarge-neyar-")
	if err != nil {
		return s.report(err)
	}
	defer os.RemoveAll(tmp)
	downloaded, err := (&release.Fetcher{Root: tmp}).Get(next, release.Current())
	if err != nil {
		return s.report(err)
	}
	if err := replaceExecutable(exe, downloaded); err != nil {
		fmt.Fprintln(s.stderr, diag.CaNay.Spelling(), exe)
		return exitCode(1)
	}
	fmt.Fprintln(s.stdout, lexicon.Spell("finend"), next)
	return nil
}

// replaceExecutable は exe を next の中身に入れ替える。書きかけの exe を残さないよう、隣に書いてから名前を変える。
// Windows は動いている exe を上書きも削除もできないが名前は変えられるので、.old に退かしてから置く。
// 退かした .old は次に入れ替えるときに消す (アンインストールでも消す。packaging/installer/fixmarge.iss)。
func replaceExecutable(exe, next string) error {
	data, err := os.ReadFile(next)
	if err != nil {
		return err
	}
	info, err := os.Stat(exe)
	if err != nil {
		return err
	}
	old, staged := exe+".old", exe+".new"
	os.Remove(old)
	if err := os.WriteFile(staged, data, info.Mode().Perm()); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		if err := os.Rename(exe, old); err != nil {
			os.Remove(staged)
			return err
		}
	}
	if err := os.Rename(staged, exe); err != nil {
		if runtime.GOOS == "windows" {
			os.Rename(old, exe)
		}
		os.Remove(staged)
		return err
	}
	return nil
}
