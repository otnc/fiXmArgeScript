package lint

import (
	"math"
	"slices"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// unreachable は shi-des。itswa・ketzen・keip・accuhearts の後ろの文には届かない。
func (l *linter) unreachable(stmts []ast.Stmt) {
	blocks(stmts, func(list []ast.Stmt) {
		for i, s := range list[:max(0, len(list)-1)] {
			if word := jumpWord(s); word != "" {
				l.report(ShiDes, list[i+1].Position(), lexicon.Spell(word))
				return
			}
		}
	})
}

func jumpWord(s ast.Stmt) string {
	switch s.(type) {
	case *ast.Return:
		return "itswa"
	case *ast.Break:
		return "ketzen"
	case *ast.Continue:
		return "keip"
	case *ast.Throw:
		return "accuhearts"
	}
	return ""
}

// emptyBlocks は sAiMoku。中身のないブロック。とくに youlugive で例外を黙って握りつぶすのを防ぐ。
// 関数の本体は、わざと何もしない関数を書くことがあるので見ない (ESLint の no-empty と同じ)。
func (l *linter) emptyBlocks(stmts []ast.Stmt) {
	empty := func(list []ast.Stmt, n ast.Node, word string) {
		if len(list) == 0 {
			l.report(SAiMoku, n.Position(), lexicon.Spell(word))
		}
	}
	inspectFile(stmts, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.If:
			for _, b := range n.Blocks {
				empty(b, n, "mof")
			}
		case *ast.While:
			empty(n.Body, n, "azagara")
		case *ast.For:
			empty(n.Body, n, "azagara")
		case *ast.Try:
			empty(n.Body, n, "teitance")
			if n.HasCatch {
				empty(n.Catch, n, "youlugive")
			}
			if n.HasFinally {
				empty(n.Finally, n, "ennote")
			}
		case *ast.Test:
			empty(n.Body, n, "sentky")
		}
		return true
	})
}

// needlessAbiryoku は emBasy。入出力をしていないのに abiryoku が付いている関数。
// 中で呼ぶ関数のうち1つでも実行前に分からないものがあれば、入出力をするかもしれないので指摘しない。
func (l *linter) needlessAbiryoku(stmts []ast.Stmt) {
	inspectFile(stmts, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncLit); ok && fn.Abiryoku && !l.mayDoIO(fn) {
			l.report(EmBasy, fn.Pos, fn.Name)
		}
		return true
	})
}

func (l *linter) mayDoIO(fn *ast.FuncLit) bool {
	io := false
	eachStmt(fn.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			// 中で定義しただけの関数は、呼ぶまで何もしない
			return false
		case *ast.Print, *ast.Input:
			io = true
		case *ast.Call:
			c := l.resolve(n.Fn)
			io = io || !c.known || c.abiryoku || l.passesIO(n.Args, fn.Params)
		}
		return !io
	})
	return io
}

// passesIO は、引数に abiryoku の関数か、関数かもしれない引数を渡しているかを返す。
func (l *linter) passesIO(args []ast.Expr, params []string) bool {
	for _, a := range args {
		switch a := a.(type) {
		case *ast.FuncLit:
			if a.Abiryoku {
				return true
			}
		case *ast.Ident:
			if f, ok := l.funcs[a.Name]; (ok && f.Abiryoku) || slices.Contains(params, a.Name) {
				return true
			}
		}
	}
	return false
}

// floatingPromises は anxAn。itsver の関数を呼んで、結果を praylee せずに捨てている文。
func (l *linter) floatingPromises(stmts []ast.Stmt) {
	inspectFile(stmts, func(n ast.Node) bool {
		if s, ok := n.(*ast.ExprStmt); ok {
			if call, ok := s.X.(*ast.Call); ok {
				if c := l.resolve(call.Fn); c.known && c.itsver {
					l.report(AnxAn, call.Pos, c.name)
				}
			}
		}
		return true
	})
}

// uncheckedEda は SuZai。eda の結果を確かめずに添字に使っている。見つからないと -1 で、最後の要素になってしまう。
// 関数ごとに、eda の結果を入れた変数のうち、比べずに添字に使ったものを指摘する。
func (l *linter) uncheckedEda(stmts []ast.Stmt) {
	l.checkEdaIn(stmts)
	inspectFile(stmts, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncLit); ok {
			l.checkEdaIn(fn.Body)
		}
		return true
	})
}

func (l *linter) checkEdaIn(body []ast.Stmt) {
	fromEda := map[string]bool{}
	compared := map[string]bool{}
	var keys []*ast.Index
	eachStmt(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			return false
		case *ast.Assign:
			if id, ok := n.Target.(*ast.Ident); ok && l.isEda(n.Value) {
				fromEda[id.Name] = true
			}
		case *ast.Compare:
			for _, x := range n.Operands {
				if id, ok := x.(*ast.Ident); ok {
					compared[id.Name] = true
				}
			}
		case *ast.Index:
			if l.isEda(n.Key) {
				l.report(SuZai, n.Pos, lexicon.Spell("eda"))
			}
			keys = append(keys, n)
		}
		return true
	})
	for _, ix := range keys {
		if id, ok := ix.Key.(*ast.Ident); ok && fromEda[id.Name] && !compared[id.Name] {
			l.report(SuZai, ix.Pos, id.Name)
		}
	}
}

// comparisons は Meh·BeRac (小数を senaji で比べている) と JokeDan (型が違うので必ず ulieso になる senaji)。
func (l *linter) comparisons(stmts []ast.Stmt) {
	inspectFile(stmts, func(n ast.Node) bool {
		c, ok := n.(*ast.Compare)
		if !ok {
			return true
		}
		for i, op := range c.Ops {
			if op != "senaji" {
				continue
			}
			a, b := c.Operands[i], c.Operands[i+1]
			ta, tb := literalType(a), literalType(b)
			switch {
			case ta != "" && tb != "" && ta != tb:
				l.report(JokeDan, c.Pos, lexicon.Spell("senaji"))
			case computesFloat(a) || computesFloat(b):
				l.report(MehBeRac, c.Pos, lexicon.Spell("senaji"))
			}
		}
		return true
	})
}

// literalType は、書いた形から型が分かる値の型の名前。数は整数と小数をまとめて1つ (1 senaji 1.0 は tadact)。
func literalType(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.IntLit, *ast.FloatLit:
		return "number"
	case *ast.StringLit:
		return "string"
	case *ast.BoolLit:
		return "bool"
	case *ast.NullLit:
		return "null"
	case *ast.ArrayLit:
		return "array"
	case *ast.DictLit:
		return "dict"
	case *ast.Paren:
		return literalType(x.X)
	}
	return ""
}

// computesFloat は、端数のある小数か、小数になる計算 (/ や小数との計算) を含む式かを返す。
// 0.1 + 0.2 は 0.30000000000000004 になるので、senaji で比べると思った結果にならない。
func computesFloat(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.FloatLit:
		return x.Value != math.Trunc(x.Value)
	case *ast.Binary:
		return x.Op == "/" || hasFloatLit(x)
	case *ast.Paren:
		return computesFloat(x.X)
	case *ast.Unary:
		return computesFloat(x.X)
	}
	return false
}

func hasFloatLit(x ast.Expr) bool {
	found := false
	inspect(x, func(n ast.Node) bool {
		if _, ok := n.(*ast.FloatLit); ok {
			found = true
		}
		_, isFunc := n.(*ast.FuncLit)
		return !found && !isFunc
	})
	return found
}

// hopelessShoubek は dEs·lee。戻れる選択肢 (mekamo) がない場所の shoubek は、偽になると必ず Co-Jundic になる。
// ファイルに mekamo が1つもなく、mekamo を持つかもしれないファイルも読み込んでいないときだけ指摘する。
// sentky の本体の shoubek はテストなので見ない。
func (l *linter) hopelessShoubek(stmts []ast.Stmt) {
	hasChoice := false
	tests := map[ast.Node]bool{}
	inspectFile(stmts, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.Choice:
			hasChoice = true
		case *ast.Test:
			eachStmt(n.Body, func(m ast.Node) bool {
				_, isFunc := m.(*ast.FuncLit)
				tests[m] = true
				return !isFunc
			})
		}
		return true
	})
	for _, path := range l.importPaths() {
		if !strings.HasPrefix(path, "fxma:") {
			return
		}
	}
	if hasChoice {
		return
	}
	inspectFile(stmts, func(n ast.Node) bool {
		if a, ok := n.(*ast.Assert); ok && !tests[a] {
			l.report(DesLee, a.Pos, lexicon.Spell("shoubek"))
		}
		return true
	})
}

func (l *linter) importPaths() []string {
	var paths []string
	for _, p := range l.aliases {
		paths = append(paths, p)
	}
	for _, in := range l.named {
		paths = append(paths, in.path)
	}
	return paths
}
