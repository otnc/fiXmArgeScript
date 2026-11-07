package parser

import (
	"strconv"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// parseExprWithEnd は式と、その後ろの式の終わりの語 (mekamo・praylee・zanfession) を読む。
func (p *parser) parseExprWithEnd() (ast.Expr, error) {
	var x ast.Expr
	if t := p.peek(); t.IsKeyword("zanfession") {
		// name izu zanfession のように、プロンプトなしで書ける
		p.next()
		x = &ast.Input{At: ast.At{Pos: t.Pos}}
	} else {
		var err error
		if x, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}
	for {
		t := p.peek()
		at := ast.At{Pos: t.Pos}
		switch {
		case t.IsKeyword("mekamo"):
			x = &ast.Choice{At: at, Options: candidates(x)}
		case t.IsKeyword("praylee"):
			x = &ast.Await{At: at, X: x}
		case t.IsKeyword("zanfession"):
			x = &ast.Input{At: at, Prompt: x}
		default:
			return x, nil
		}
		p.next()
	}
}

// candidates は mekamo の文の owk を候補の区切りとして読む。括弧の中の owk は「または」のまま。
func candidates(x ast.Expr) []ast.Expr {
	if l, ok := x.(*ast.Logical); ok && l.Op == "owk" {
		return append(candidates(l.X), candidates(l.Y)...)
	}
	return []ast.Expr{x}
}

func (p *parser) parseExpr() (ast.Expr, error) { return p.parseOr() }

func (p *parser) parseOr() (ast.Expr, error)  { return p.parseLogical("owk", p.parseAnd) }
func (p *parser) parseAnd() (ast.Expr, error) { return p.parseLogical("sond", p.parseNot) }

func (p *parser) parseLogical(op string, operand func() (ast.Expr, error)) (ast.Expr, error) {
	x, err := operand()
	if err != nil {
		return nil, err
	}
	for p.peek().IsKeyword(op) {
		at := ast.At{Pos: p.next().Pos}
		y, err := operand()
		if err != nil {
			return nil, err
		}
		x = &ast.Logical{At: at, Op: op, X: x, Y: y}
	}
	return x, nil
}

func (p *parser) parseNot() (ast.Expr, error) {
	if t := p.peek(); t.IsKeyword("nii") {
		p.next()
		x, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{At: ast.At{Pos: t.Pos}, Op: "nii", X: x}, nil
	}
	return p.parseComparison()
}

func (p *parser) parseComparison() (ast.Expr, error) {
	first, err := p.parseRange()
	if err != nil {
		return nil, err
	}
	cmp := &ast.Compare{At: ast.At{Pos: first.Position()}, Operands: []ast.Expr{first}}
	for {
		t := p.peek()
		op := ""
		switch {
		case t.Kind == token.Op && (t.Name == "<" || t.Name == ">" || t.Name == "<=" || t.Name == ">="):
			op = t.Name
		case t.IsKeyword("senaji"), t.IsKeyword("exi"):
			op = t.Name
		}
		if op == "" || p.peekAt(1).IsKeyword("izu") {
			break
		}
		if op == "exi" && !canStartExpr(p.peekAt(1)) {
			// 右に式がない x exi は「値があるか」
			if len(cmp.Ops) > 0 {
				return nil, p.unexpected(t)
			}
			p.next()
			return &ast.Exists{At: ast.At{Pos: t.Pos}, X: first}, nil
		}
		p.next()
		y, err := p.parseRange()
		if err != nil {
			return nil, err
		}
		cmp.Ops = append(cmp.Ops, op)
		cmp.Operands = append(cmp.Operands, y)
	}
	if len(cmp.Ops) == 0 {
		return first, nil
	}
	return cmp, nil
}

func (p *parser) parseRange() (ast.Expr, error) {
	lo, err := p.parseBinary(0)
	if err != nil {
		return nil, err
	}
	t := p.peek()
	if !t.IsOp("..") && !t.IsOp("...") {
		return lo, nil
	}
	p.next()
	r := &ast.Range{At: ast.At{Pos: t.Pos}, Lo: lo, Inclusive: t.Name == ".."}
	if r.Hi, err = p.parseBinary(0); err != nil {
		return nil, err
	}
	if p.acceptKeyword("batte") {
		if r.Step, err = p.parseBinary(0); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// binaryLevels は二項演算子を優先順位の低い順に並べたもの (Python と同じ)。
var binaryLevels = [][]string{
	{"|"},
	{"^"},
	{"&"},
	{"<<", ">>"},
	{"+", "-"},
	{"*", "/", "//", "%"},
}

func (p *parser) parseBinary(level int) (ast.Expr, error) {
	if level == len(binaryLevels) {
		return p.parseUnary()
	}
	x, err := p.parseBinary(level + 1)
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		if !p.isBinaryOp(t, binaryLevels[level]) {
			return x, nil
		}
		p.next()
		y, err := p.parseBinary(level + 1)
		if err != nil {
			return nil, err
		}
		x = &ast.Binary{At: ast.At{Pos: t.Pos}, Op: t.Name, X: x, Y: y}
	}
}

// isBinaryOp は t が ops のどれかかを調べる。x +izu 1 の + は演算子ではなく複合代入なので外す。
func (p *parser) isBinaryOp(t token.Token, ops []string) bool {
	if t.Kind != token.Op || p.peekAt(1).IsKeyword("izu") {
		return false
	}
	for _, op := range ops {
		if t.Name == op {
			return true
		}
	}
	return false
}

func (p *parser) parseUnary() (ast.Expr, error) {
	t := p.peek()
	if t.IsOp("-") || t.IsOp("+") || t.IsOp("~") {
		p.next()
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{At: ast.At{Pos: t.Pos}, Op: t.Name, X: x}, nil
	}
	return p.parsePower()
}

func (p *parser) parsePower() (ast.Expr, error) {
	x, err := p.parsePostfix()
	if err != nil {
		return nil, err
	}
	if t := p.peek(); p.isBinaryOp(t, []string{"**"}) {
		p.next()
		// 右結合で、右には符号を書ける (Python の 2 ** -1)
		y, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.Binary{At: ast.At{Pos: t.Pos}, Op: "**", X: x, Y: y}, nil
	}
	return x, nil
}

func (p *parser) parsePostfix() (ast.Expr, error) {
	x, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		switch {
		case t.IsOp("("):
			p.next()
			args, err := p.parseList(")")
			if err != nil {
				return nil, err
			}
			x = &ast.Call{At: ast.At{Pos: t.Pos}, Fn: x, Args: args}
		case t.IsKeyword("no"):
			p.next()
			if x, err = p.parseNo(x, ast.At{Pos: t.Pos}); err != nil {
				return nil, err
			}
		default:
			return x, nil
		}
	}
}

// parseNo は no の右側を読む。右に書けるのは数・識別子・文字列・括弧の式と、それを .. か ... でつないだ範囲。
func (p *parser) parseNo(x ast.Expr, at ast.At) (ast.Expr, error) {
	var lo ast.Expr
	if t := p.peek(); !t.IsOp("..") && !t.IsOp("...") {
		var err error
		if lo, err = p.parseNoOperand(); err != nil {
			return nil, err
		}
	}
	t := p.peek()
	if !t.IsOp("..") && !t.IsOp("...") {
		if lo == nil {
			return nil, p.unexpected(t)
		}
		return &ast.Index{At: at, X: x, Key: lo}, nil
	}
	p.next()
	s := &ast.Slice{At: at, X: x, Lo: lo, Inclusive: t.Name == ".."}
	if canStartNoOperand(p.peek(), p.peekAt(1)) {
		var err error
		if s.Hi, err = p.parseNoOperand(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func canStartNoOperand(t, next token.Token) bool {
	switch t.Kind {
	case token.Int, token.Float, token.Ident, token.String:
		return true
	case token.Op:
		return t.Name == "(" || (t.Name == "-" && (next.Kind == token.Int || next.Kind == token.Float))
	}
	return false
}

func (p *parser) parseNoOperand() (ast.Expr, error) {
	t := p.peek()
	if !canStartNoOperand(t, p.peekAt(1)) {
		return nil, p.unexpected(t)
	}
	if t.IsOp("-") {
		p.next()
		x, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{At: ast.At{Pos: t.Pos}, Op: "-", X: x}, nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (ast.Expr, error) {
	t := p.peek()
	at := ast.At{Pos: t.Pos}
	switch t.Kind {
	case token.Int:
		p.next()
		n, err := strconv.ParseInt(t.Name, 10, 64)
		if err != nil {
			return nil, diag.New(diag.Bavellabion, p.file, t.Pos, t.Name)
		}
		return &ast.IntLit{At: at, Value: n}, nil
	case token.Float:
		p.next()
		f, err := strconv.ParseFloat(t.Name, 64)
		if err != nil {
			return nil, diag.New(diag.Bavellabion, p.file, t.Pos, t.Name)
		}
		return &ast.FloatLit{At: at, Value: f}, nil
	case token.String:
		p.next()
		return p.parseString(t)
	case token.Ident:
		p.next()
		return &ast.Ident{At: at, Name: t.Name}, nil
	case token.Keyword:
		switch t.Name {
		case "tadact", "ulieso":
			p.next()
			return &ast.BoolLit{At: at, Value: t.Name == "tadact"}, nil
		case "vanita":
			p.next()
			return &ast.NullLit{At: at}, nil
		case "weho", "abiryoku", "itsver":
			return p.parseFuncLit()
		}
	case token.Op:
		switch t.Name {
		case "(":
			p.next()
			x, err := p.parseExprWithEnd()
			if err != nil {
				return nil, err
			}
			if err := p.expectOp(")"); err != nil {
				return nil, err
			}
			return &ast.Paren{At: at, X: x}, nil
		case "[":
			p.next()
			elems, err := p.parseList("]")
			if err != nil {
				return nil, err
			}
			return &ast.ArrayLit{At: at, Elems: elems}, nil
		case "{":
			p.next()
			return p.parseDict(at)
		}
	}
	return nil, p.unexpected(t)
}

func (p *parser) parseFuncLit() (ast.Expr, error) {
	at := ast.At{Pos: p.peek().Pos}
	var abiryoku, itsver bool
	for {
		switch t := p.peek(); {
		case t.IsKeyword("abiryoku") && !abiryoku:
			abiryoku = true
		case t.IsKeyword("itsver") && !itsver:
			itsver = true
		case t.IsKeyword("weho"):
			if p.peekAt(1).Kind == token.Ident {
				// 式の中の関数は名前を書かない
				return nil, p.unexpected(p.peekAt(1))
			}
			return p.parseFunc(at, abiryoku, itsver)
		default:
			return nil, p.unexpected(t)
		}
		p.next()
	}
}

// parseList は close までの、カンマで区切った式の並びを読む。最後のカンマは書いてもよい。
func (p *parser) parseList(close string) ([]ast.Expr, error) {
	var list []ast.Expr
	for !p.acceptOp(close) {
		x, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, x)
		if !p.acceptOp(",") {
			if err := p.expectOp(close); err != nil {
				return nil, err
			}
			break
		}
	}
	return list, nil
}

func (p *parser) parseDict(at ast.At) (ast.Expr, error) {
	d := &ast.DictLit{At: at}
	for !p.acceptOp("}") {
		k, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expectOp(":"); err != nil {
			return nil, err
		}
		v, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		d.Keys = append(d.Keys, k)
		d.Values = append(d.Values, v)
		if !p.acceptOp(",") {
			if err := p.expectOp("}"); err != nil {
				return nil, err
			}
			break
		}
	}
	return d, nil
}

// parseString は文字列を読み、#{} の中を式として読む。
func (p *parser) parseString(t token.Token) (ast.Expr, error) {
	s := &ast.StringLit{At: ast.At{Pos: t.Pos}}
	for _, part := range t.Parts {
		if !part.IsCode {
			s.Parts = append(s.Parts, ast.StringPiece{Text: part.Text})
			continue
		}
		tokens, _, err := lexer.LexRange(p.src, part.Start, part.Start+len(part.Text), part.Pos, p.file, p.mode)
		if err != nil {
			return nil, err
		}
		if p.macros != nil {
			if tokens, err = p.macros.expand(tokens, 0); err != nil {
				return nil, err
			}
		}
		sub := &parser{src: p.src, file: p.file, mode: p.mode, tokens: tokens, macros: p.macros}
		sub.skipNewlines()
		x, err := sub.parseExprWithEnd()
		if err != nil {
			return nil, err
		}
		sub.skipNewlines()
		if end := sub.peek(); end.Kind != token.EOF {
			return nil, sub.unexpected(end)
		}
		s.Parts = append(s.Parts, ast.StringPiece{Expr: x})
	}
	return s, nil
}

func canStartExpr(t token.Token) bool {
	switch t.Kind {
	case token.Int, token.Float, token.String, token.Ident:
		return true
	case token.Op:
		switch t.Name {
		case "(", "[", "{", "-", "+", "~":
			return true
		}
	case token.Keyword:
		switch t.Name {
		case "tadact", "ulieso", "vanita", "nii", "weho", "abiryoku", "itsver", "zanfession":
			return true
		}
	}
	return false
}
