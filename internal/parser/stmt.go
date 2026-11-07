package parser

import (
	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// blockEnds はブロックを閉じる語。ブロックの中の文はここで止まる。
var blockEnds = map[string]bool{"finend": true, "siva": true, "youlugive": true, "ennote": true}

func atBlockEnd(t token.Token) bool {
	return t.Kind == token.EOF || (t.Kind == token.Keyword && blockEnds[t.Name])
}

// parseStmts はブロックの終わりかファイルの終わりまで文を読む。
func (p *parser) parseStmts() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	for {
		p.skipNewlines()
		if atBlockEnd(p.peek()) {
			return stmts, nil
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, s)
		// 1行に1文。ただし weho (n) n * 2 itswa finend のように、ブロックの終わりの語は同じ行に続けられる
		if t := p.peek(); t.Kind != token.Newline && !atBlockEnd(t) {
			return nil, p.unexpected(t)
		}
	}
}

// parseBlock は文を読んで finend で閉じる。
func (p *parser) parseBlock() ([]ast.Stmt, error) {
	body, err := p.parseStmts()
	if err != nil {
		return nil, err
	}
	return body, p.expectKeyword("finend")
}

func (p *parser) parseStmt() (ast.Stmt, error) {
	t := p.peek()
	at := ast.At{Pos: t.Pos}
	if t.Kind == token.Keyword {
		switch t.Name {
		case "mof":
			return p.parseIf()
		case "weho", "abiryoku", "itsver", "oqun", "showence":
			if p.atAnonymousFunc() {
				return p.parseSimpleStmt()
			}
			return p.parseDecl()
		case "teitance":
			return p.parseTry()
		case "sentky":
			p.next()
			name := p.peek()
			if name.Kind != token.String || len(name.Parts) > 1 || (len(name.Parts) == 1 && name.Parts[0].IsCode) {
				return nil, p.unexpected(name)
			}
			p.next()
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			return &ast.Test{At: at, Name: stringText(name), Body: body}, nil
		case "keip":
			p.next()
			return &ast.Continue{At: at}, nil
		case "ketzen":
			p.next()
			return &ast.Break{At: at}, nil
		case "itswa":
			p.next()
			return &ast.Return{At: at}, nil
		}
	}
	if t.Kind == token.String && p.lineEndsWith("werer") {
		return p.parseImport()
	}
	return p.parseSimpleStmt()
}

// parseSimpleStmt は式で始まる文を読む。文の形は、式の後ろに来る語で決まる。
func (p *parser) parseSimpleStmt() (ast.Stmt, error) {
	start := p.peek()
	at := ast.At{Pos: start.Pos}
	x, err := p.parseExprWithEnd()
	if err != nil {
		return nil, err
	}

	t := p.peek()
	if op, ok := p.compoundOp(); ok || t.IsKeyword("izu") {
		return p.parseAssign(x, op, ast.At{Pos: t.Pos})
	}
	if t.Kind != token.Keyword {
		return &ast.ExprStmt{At: at, X: x}, nil
	}
	switch t.Name {
	case "azagara":
		p.next()
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		return &ast.While{At: at, Cond: x, Body: body}, nil
	case "aste":
		p.next()
		name, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		if err := p.expectKeyword("azagara"); err != nil {
			return nil, err
		}
		body, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		return &ast.For{At: at, Iter: x, Var: name.Name, Body: body}, nil
	case "redogh", "canosiri":
		p.next()
		return &ast.Print{At: at, X: x, Stderr: t.Name == "canosiri"}, nil
	case "itswa":
		p.next()
		return &ast.Return{At: at, X: x}, nil
	case "accuhearts":
		p.next()
		return &ast.Throw{At: at, X: x}, nil
	case "shoubek":
		p.next()
		return &ast.Assert{At: ast.At{Pos: t.Pos}, X: x}, nil
	}
	return &ast.ExprStmt{At: at, X: x}, nil
}

// compoundOp は今の位置が +izu のような複合代入なら、その記号を読んで返す。
func (p *parser) compoundOp() (string, bool) {
	t := p.peek()
	if t.Kind == token.Op && token.CompoundOps[t.Name] && p.peekAt(1).IsKeyword("izu") {
		p.next()
		return t.Name, true
	}
	return "", false
}

func (p *parser) parseAssign(target ast.Expr, op string, at ast.At) (*ast.Assign, error) {
	switch target.(type) {
	case *ast.Ident, *ast.Index, *ast.Slice:
	default:
		return nil, diag.New(diag.AyAmis, p.file, target.Position(), "izu")
	}
	p.next() // izu
	value, err := p.parseExprWithEnd()
	if err != nil {
		return nil, err
	}
	return &ast.Assign{At: at, Target: target, Op: op, Value: value}, nil
}

// parseDecl は weho と showence を、前に付く abiryoku・itsver・oqun と一緒に読む。
func (p *parser) parseDecl() (ast.Stmt, error) {
	start := p.peek()
	at := ast.At{Pos: start.Pos}
	var export, abiryoku, itsver bool
	for {
		t := p.peek()
		switch {
		case t.IsKeyword("oqun") && !export:
			export = true
		case t.IsKeyword("abiryoku") && !abiryoku:
			abiryoku = true
		case t.IsKeyword("itsver") && !itsver:
			itsver = true
		case t.IsKeyword("showence"):
			if abiryoku || itsver {
				return nil, p.unexpected(t)
			}
			p.next()
			name, err := p.expectIdent()
			if err != nil {
				return nil, err
			}
			if !p.peek().IsKeyword("izu") {
				return nil, p.unexpected(p.peek())
			}
			a, err := p.parseAssign(&ast.Ident{At: ast.At{Pos: name.Pos}, Name: name.Name}, "", at)
			if err != nil {
				return nil, err
			}
			a.Const, a.Export = true, export
			return a, nil
		case t.IsKeyword("weho"):
			if p.peekAt(1).Kind != token.Ident {
				return nil, p.unexpected(p.peekAt(1))
			}
			fn, err := p.parseFunc(at, abiryoku, itsver)
			if err != nil {
				return nil, err
			}
			return &ast.FuncDecl{At: at, Fn: fn, Export: export}, nil
		default:
			return nil, p.unexpected(t)
		}
		p.next()
	}
}

// atAnonymousFunc は、今の位置が名前のない関数の式 (weho (n) ...) で始まるかを調べる。
func (p *parser) atAnonymousFunc() bool {
	i := 0
	for p.peekAt(i).IsKeyword("abiryoku") || p.peekAt(i).IsKeyword("itsver") {
		i++
	}
	return p.peekAt(i).IsKeyword("weho") && p.peekAt(i+1).IsOp("(")
}

// parseFunc は weho から finend までを読む。名前のない関数も同じ形で読む。
func (p *parser) parseFunc(at ast.At, abiryoku, itsver bool) (*ast.FuncLit, error) {
	if err := p.expectKeyword("weho"); err != nil {
		return nil, err
	}
	fn := &ast.FuncLit{At: at, Abiryoku: abiryoku, Itsver: itsver}
	if p.peek().Kind == token.Ident {
		fn.Name = p.next().Name
	}
	if err := p.expectOp("("); err != nil {
		return nil, err
	}
	for !p.acceptOp(")") {
		if len(fn.Params) > 0 {
			if err := p.expectOp(","); err != nil {
				return nil, err
			}
		}
		name, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		for _, prev := range fn.Params {
			if prev == name.Name {
				return nil, diag.New(diag.SuReady, p.file, name.Pos, name.Name)
			}
		}
		fn.Params = append(fn.Params, name.Name)
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	fn.Body = body
	return fn, nil
}

func (p *parser) parseIf() (ast.Stmt, error) {
	s := &ast.If{At: ast.At{Pos: p.next().Pos}}
	for {
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		body, err := p.parseStmts()
		if err != nil {
			return nil, err
		}
		s.Conds = append(s.Conds, cond)
		s.Blocks = append(s.Blocks, body)
		if !p.acceptKeyword("siva") {
			break
		}
		if !p.acceptKeyword("mof") {
			if s.Else, err = p.parseStmts(); err != nil {
				return nil, err
			}
			if s.Else == nil {
				s.Else = []ast.Stmt{}
			}
			break
		}
	}
	return s, p.expectKeyword("finend")
}

func (p *parser) parseTry() (ast.Stmt, error) {
	s := &ast.Try{At: ast.At{Pos: p.next().Pos}}
	var err error
	if s.Body, err = p.parseStmts(); err != nil {
		return nil, err
	}
	if p.acceptKeyword("youlugive") {
		s.HasCatch = true
		if p.peek().Kind == token.Ident {
			s.CatchVar = p.next().Name
		}
		if s.Catch, err = p.parseStmts(); err != nil {
			return nil, err
		}
	}
	if p.acceptKeyword("ennote") {
		s.HasFinally = true
		if s.Finally, err = p.parseStmts(); err != nil {
			return nil, err
		}
	}
	if !s.HasCatch && !s.HasFinally {
		return nil, p.unexpected(p.peek())
	}
	return s, p.expectKeyword("finend")
}

// lineEndsWith は、今の行が name の語で終わるかを調べる。読み込みの文と、文字列で始まる式の文を見分けるのに使う。
func (p *parser) lineEndsWith(name string) bool {
	last := p.peek()
	for i := p.i; i < len(p.tokens); i++ {
		t := p.tokens[i]
		if t.Kind == token.Newline || t.Kind == token.EOF {
			break
		}
		last = t
	}
	return last.IsKeyword(name)
}

// parseImport は "./util" aste util werer と "./util" no a, b aste c werer を読む。
func (p *parser) parseImport() (ast.Stmt, error) {
	pathTok := p.next()
	if len(pathTok.Parts) > 1 || (len(pathTok.Parts) == 1 && pathTok.Parts[0].IsCode) {
		return nil, p.unexpected(pathTok)
	}
	s := &ast.Import{At: ast.At{Pos: pathTok.Pos}, Path: stringText(pathTok)}
	switch {
	case p.acceptKeyword("aste"):
		alias, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		s.Alias = alias.Name
	case p.acceptKeyword("no"):
		for {
			name, err := p.expectIdent()
			if err != nil {
				return nil, err
			}
			in := ast.ImportName{Name: name.Name, Alias: name.Name}
			if p.acceptKeyword("aste") {
				alias, err := p.expectIdent()
				if err != nil {
					return nil, err
				}
				in.Alias = alias.Name
			}
			s.Names = append(s.Names, in)
			if !p.acceptOp(",") {
				break
			}
		}
	default:
		// "./util" werer だけの書き方はない (RFC のパッケージとモジュール)
		return nil, p.unexpected(p.peek())
	}
	return s, p.expectKeyword("werer")
}

func stringText(t token.Token) string {
	if len(t.Parts) == 0 {
		return ""
	}
	return t.Parts[0].Text
}
