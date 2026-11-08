package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// hak は対話モード。1行ずつ読んで実行する。ブロックの途中などで入力が終わっていたら、続きの行を待つ。
// 続きを待っているときに空行を打つと、そこまでで実行してエラーを出す (Python と同じ)。
func (s streams) hak(context.Context, *cli.Command) error {
	// プログラムの zanfession と同じ読み手を使い、先読みした行を取り合わないようにする
	stdin := bufio.NewReader(s.stdin)
	opts := s.options(nil)
	opts.Stdin = stdin
	h := runner.NewHak(opts)
	var entry strings.Builder
	for {
		prompt := "> "
		if entry.Len() > 0 {
			prompt = "... "
		}
		fmt.Fprint(s.stdout, prompt)
		line, readErr := stdin.ReadString('\n')
		if line == "" && readErr != nil {
			fmt.Fprintln(s.stdout)
			return nil
		}
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}
		blank := strings.TrimSpace(line) == ""
		if blank && entry.Len() == 0 {
			continue
		}
		entry.WriteString(line)
		warnings, err := h.Feed(entry.String())
		if runner.IsUnfinished(err) {
			if !blank && readErr == nil {
				continue
			}
			h.Drop(entry.String())
		}
		entry.Reset()
		s.printWarnings(warnings)
		if exit, ok := errors.AsType[*interp.Exit](err); ok {
			return exitCode(exit.Code)
		}
		if err != nil {
			fmt.Fprintln(s.stderr, err)
		}
	}
}
