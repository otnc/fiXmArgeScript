// Package scope は、ファイルの中の名前がどの定義を指すかを決める。
// 決まりは実行時と同じ (RFC の変数のスコープ)。lint の未使用の指摘と、LSP の定義へ移動で共有する。
package scope

import (
	"slices"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Kind は定義の種類。
type Kind int

const (
	Var      Kind = iota // izu で作った変数
	Const                // showence
	Func                 // weho で定義した関数
	Param                // 関数の引数
	LoopVar              // aste の変数
	CatchVar             // youlugive の変数
	Import               // "./util" no double werer で読み込んだ名前
	Module               // "./util" aste util werer で読み込んだモジュール
)

// Binding は定義1つ。
type Binding struct {
	Name string
	Kind Kind
	// Pos は定義した文の位置。引数は関数の、読み込みは werer の文の位置で、名前そのものの位置とは限らない
	Pos    token.Pos
	Export bool        // oqun
	Path   string      // Import と Module の読み込み先
	Remote string      // Import の、読み込み先での名前
	Reads  []token.Pos // 値を読んだ位置
}

// Member はモジュールの中の名前 (util no double の double)。
type Member struct {
	Module *Binding
	Name   string
}

// Info は Resolve の結果。
type Info struct {
	Bindings []*Binding
	// Refs は識別子の位置から、その識別子が指す定義を引く。読む所も書く所も、定義する所も入る
	Refs    map[token.Pos]*Binding
	Members map[token.Pos]Member
}

// Resolve はファイルの文を調べ、名前と定義を結びつける。
func Resolve(stmts []ast.Stmt) *Info {
	r := &resolver{info: &Info{Refs: map[token.Pos]*Binding{}, Members: map[token.Pos]Member{}}}
	r.block(stmts, false, nil)
	return r.info
}

type frame struct {
	vars    map[string]*Binding
	funcTop bool
	// funcs は後で調べる関数。関数の中から、後で定義する外の変数を読めるので、ブロックの最後まで待つ
	funcs []*ast.FuncLit
}

type resolver struct {
	info   *Info
	frames []*frame
}

func (r *resolver) top() *frame { return r.frames[len(r.frames)-1] }

func (r *resolver) block(stmts []ast.Stmt, funcTop bool, defs []*Binding) {
	r.frames = append(r.frames, &frame{vars: map[string]*Binding{}, funcTop: funcTop})
	for _, b := range defs {
		r.define(b, nil)
	}
	for _, s := range stmts {
		if d, ok := s.(*ast.FuncDecl); ok {
			r.define(&Binding{Name: d.Fn.Name, Kind: Func, Pos: d.Pos, Export: d.Export}, nil)
		}
	}
	for _, s := range stmts {
		r.stmt(s)
	}
	for f := r.top(); len(f.funcs) > 0; {
		fn := f.funcs[0]
		f.funcs = f.funcs[1:]
		r.function(fn)
	}
	r.frames = r.frames[:len(r.frames)-1]
}

func (r *resolver) function(fn *ast.FuncLit) {
	var defs []*Binding
	for _, p := range fn.Params {
		defs = append(defs, &Binding{Name: p, Kind: Param, Pos: fn.Pos})
	}
	r.block(fn.Body, true, defs)
}

// define は定義を今のブロックに置く。at は名前そのものの識別子 (あれば)。
func (r *resolver) define(b *Binding, at *ast.Ident) {
	r.top().vars[b.Name] = b
	r.info.Bindings = append(r.info.Bindings, b)
	if at != nil {
		r.info.Refs[at.Pos] = b
	}
}

func (r *resolver) lookup(name string) *Binding {
	for _, f := range slices.Backward(r.frames) {
		if b, ok := f.vars[name]; ok {
			return b
		}
	}
	return nil
}

func (r *resolver) read(id *ast.Ident) {
	if b := r.lookup(id.Name); b != nil {
		b.Reads = append(b.Reads, id.Pos)
		r.info.Refs[id.Pos] = b
	}
}

// existing は izu の左の名前が、今の関数の中にもうある変数ならそれを返す。なければ新しく作る (実行時の resolveAssign と同じ)。
func (r *resolver) existing(name string) *Binding {
	for _, f := range slices.Backward(r.frames) {
		if b, ok := f.vars[name]; ok {
			return b
		}
		if f.funcTop {
			return nil
		}
	}
	return nil
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
		r.block(s.Body, false, []*Binding{{Name: s.Var, Kind: LoopVar, Pos: s.Pos}})
	case *ast.FuncDecl:
		r.top().funcs = append(r.top().funcs, s.Fn)
	case *ast.Try:
		r.block(s.Body, false, nil)
		var defs []*Binding
		if s.CatchVar != "" {
			defs = append(defs, &Binding{Name: s.CatchVar, Kind: CatchVar, Pos: s.Pos})
		}
		r.block(s.Catch, false, defs)
		r.block(s.Finally, false, nil)
	case *ast.Test:
		r.block(s.Body, false, nil)
	case *ast.Import:
		if s.Alias != "" {
			r.define(&Binding{Name: s.Alias, Kind: Module, Pos: s.Pos, Path: s.Path}, nil)
		}
		for _, n := range s.Names {
			local := n.Name
			if n.Alias != "" {
				local = n.Alias
			}
			r.define(&Binding{Name: local, Kind: Import, Pos: s.Pos, Path: s.Path, Remote: n.Name}, nil)
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
	switch b := r.existing(id.Name); {
	case s.Op != "":
		r.read(id)
	case s.Const || b == nil:
		kind := Var
		if s.Const {
			kind = Const
		}
		r.define(&Binding{Name: id.Name, Kind: kind, Pos: id.Pos, Export: s.Export}, id)
	default:
		r.info.Refs[id.Pos] = b
	}
}

func (r *resolver) expr(x ast.Expr) {
	if x == nil {
		return
	}
	ast.Inspect(x, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Ident:
			r.read(n)
		case *ast.FuncLit:
			r.top().funcs = append(r.top().funcs, n)
			return false
		case *ast.Index:
			// util no double の double は、util がモジュールなら中の名前で、変数ではない
			mod, ok := n.X.(*ast.Ident)
			if !ok {
				return true
			}
			b := r.lookup(mod.Name)
			if b == nil || b.Kind != Module {
				return true
			}
			r.read(mod)
			if key, ok := n.Key.(*ast.Ident); ok {
				r.info.Members[key.Pos] = Member{Module: b, Name: key.Name}
			}
			return false
		}
		return true
	})
}
