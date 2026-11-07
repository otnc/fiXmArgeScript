package build

import (
	"errors"
	"io/fs"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// Output は清書したファイル1つ。Name は下書きの名前に c を付けたもの。
type Output struct {
	Name   string
	Source string
}

// Project は下書き entry と、そこから werer でたどれる下書きをすべて清書する (tsc がエントリから import 先をたどるのと同じ)。
// 清書を実行すると読み込み先も清書を読むので、まとめて書き出さないと動かない。
// どれか1つでも失敗したら、何も返さない。
func Project(entry string, readFile func(name string) ([]byte, error)) ([]Output, []diag.Warning, error) {
	var outputs []Output
	var warnings []diag.Warning
	seen := map[string]bool{entry: true}
	queue := []string{entry}
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		src, err := readFile(name)
		if err != nil {
			return nil, nil, err
		}
		fair, w, err := Build(string(src), name)
		warnings = append(warnings, w...)
		if err != nil {
			return nil, warnings, err
		}
		outputs = append(outputs, Output{Name: name + "c", Source: fair})
		imports, err := localImports(string(src), name, readFile)
		if err != nil {
			return nil, warnings, err
		}
		for _, next := range imports {
			if !seen[next] {
				seen[next] = true
				queue = append(queue, next)
			}
		}
	}
	return outputs, warnings, nil
}

// localImports はファイルの一番外側の werer のうち、./ か ../ で始まる読み込み先を返す。
// 読み込み先がなければ、実行するときと同じく Nii eXi'。
func localImports(src, name string, readFile func(string) ([]byte, error)) ([]string, error) {
	f, err := parser.Parse(src, name)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, s := range f.Stmts {
		imp, ok := s.(*ast.Import)
		if !ok {
			continue
		}
		next, ok := modpath.Local(name, imp.Path)
		if !ok {
			continue
		}
		if _, err := readFile(next); errors.Is(err, fs.ErrNotExist) {
			return nil, diag.New(diag.NiiExi, name, imp.Pos, imp.Path)
		}
		out = append(out, next)
	}
	return out, nil
}
