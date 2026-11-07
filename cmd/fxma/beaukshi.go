package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/format"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// beaukshi は下書きを整えて書き換える (go fmt や deno fmt と同じ)。変わったファイルだけ書き換えて名前を出す。
// --curoi のときは、呪いを展開した結果を画面に出すだけで書き換えない。
func (s streams) beaukshi(_ context.Context, cmd *cli.Command) error {
	if cmd.Bool("curoi") {
		return s.expand(cmd)
	}
	targets := cmd.Args().Slice()
	if len(targets) == 0 {
		targets = []string{"."}
	}
	failed := false
	for _, t := range targets {
		files, err := findFiles(t, ".fxma")
		if err != nil {
			fmt.Fprintln(s.stderr, readError(err))
			failed = true
			continue
		}
		for _, path := range files {
			if err := s.formatFile(path); err != nil {
				fmt.Fprintln(s.stderr, err)
				failed = true
			}
		}
	}
	if failed {
		return exitCode(1)
	}
	return nil
}

func (s streams) formatFile(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return readError(err)
	}
	out, err := format.Format(string(src), path)
	if err != nil || out == string(src) {
		return err
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return fmt.Errorf("%s %s", diag.CaNay.Spelling(), path)
	}
	fmt.Fprintln(s.stdout, lexicon.Spell("beaukshi"), path)
	return nil
}

func (s streams) expand(cmd *cli.Command) error {
	if cmd.Args().Len() != 1 || !isSource(cmd.Args().First()) {
		s.usage(s.stderr)
		return exitCode(2)
	}
	path := cmd.Args().First()
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(s.stderr, readError(err))
		return exitCode(1)
	}
	out, err := format.Expand(string(src), path)
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return exitCode(1)
	}
	fmt.Fprint(s.stdout, out)
	return nil
}
