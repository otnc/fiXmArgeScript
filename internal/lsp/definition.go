package lsp

import (
	"os"
	"path/filepath"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modload"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/scope"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// definition は位置の名前の定義の場所を返す。werer で読み込んだ名前は、読み込み先のファイルの定義に飛ぶ。
// 標準パッケージ (":名前") の中身は fxma に組み込まれているので、飛び先はない。
func (s *Server) definition(d *document, p protocol.Position) (protocol.Location, bool) {
	t, ok := d.tokenAt(p)
	if !ok || t.Kind != token.Ident {
		return protocol.Location{}, false
	}
	if m, ok := d.member(t.Pos); ok {
		return s.exportedLocation(d, m.Module.Path, m.Name)
	}
	b, ok := d.refs[t.Pos]
	if !ok {
		return protocol.Location{}, false
	}
	if b.Kind == scope.Import {
		return s.exportedLocation(d, b.Path, b.Remote)
	}
	return d.location(b)
}

func (d *document) location(b *scope.Binding) (protocol.Location, bool) {
	name, ok := d.nameToken(b)
	if !ok {
		return protocol.Location{}, false
	}
	return protocol.Location{URI: d.uri, Range: d.tokenRange(name)}, true
}

func (s *Server) exportedLocation(from *document, path, name string) (protocol.Location, bool) {
	for _, ex := range s.exports(from, path) {
		if ex.binding.Name == name {
			return ex.doc.location(ex.binding)
		}
	}
	return protocol.Location{}, false
}

// exported は読み込み先のファイルで oqun を付けて公開した名前1つ。
type exported struct {
	doc     *document
	binding *scope.Binding
}

// exports は werer の読み込み先 path で公開している名前を返す。開いているファイルならエディタの中身を、なければディスクを読む。
func (s *Server) exports(from *document, path string) []exported {
	name, ok := importFile(from.path, path)
	if !ok {
		return nil
	}
	target := s.open(uri.File(filepath.FromSlash(name)))
	if target == nil || target.file == nil {
		return nil
	}
	var out []exported
	for _, b := range target.scope.Bindings {
		if b.Export && topLevel(target.file.Stmts, b) {
			out = append(out, exported{target, b})
		}
	}
	return out
}

// open は開いているファイルを返す。開いていなければディスクから読んで解析する (覚えてはおかない)。
func (s *Server) open(u uri.URI) *document {
	if d := s.document(u); d != nil {
		return d
	}
	src, err := os.ReadFile(u.FsPath())
	if err != nil {
		return nil
	}
	return newDocument(u, string(src))
}

// topLevel は、b がファイルの一番外側で定義した名前かを返す。公開できるのは一番外側の名前だけ。
func topLevel(stmts []ast.Stmt, b *scope.Binding) bool {
	for _, s := range stmts {
		switch s := s.(type) {
		case *ast.FuncDecl:
			if s.Fn.Name == b.Name && s.Pos == b.Pos {
				return true
			}
		case *ast.Assign:
			if id, ok := s.Target.(*ast.Ident); ok && id.Name == b.Name && id.Pos == b.Pos {
				return true
			}
		}
	}
	return false
}

// importFile は werer の読み込み先を、読むファイルの名前にする。外部パッケージは from から上にたどった fixmarge.mod で探す。
// 中身が fixmarge.sum と合うかは確かめない (それは実行するときに見る)。
func importFile(from, path string) (string, bool) {
	if name, ok := modpath.Local(from, path); ok {
		return name, true
	}
	if !modpath.External(path) {
		return "", false
	}
	root, err := modfetch.DefaultRoot()
	if err != nil {
		return "", false
	}
	m, err := modload.Load(filepath.Dir(filepath.FromSlash(from)), &modfetch.Fetcher{Root: root})
	if err != nil {
		return "", false
	}
	name, _, err := m.Locate(path)
	return name, err == nil
}
