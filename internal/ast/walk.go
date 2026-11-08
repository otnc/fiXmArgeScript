package ast

// Inspect は構文木を深さ優先でたどり、各部品で visit を呼ぶ。visit が false を返すと、その部品の中には入らない (go/ast の Inspect と同じ)。
func Inspect(n Node, visit func(Node) bool) {
	if n == nil || !visit(n) {
		return
	}
	each := func(nodes ...Node) {
		for _, c := range nodes {
			if c != nil {
				Inspect(c, visit)
			}
		}
	}
	switch n := n.(type) {
	case *StringLit:
		for _, p := range n.Parts {
			if p.Expr != nil {
				Inspect(p.Expr, visit)
			}
		}
	case *ArrayLit:
		eachExpr(n.Elems, visit)
	case *DictLit:
		eachExpr(n.Keys, visit)
		eachExpr(n.Values, visit)
	case *Paren:
		each(n.X)
	case *Unary:
		each(n.X)
	case *Binary:
		each(n.X, n.Y)
	case *Logical:
		each(n.X, n.Y)
	case *Compare:
		eachExpr(n.Operands, visit)
	case *Exists:
		each(n.X)
	case *Index:
		each(n.X, n.Key)
	case *Slice:
		each(n.X)
		eachOptional(visit, n.Lo, n.Hi)
	case *Range:
		each(n.Lo, n.Hi)
		eachOptional(visit, n.Step)
	case *Call:
		each(n.Fn)
		eachExpr(n.Args, visit)
	case *FuncLit:
		eachStmt(n.Body, visit)
	case *Choice:
		eachExpr(n.Options, visit)
	case *Await:
		each(n.X)
	case *Input:
		eachOptional(visit, n.Prompt)

	case *ExprStmt:
		each(n.X)
	case *Assign:
		each(n.Target, n.Value)
	case *Print:
		each(n.X)
	case *Return:
		eachOptional(visit, n.X)
	case *Throw:
		each(n.X)
	case *Assert:
		each(n.X)
	case *If:
		eachExpr(n.Conds, visit)
		for _, b := range n.Blocks {
			eachStmt(b, visit)
		}
		eachStmt(n.Else, visit)
	case *While:
		each(n.Cond)
		eachStmt(n.Body, visit)
	case *For:
		each(n.Iter)
		eachStmt(n.Body, visit)
	case *FuncDecl:
		each(n.Fn)
	case *Try:
		eachStmt(n.Body, visit)
		eachStmt(n.Catch, visit)
		eachStmt(n.Finally, visit)
	case *Test:
		eachStmt(n.Body, visit)
	}
}

func eachExpr(xs []Expr, visit func(Node) bool) {
	for _, x := range xs {
		Inspect(x, visit)
	}
}

func eachStmt(stmts []Stmt, visit func(Node) bool) {
	for _, s := range stmts {
		Inspect(s, visit)
	}
}

// eachOptional は省ける部品 (nil の式) を飛ばしてたどる。
func eachOptional(visit func(Node) bool, xs ...Expr) {
	for _, x := range xs {
		if x != nil {
			Inspect(x, visit)
		}
	}
}

// InspectStmts は文の並びを全部たどる。
func InspectStmts(stmts []Stmt, visit func(Node) bool) {
	eachStmt(stmts, visit)
}
