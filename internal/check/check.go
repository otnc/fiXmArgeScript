// Package check は実行する前に見つけられるエラーを探す。
// abiryoku の付け忘れ (kiNbOO)、showence への再代入 (blasTOac)、_coatain.fxma 以外の sentky と persemon (Ca-nay)。
package check

import (
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// IsTestFile は sentky と persemon を書けるファイルかを返す。
func IsTestFile(name string) bool {
	return strings.HasSuffix(name, "_coatain.fxma") || strings.HasSuffix(name, "_coatain.fxmac")
}

// File はファイル1つを検査し、最初に見つけたエラーを返す。
func File(f *parser.File) error {
	c := &checker{file: f.Name, abiryokuFuncs: map[string]bool{}, testFile: IsTestFile(f.Name)}
	collectAbiryoku(f.Stmts, c.abiryokuFuncs)
	c.pushScope()
	c.stmts(f.Stmts, context{abiryoku: true})
	return c.err
}

// context は今いる関数の性質。ファイルの一番外側は入出力をしてよい。
type context struct {
	abiryoku bool
	params   map[string]bool // 名前で関数を探すとき、引数に隠された名前は外す
}

type checker struct {
	file          string
	testFile      bool
	abiryokuFuncs map[string]bool
	scopes        []map[string]bool // 名前 → showence か
	err           error
}

// collectAbiryoku は weho で定義した abiryoku 付きの関数の名前を集める。巻き上げがあるので先に全部集める。
func collectAbiryoku(stmts []ast.Stmt, out map[string]bool) {
	walkStmts(stmts, func(s ast.Stmt) {
		if d, ok := s.(*ast.FuncDecl); ok && d.Fn.Abiryoku {
			out[d.Fn.Name] = true
		}
	})
}

func (c *checker) fail(kind diag.Kind, n ast.Node, target string) {
	if c.err == nil {
		c.err = diag.New(kind, c.file, n.Position(), target)
	}
}

func (c *checker) pushScope() { c.scopes = append(c.scopes, map[string]bool{}) }
func (c *checker) popScope()  { c.scopes = c.scopes[:len(c.scopes)-1] }

func (c *checker) block(stmts []ast.Stmt, ctx context) {
	c.pushScope()
	c.stmts(stmts, ctx)
	c.popScope()
}

func (c *checker) stmts(stmts []ast.Stmt, ctx context) {
	for _, s := range stmts {
		c.stmt(s, ctx)
	}
}

func (c *checker) isConst(name string) bool {
	for _, v := range slices.Backward(c.scopes) {
		if isConst, ok := v[name]; ok {
			return isConst
		}
	}
	return false
}

func (c *checker) stmt(s ast.Stmt, ctx context) {
	switch s := s.(type) {
	case *ast.ExprStmt:
		c.expr(s.X, ctx)
	case *ast.Assign:
		c.expr(s.Value, ctx)
		if id, ok := s.Target.(*ast.Ident); ok {
			if !s.Const && c.isConst(id.Name) {
				c.fail(diag.BlasTOac, s, id.Name)
			}
			if _, declared := c.scopes[len(c.scopes)-1][id.Name]; !declared || s.Const {
				c.scopes[len(c.scopes)-1][id.Name] = s.Const
			}
		} else {
			c.expr(s.Target, ctx)
		}
	case *ast.Print:
		c.io(s, ctx)
		c.expr(s.X, ctx)
	case *ast.Return:
		c.expr(s.X, ctx)
	case *ast.Throw:
		c.expr(s.X, ctx)
	case *ast.Assert:
		c.expr(s.X, ctx)
	case *ast.If:
		for i, cond := range s.Conds {
			c.expr(cond, ctx)
			c.block(s.Blocks[i], ctx)
		}
		c.block(s.Else, ctx)
	case *ast.While:
		c.expr(s.Cond, ctx)
		c.block(s.Body, ctx)
	case *ast.For:
		c.expr(s.Iter, ctx)
		c.pushScope()
		c.scopes[len(c.scopes)-1][s.Var] = false
		c.stmts(s.Body, ctx)
		c.popScope()
	case *ast.FuncDecl:
		c.scopes[len(c.scopes)-1][s.Fn.Name] = false
		c.funcLit(s.Fn)
	case *ast.Try:
		c.block(s.Body, ctx)
		c.pushScope()
		if s.CatchVar != "" {
			c.scopes[len(c.scopes)-1][s.CatchVar] = false
		}
		c.stmts(s.Catch, ctx)
		c.popScope()
		c.block(s.Finally, ctx)
	case *ast.Import:
		// 読み込んだ名前は showence と同じく書き換えられない
		if s.Alias != "" {
			c.scopes[len(c.scopes)-1][s.Alias] = true
		}
		for _, n := range s.Names {
			name := n.Name
			if n.Alias != "" {
				name = n.Alias
			}
			c.scopes[len(c.scopes)-1][name] = true
		}
	case *ast.Test:
		if !c.testFile {
			c.fail(diag.CaNay, s, s.Keyword())
		}
		c.block(s.Body, ctx)
	}
}

func (c *checker) io(n ast.Node, ctx context) {
	if !ctx.abiryoku {
		c.fail(diag.KiNbOO, n, "")
	}
}

func (c *checker) funcLit(fn *ast.FuncLit) {
	ctx := context{abiryoku: fn.Abiryoku, params: map[string]bool{}}
	c.pushScope()
	for _, p := range fn.Params {
		ctx.params[p] = true
		c.scopes[len(c.scopes)-1][p] = false
	}
	c.stmts(fn.Body, ctx)
	c.popScope()
}

// isAbiryokuName は、名前が abiryoku 付きの関数を指していると実行前に分かるかを返す。
func (c *checker) isAbiryokuName(x ast.Expr, ctx context) bool {
	id, ok := x.(*ast.Ident)
	return ok && !ctx.params[id.Name] && c.abiryokuFuncs[id.Name]
}

func (c *checker) expr(x ast.Expr, ctx context) {
	if x == nil {
		return
	}
	switch x := x.(type) {
	case *ast.Call:
		if !ctx.abiryoku {
			if c.isAbiryokuName(x.Fn, ctx) {
				c.fail(diag.KiNbOO, x, x.Fn.(*ast.Ident).Name)
			}
			for _, a := range x.Args {
				if c.isAbiryokuName(a, ctx) {
					c.fail(diag.KiNbOO, a, a.(*ast.Ident).Name)
				}
			}
		}
		c.expr(x.Fn, ctx)
		c.exprs(x.Args, ctx)
	case *ast.Input:
		c.io(x, ctx)
		c.expr(x.Prompt, ctx)
	case *ast.FuncLit:
		c.funcLit(x)
	case *ast.StringLit:
		for _, p := range x.Parts {
			c.expr(p.Expr, ctx)
		}
	case *ast.ArrayLit:
		c.exprs(x.Elems, ctx)
	case *ast.DictLit:
		c.exprs(x.Keys, ctx)
		c.exprs(x.Values, ctx)
	case *ast.Paren:
		c.expr(x.X, ctx)
	case *ast.Unary:
		c.expr(x.X, ctx)
	case *ast.Binary:
		c.expr(x.X, ctx)
		c.expr(x.Y, ctx)
	case *ast.Logical:
		c.expr(x.X, ctx)
		c.expr(x.Y, ctx)
	case *ast.Compare:
		c.exprs(x.Operands, ctx)
	case *ast.Exists:
		c.expr(x.X, ctx)
	case *ast.Index:
		c.expr(x.X, ctx)
		c.expr(x.Key, ctx)
	case *ast.Slice:
		c.expr(x.X, ctx)
		c.expr(x.Lo, ctx)
		c.expr(x.Hi, ctx)
	case *ast.Range:
		c.expr(x.Lo, ctx)
		c.expr(x.Hi, ctx)
		c.expr(x.Step, ctx)
	case *ast.Choice:
		c.exprs(x.Options, ctx)
	case *ast.Await:
		c.expr(x.X, ctx)
	}
}

func (c *checker) exprs(xs []ast.Expr, ctx context) {
	for _, x := range xs {
		c.expr(x, ctx)
	}
}

// walkStmts は関数の中も含めて、すべての文に f を呼ぶ。
func walkStmts(stmts []ast.Stmt, f func(ast.Stmt)) {
	for _, s := range stmts {
		f(s)
		switch s := s.(type) {
		case *ast.If:
			for _, b := range s.Blocks {
				walkStmts(b, f)
			}
			walkStmts(s.Else, f)
		case *ast.While:
			walkStmts(s.Body, f)
		case *ast.For:
			walkStmts(s.Body, f)
		case *ast.FuncDecl:
			walkStmts(s.Fn.Body, f)
		case *ast.Try:
			walkStmts(s.Body, f)
			walkStmts(s.Catch, f)
			walkStmts(s.Finally, f)
		case *ast.Test:
			walkStmts(s.Body, f)
		}
	}
}
