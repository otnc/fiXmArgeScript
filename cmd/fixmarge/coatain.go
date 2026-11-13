package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// coatain は _coatain.fxma の sentky を実行して判決を出す。有罪が1つでもあれば終了コード1。
func (s streams) coatain(_ context.Context, cmd *cli.Command) error {
	return s.runTestFiles(cmd, interp.RunTests)
}

// persemon は _coatain.fxma の persemon を実行して速さを出す。本体の shoubek が落ちたものがあれば終了コード1。
func (s streams) persemon(_ context.Context, cmd *cli.Command) error {
	return s.runTestFiles(cmd, interp.RunBenches)
}

// runTestFiles は渡されたファイルかフォルダの _coatain.fxma を run で1つずつ実行する。
func (s streams) runTestFiles(cmd *cli.Command, run func(*parser.File, interp.Options) (bool, error)) error {
	// persemon には --laikill がなく、0 (上限なし) になる
	maxSteps := cmd.Int("laikill")
	if maxSteps < 0 {
		s.usage(s.stderr)
		return exitCode(2)
	}
	targets := cmd.Args().Slice()
	if len(targets) == 0 {
		targets = []string{"."}
	}
	var files []string
	for _, t := range targets {
		found, err := findFiles(t, "_coatain.fxma")
		if err != nil {
			fmt.Fprintln(s.stderr, readError(err))
			return exitCode(1)
		}
		files = append(files, found...)
	}
	failed := false
	for _, path := range files {
		f, err := s.prepare(path)
		if err == nil {
			var guilty bool
			opts := s.options(path, nil)
			opts.MaxSteps = maxSteps
			guilty, err = run(f, opts)
			failed = failed || guilty
		}
		if exit, ok := errors.AsType[*interp.Exit](err); ok {
			return exitCode(exit.Code)
		}
		if err != nil {
			fmt.Fprintln(s.stderr, err)
			failed = true
		}
	}
	if failed {
		return exitCode(1)
	}
	return nil
}

// findFiles は target がファイルならそれを、フォルダなら下のフォルダまで探した、名前が suffix で終わるファイルを返す
// (vitest や deno fmt と同じ)。隠しフォルダと node_modules は飛ばす。
func findFiles(target, suffix string) ([]string, error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{filepath.ToSlash(target)}, nil
	}
	var files []string
	err = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() && path != target && (strings.HasPrefix(name, ".") || name == "node_modules") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(name, suffix) {
			files = append(files, filepath.ToSlash(path))
		}
		return nil
	})
	return files, err
}
