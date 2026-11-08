package lint

import "github.com/otnc/fiXmArgeScript/internal/ast"

// blocks はファイルの中のすべての文の並び (ファイルの一番外側、ブロック、関数の本体) に f を呼ぶ。
func blocks(stmts []ast.Stmt, f func([]ast.Stmt)) {
	f(stmts)
	ast.InspectStmts(stmts, func(n ast.Node) bool {
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
