package ast

// MaxDepth は構文木の深さの上限。構文木をたどる処理はどれも再帰で書いてあり、深すぎる木をたどると
// Go のスタックを使い切って落ちるので、構文解析がこれより深い木を作らない。
const MaxDepth = 10000

// Inspect は構文木を深さ優先でたどり、各部品で visit を呼ぶ。visit が false を返すと、その部品の中には入らない (go/ast の Inspect と同じ)。
func Inspect(n Node, visit func(Node) bool) {
	var walk func(Node)
	walk = func(n Node) {
		if n != nil && visit(n) {
			eachChild(n, walk)
		}
	}
	walk(n)
}

// InspectStmts は文の並びを全部たどる。
func InspectStmts(stmts []Stmt, visit func(Node) bool) {
	for _, s := range stmts {
		Inspect(s, visit)
	}
}

// Imports は stmts の werer をすべて返す。mof や関数の中に書いたものも含む (実行するときはどこに書いても読み込むため)。
func Imports(stmts []Stmt) []*Import {
	var out []*Import
	InspectStmts(stmts, func(n Node) bool {
		if imp, ok := n.(*Import); ok {
			out = append(out, imp)
		}
		return true
	})
	return out
}

// Measure は stmts の木の深さを返し、関数ごとに本体の深さを書き込む (FuncLit.Depth)。
// MaxDepth より深いところがあれば、そこの部品を tooDeep に返す。
func Measure(stmts []Stmt) (depth int, tooDeep Node) {
	m := &measurer{}
	for _, s := range stmts {
		depth = max(depth, m.depth(s, 1))
	}
	return depth, m.tooDeep
}

type measurer struct{ tooDeep Node }

// depth は n から下の木の深さを返す。level は根から n までの深さで、MaxDepth を超えたらそれより下には入らない。
func (m *measurer) depth(n Node, level int) int {
	if level > MaxDepth {
		if m.tooDeep == nil {
			m.tooDeep = n
		}
		return 0
	}
	below := 0
	eachChild(n, func(c Node) { below = max(below, m.depth(c, level+1)) })
	if f, ok := n.(*FuncLit); ok {
		f.Depth = below
	}
	return below + 1
}

// eachChild は n のすぐ下の部品を、書かれた順に visit へ渡す。省かれた部品 (nil の式) は渡さない。
func eachChild(n Node, visit func(Node)) {
	exprs := func(xs ...Expr) {
		for _, x := range xs {
			if x != nil {
				visit(x)
			}
		}
	}
	stmts := func(blocks ...[]Stmt) {
		for _, b := range blocks {
			for _, s := range b {
				visit(s)
			}
		}
	}
	switch n := n.(type) {
	case *StringLit:
		for _, p := range n.Parts {
			exprs(p.Expr)
		}
	case *ArrayLit:
		exprs(n.Elems...)
	case *DictLit:
		exprs(n.Keys...)
		exprs(n.Values...)
	case *Paren:
		exprs(n.X)
	case *Unary:
		exprs(n.X)
	case *Binary:
		exprs(n.X, n.Y)
	case *Logical:
		exprs(n.X, n.Y)
	case *Compare:
		exprs(n.Operands...)
	case *Exists:
		exprs(n.X)
	case *Index:
		exprs(n.X, n.Key)
	case *Slice:
		exprs(n.X, n.Lo, n.Hi)
	case *Range:
		exprs(n.Lo, n.Hi, n.Step)
	case *Call:
		exprs(n.Fn)
		exprs(n.Args...)
	case *FuncLit:
		stmts(n.Body)
	case *Choice:
		exprs(n.Options...)
	case *Await:
		exprs(n.X)
	case *Input:
		exprs(n.Prompt)

	case *ExprStmt:
		exprs(n.X)
	case *Assign:
		exprs(n.Target, n.Value)
	case *Print:
		exprs(n.X)
	case *Return:
		exprs(n.X)
	case *Throw:
		exprs(n.X)
	case *Assert:
		exprs(n.X)
	case *If:
		exprs(n.Conds...)
		stmts(n.Blocks...)
		stmts(n.Else)
	case *While:
		exprs(n.Cond)
		stmts(n.Body)
	case *For:
		exprs(n.Iter)
		stmts(n.Body)
	case *FuncDecl:
		visit(n.Fn)
	case *Try:
		stmts(n.Body, n.Catch, n.Finally)
	case *Test:
		stmts(n.Body)
	}
}
