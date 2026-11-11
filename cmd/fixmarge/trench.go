package main

import (
	"context"
	"io"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/lsp"
)

// trench はエディタのための言語サーバーを、標準入出力で動かす (gopls と同じ)。
func (s streams) trench(ctx context.Context, _ *cli.Command) error {
	return lsp.Serve(ctx, stdio{s.stdin, s.stdout}, moduleVersion())
}

// stdio は標準入力と標準出力を、LSP の接続 (読み書きと閉じる) としてまとめる。
type stdio struct {
	io.Reader
	io.Writer
}

// Close は何もしない。標準入出力はプロセスが終わるときに閉じる。
func (stdio) Close() error { return nil }
