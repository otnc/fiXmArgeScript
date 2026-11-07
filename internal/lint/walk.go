package lint

import "github.com/otnc/fiXmArgeScript/internal/ast"

// inspect は構文木を深さ優先でたどり、各部品で visit を呼ぶ。visit が false を返すと、その部品の中には入らない。
func inspect(n ast.Node, visit func(ast.Node) bool) {
	if n == nil || !visit(n) {
		return
	}
	each := func(nodes ...ast.Node) {
		for _, c := range nodes {
			if c != nil {
				inspect(c, visit)
			}
		}
	}
	switch n := n.(type) {
	case *ast.StringLit:
		for _, p := range n.Parts {
			if p.Expr != nil {
				inspect(p.Expr, visit)
			}
		}
	case *ast.ArrayLit:
		eachExpr(n.Elems, visit)
	case *ast.DictLit:
		eachExpr(n.Keys, visit)
		eachExpr(n.Values, visit)
	case *ast.Paren:
		each(n.X)
	case *ast.Unary:
		each(n.X)
	case *ast.Binary:
		each(n.X, n.Y)
	case *ast.Logical:
		each(n.X, n.Y)
	case *ast.Compare:
		eachExpr(n.Operands, visit)
	case *ast.Exists:
		each(n.X)
	case *ast.Index:
		each(n.X, n.Key)
	case *ast.Slice:
		each(n.X)
		eachOptional(visit, n.Lo, n.Hi)
	case *ast.Range:
		each(n.Lo, n.Hi)
		eachOptional(visit, n.Step)
	case *ast.Call:
		each(n.Fn)
		eachExpr(n.Args, visit)
	case *ast.FuncLit:
		eachStmt(n.Body, visit)
	case *ast.Choice:
		eachExpr(n.Options, visit)
	case *ast.Await:
		each(n.X)
	case *ast.Input:
		eachOptional(visit, n.Prompt)

	case *ast.ExprStmt:
		each(n.X)
	case *ast.Assign:
		each(n.Target, n.Value)
	case *ast.Print:
		each(n.X)
	case *ast.Return:
		eachOptional(visit, n.X)
	case *ast.Throw:
		each(n.X)
	case *ast.Assert:
		each(n.X)
	case *ast.If:
		eachExpr(n.Conds, visit)
		for _, b := range n.Blocks {
			eachStmt(b, visit)
		}
		eachStmt(n.Else, visit)
	case *ast.While:
		each(n.Cond)
		eachStmt(n.Body, visit)
	case *ast.For:
		each(n.Iter)
		eachStmt(n.Body, visit)
	case *ast.FuncDecl:
		each(n.Fn)
	case *ast.Try:
		eachStmt(n.Body, visit)
		eachStmt(n.Catch, visit)
		eachStmt(n.Finally, visit)
	case *ast.Test:
		eachStmt(n.Body, visit)
	}
}

func eachExpr(xs []ast.Expr, visit func(ast.Node) bool) {
	for _, x := range xs {
		inspect(x, visit)
	}
}

func eachStmt(stmts []ast.Stmt, visit func(ast.Node) bool) {
	for _, s := range stmts {
		inspect(s, visit)
	}
}

// eachOptional は省ける部品 (nil の式) を飛ばしてたどる。
func eachOptional(visit func(ast.Node) bool, xs ...ast.Expr) {
	for _, x := range xs {
		if x != nil {
			inspect(x, visit)
		}
	}
}

// inspectFile はファイルの文を全部たどる。
func inspectFile(stmts []ast.Stmt, visit func(ast.Node) bool) {
	eachStmt(stmts, visit)
}

// blocks はファイルの中のすべての文の並び (ファイルの一番外側、ブロック、関数の本体) に f を呼ぶ。
func blocks(stmts []ast.Stmt, f func([]ast.Stmt)) {
	f(stmts)
	inspectFile(stmts, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.If:
			for _, b := range n.Blocks {
				f(b)
			}
			if n.Else != nil {
				f(n.Else)
			}
		case *ast.While:
			f(n.Body)
		case *ast.For:
			f(n.Body)
		case *ast.FuncLit:
			f(n.Body)
		case *ast.Try:
			f(n.Body)
			f(n.Catch)
			f(n.Finally)
		case *ast.Test:
			f(n.Body)
		}
		return true
	})
}
