// Package runner は、ソースを読んで検査し、実行するまでの流れをまとめる。CLI と WebAssembly 版で共有する。
package runner

import (
	"github.com/otnc/fiXmArgeScript/internal/check"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// Prepare はソースを読んで実行前の検査をする。
// 下書きは清書と同じ字句の並びになるので、清書を書き出さずにそのまま読む。
func Prepare(src, file string) (*parser.File, error) {
	f, err := parser.Parse(src, file)
	if err != nil {
		return nil, err
	}
	return f, check.File(f)
}

// Run はソースを実行する。方言の警告は、エラーがあっても返す。
func Run(src, file string, opts interp.Options) ([]diag.Warning, error) {
	f, err := Prepare(src, file)
	if f == nil {
		return nil, err
	}
	if err != nil {
		return f.Warnings, err
	}
	return f.Warnings, interp.Run(f, opts)
}
