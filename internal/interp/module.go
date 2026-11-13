package interp

import (
	"errors"
	"io"
	"io/fs"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Module は werer で読み込んだファイル1つ、または標準パッケージ1つ。
type Module struct {
	Path     string // werer に書いた読み込み先
	members  map[string]Value
	exported map[string]bool
}

// execImport は werer の文を実行する。読み込んだ名前は showence と同じく書き換えられない (JavaScript の import と同じ)。
func (in *Interp) execImport(s *ast.Import, env *Env) error {
	m, err := in.loadModule(s.Path, s.Pos)
	if err != nil {
		return err
	}
	if s.Alias != "" {
		return in.bindImport(env, s.Alias, m, s.Pos)
	}
	for _, n := range s.Names {
		v, err := in.memberNamed(m, n.Name, s.Pos)
		if err != nil {
			return err
		}
		name := n.Name
		if n.Alias != "" {
			name = n.Alias
		}
		if err := in.bindImport(env, name, v, s.Pos); err != nil {
			return err
		}
	}
	return nil
}

// bindImport は読み込んだ名前を定義する。同じ名前がもうあれば、aste で別名を付けるよう Su·Ready で知らせる。
func (in *Interp) bindImport(env *Env, name string, v Value, pos token.Pos) error {
	if env.declared(name) {
		return in.errorf(diag.SuReady, pos, name)
	}
	env.define(name, v, true)
	return nil
}

// member は A no B の A がモジュールのとき。右の識別子は変数ではなく、モジュールの中の名前として読む。
func (in *Interp) member(m *Module, key ast.Expr, pos token.Pos) (Value, error) {
	id, ok := key.(*ast.Ident)
	if !ok {
		return nil, in.errorf(diag.ImJune, pos, lexicon.Spell("no"))
	}
	return in.memberNamed(m, id.Name, pos)
}

func (in *Interp) memberNamed(m *Module, name string, pos token.Pos) (Value, error) {
	v, ok := m.members[name]
	switch {
	case !ok:
		return nil, in.errorf(diag.NiiExi, pos, name)
	case !m.exported[name]:
		return nil, in.errorf(diag.CaNay, pos, name)
	}
	return v, nil
}

// loadModule は読み込み先のモジュールを返す。同じ実行の中では一度だけ読む。
func (in *Interp) loadModule(path string, pos token.Pos) (*Module, error) {
	if pkg, ok := strings.CutPrefix(path, modpath.StdPrefix); ok {
		return in.loadStd(path, pkg, pos)
	}
	name, err := in.moduleFile(in.frame().file, path)
	if u, ok := errors.AsType[*diag.Unplaced](err); ok {
		return nil, in.errorf(u.Kind, pos, u.Target)
	}
	if err != nil {
		return nil, err
	}
	if m, ok := in.modules[name]; ok {
		return m, nil
	}
	if in.loading[name] {
		return nil, in.errorf(diag.SuReady, pos, path)
	}
	f, err := in.parseModule(name, path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, in.errorf(diag.NiiExi, pos, path)
	}
	if err != nil {
		return nil, err
	}
	in.loading[name] = true
	defer delete(in.loading, name)
	m, err := in.runModule(path, f)
	if err != nil {
		return nil, err
	}
	in.modules[name] = m
	return m, nil
}

// moduleFile は from に書いた werer の読み込み先を、読むファイルの名前にする。
// 見つからないときなどは *diag.Unplaced を返す。
func (in *Interp) moduleFile(from, path string) (string, error) {
	if name, ok := modpath.Local(from, path); ok {
		return name, nil
	}
	if in.opts.External == nil {
		return "", &diag.Unplaced{Kind: diag.CaNay, Target: path}
	}
	return in.opts.External(path)
}

// mayChoose は f と、f から読み込むファイルのどこかに mekamo があるかを調べる。seen は調べ終えたファイル。
// 読み込めないファイルはあるものとして扱う。エラーは実行したときに werer の場所で出す。
func (in *Interp) mayChoose(f *parser.File, seen map[string]bool) bool {
	found := false
	ast.InspectStmts(f.Stmts, func(n ast.Node) bool {
		_, isChoice := n.(*ast.Choice)
		found = found || isChoice
		return !found
	})
	if found {
		return true
	}
	for _, imp := range ast.Imports(f.Stmts) {
		path := imp.Path
		if strings.HasPrefix(path, modpath.StdPrefix) {
			continue
		}
		name, err := in.moduleFile(f.Name, path)
		if err != nil || in.opts.Load == nil {
			return true
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		m, err := in.parseModule(name, path)
		if err != nil || in.mayChoose(m, seen) {
			return true
		}
	}
	return false
}

func (in *Interp) loadStd(path, pkg string, pos token.Pos) (*Module, error) {
	if m, ok := in.modules[path]; ok {
		return m, nil
	}
	newPkg, ok := stdPackages[pkg]
	if !ok {
		return nil, in.errorf(diag.NiiExi, pos, path)
	}
	m := newPkg()
	m.Path = path
	in.modules[path] = m
	return m, nil
}

// parseModule は読み込み先 path のファイル name を読む。やり直しのたびに読み直さないよう、構文木は覚えておく。
// 方言の警告は、実行するファイルと同じく初めて読んだときに出す。外部パッケージの綴りは使う側では直せないので出さない。
func (in *Interp) parseModule(name, path string) (*parser.File, error) {
	if f, ok := in.parsed[name]; ok {
		return f, nil
	}
	if in.opts.Load == nil {
		return nil, in.errorf(diag.CaNay, token.Pos{}, name)
	}
	f, err := in.opts.Load(name)
	if err != nil {
		return nil, err
	}
	in.parsed[name] = f
	if !modpath.External(path) && in.opts.Stderr != nil {
		for _, w := range f.Warnings {
			io.WriteString(in.opts.Stderr, w.String()+"\n")
		}
	}
	return f, nil
}

// runModule は読み込んだファイルの定義を実行する。読み込まれるファイルに書けるのは定義と werer だけ。
func (in *Interp) runModule(path string, f *parser.File) (*Module, error) {
	m := &Module{Path: path, members: map[string]Value{}, exported: map[string]bool{}}
	for _, s := range f.Stmts {
		switch s := s.(type) {
		case *ast.FuncDecl:
			m.exported[s.Fn.Name] = s.Export
			continue
		case *ast.Assign:
			if id, ok := s.Target.(*ast.Ident); ok && s.Const {
				m.exported[id.Name] = s.Export
				continue
			}
		case *ast.Import:
			continue
		}
		return nil, diag.New(diag.CaNay, f.Name, s.Position(), "")
	}
	env := newEnv(nil, false)
	in.frames = append(in.frames, frame{file: f.Name, nesting: in.frame().nesting + f.Depth})
	_, _, err := in.execBlock(f.Stmts, env)
	in.popFrame()
	if err != nil {
		return nil, err
	}
	for _, v := range env.vars {
		m.members[v.name] = v.value
	}
	return m, nil
}
