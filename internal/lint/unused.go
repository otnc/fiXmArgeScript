package lint

import (
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// binding は変数1つ。report が false のもの (引数や関数の名前、公開した showence) は使わなくても指摘しない。
type binding struct {
	pos    token.Pos
	used   bool
	report bool
}

// scope はブロック1つ分の変数。実行時と同じくブロックごとに分ける (RFC の変数のスコープ)。
type scope struct {
	vars    map[string]*binding
	funcTop bool
	// funcs は後で調べる関数。関数の中から、後で定義する外の変数を読めるので、ブロックの最後まで待つ
	funcs []*ast.FuncLit
}

type definition struct {
	name   string
	pos    token.Pos
	report bool
}

// unused は Was·Da。使っていない変数と読み込み。_ で始まる名前は、わざと使わない名前として見ない (ESLint と同じ)。
func (l *linter) unused(stmts []ast.Stmt) {
	(&resolver{l: l}).block(stmts, false, nil)
}

type resolver struct {
	l      *linter
	scopes []*scope
}

func (r *resolver) top() *scope { return r.scopes[len(r.scopes)-1] }

func (r *resolver) block(stmts []ast.Stmt, funcTop bool, defs []definition) {
	r.scopes = append(r.scopes, &scope{vars: map[string]*binding{}, funcTop: funcTop})
	for _, d := range defs {
		r.define(d)
	}
	for _, s := range stmts {
		if d, ok := s.(*ast.FuncDecl); ok {
			r.define(definition{name: d.Fn.Name, pos: d.Pos})
		}
	}
	for _, s := range stmts {
		r.stmt(s)
	}
	for sc := r.top(); len(sc.funcs) > 0; {
		fn := sc.funcs[0]
		sc.funcs = sc.funcs[1:]
		r.function(fn)
	}
	for name, b := range r.top().vars {
		if b.report && !b.used && !strings.HasPrefix(name, "_") {
			r.l.report(WasDa, b.pos, name)
		}
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *resolver) function(fn *ast.FuncLit) {
	var defs []definition
	for _, p := range fn.Params {
		defs = append(defs, definition{name: p, pos: fn.Pos})
	}
	r.block(fn.Body, true, defs)
}

func (r *resolver) define(d definition) {
	r.top().vars[d.name] = &binding{pos: d.pos, report: d.report}
}

// read は名前を読んだ印を付ける。関数の外の変数も読める。
func (r *resolver) read(name string) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if b, ok := r.scopes[i].vars[name]; ok {
			b.used = true
			return
		}
	}
}

// assignable は izu の左の名前が、今の関数の中にもうある変数かを返す。なければ新しく作る (実行時の resolveAssign と同じ)。
func (r *resolver) assignable(name string) bool {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i].vars[name]; ok {
			return true
		}
		if r.scopes[i].funcTop {
			return false
		}
	}
	return false
}

func (r *resolver) stmt(s ast.Stmt) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		r.expr(s.X)
	case *ast.Assign:
		r.expr(s.Value)
		r.assign(s)
	case *ast.Print:
		r.expr(s.X)
	case *ast.Return:
		r.expr(s.X)
	case *ast.Throw:
		r.expr(s.X)
	case *ast.Assert:
		r.expr(s.X)
	case *ast.If:
		for i, c := range s.Conds {
			r.expr(c)
			r.block(s.Blocks[i], false, nil)
		}
		r.block(s.Else, false, nil)
	case *ast.While:
		r.expr(s.Cond)
		r.block(s.Body, false, nil)
	case *ast.For:
		r.expr(s.Iter)
		r.block(s.Body, false, []definition{{name: s.Var, pos: s.Pos, report: true}})
	case *ast.FuncDecl:
		r.top().funcs = append(r.top().funcs, s.Fn)
	case *ast.Try:
		r.block(s.Body, false, nil)
		var defs []definition
		if s.CatchVar != "" {
			defs = append(defs, definition{name: s.CatchVar, pos: s.Pos, report: true})
		}
		r.block(s.Catch, false, defs)
		r.block(s.Finally, false, nil)
	case *ast.Test:
		r.block(s.Body, false, nil)
	case *ast.Import:
		if s.Alias != "" {
			r.define(definition{name: s.Alias, pos: s.Pos, report: true})
		}
		for _, n := range s.Names {
			name := n.Name
			if n.Alias != "" {
				name = n.Alias
			}
			r.define(definition{name: name, pos: s.Pos, report: true})
		}
	}
}

func (r *resolver) assign(s *ast.Assign) {
	id, ok := s.Target.(*ast.Ident)
	if !ok {
		// xs no 0 izu 1 は xs を読んでいる
		r.expr(s.Target)
		return
	}
	switch {
	case s.Op != "":
		r.read(id.Name)
	case s.Const || !r.assignable(id.Name):
		// 公開した showence は外から使うので指摘しない
		r.define(definition{name: id.Name, pos: id.Pos, report: !s.Export})
	}
}

func (r *resolver) expr(x ast.Expr) {
	if x == nil {
		return
	}
	inspect(x, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			r.read(n.Name)
		case *ast.FuncLit:
			r.top().funcs = append(r.top().funcs, n)
			return false
		case *ast.Index:
			// util no double の double は、util がモジュールなら中の名前で、変数ではない
			if id, ok := n.X.(*ast.Ident); ok {
				if _, isModule := r.l.aliases[id.Name]; isModule {
					r.read(id.Name)
					return false
				}
			}
		}
		return true
	})
}
