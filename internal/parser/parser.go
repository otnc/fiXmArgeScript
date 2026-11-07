// Package parser は字句の並びから構文木を作る (再帰下降法)。
package parser

import (
	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// File は読んだファイル1つ。
type File struct {
	Name     string
	Stmts    []ast.Stmt
	Warnings []diag.Warning
}

// Parse はソースを読んで構文木にする。下書きか清書かはファイル名の拡張子で決める。
func Parse(src, file string) (*File, error) {
	mode := lexer.ModeOf(file)
	tokens, warnings, err := lexer.Lex(src, file, mode)
	if err != nil {
		return nil, err
	}
	p := &parser{src: src, file: file, mode: mode, tokens: tokens}
	stmts, err := p.parseStmts()
	if err != nil {
		return nil, err
	}
	if t := p.peek(); t.Kind != token.EOF {
		return nil, p.unexpected(t)
	}
	return &File{Name: file, Stmts: stmts, Warnings: warnings}, nil
}

type parser struct {
	src    string
	file   string
	mode   lexer.Mode
	tokens []token.Token
	i      int
}

func (p *parser) peek() token.Token { return p.tokens[p.i] }

func (p *parser) peekAt(ahead int) token.Token {
	if p.i+ahead < len(p.tokens) {
		return p.tokens[p.i+ahead]
	}
	return p.tokens[len(p.tokens)-1]
}

func (p *parser) next() token.Token {
	t := p.tokens[p.i]
	if t.Kind != token.EOF {
		p.i++
	}
	return t
}

func (p *parser) acceptKeyword(name string) bool {
	if p.peek().IsKeyword(name) {
		p.next()
		return true
	}
	return false
}

func (p *parser) acceptOp(name string) bool {
	if p.peek().IsOp(name) {
		p.next()
		return true
	}
	return false
}

func (p *parser) expectKeyword(name string) error {
	if !p.acceptKeyword(name) {
		return p.unexpected(p.peek())
	}
	return nil
}

func (p *parser) expectOp(name string) error {
	if !p.acceptOp(name) {
		return p.unexpected(p.peek())
	}
	return nil
}

func (p *parser) expectIdent() (token.Token, error) {
	t := p.peek()
	if t.Kind != token.Ident {
		return t, p.unexpected(t)
	}
	return p.next(), nil
}

func (p *parser) skipNewlines() {
	for p.peek().Kind == token.Newline {
		p.next()
	}
}

// unexpected は、その字句がそこに書けないことを AyAmis で知らせる。
func (p *parser) unexpected(t token.Token) error {
	return diag.New(diag.AyAmis, p.file, t.Pos, p.text(t))
}

func (p *parser) text(t token.Token) string {
	switch t.Kind {
	case token.EOF:
		return "FineNd"
	case token.Newline:
		return `\n`
	}
	return p.src[t.Start:t.End]
}
