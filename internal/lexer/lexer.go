// Package lexer はソースを字句に分ける。下書き (.fxma) と清書 (.fxmac) で語の読み方が違う。
package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Mode は語の読み方。ファイルの拡張子で決まる (ModeOf)。
type Mode int

const (
	Draft Mode = iota // 下書き。語は骨格で書き、方言も受け付ける
	Fair              // 清書。語は正式な綴りだけを受け付ける
)

// ModeOf はファイル名の拡張子からモードを決める。
func ModeOf(file string) Mode {
	if strings.HasSuffix(file, ".fxmac") {
		return Fair
	}
	return Draft
}

// Lex はソース全体を字句に分ける。
func Lex(src, file string, mode Mode) ([]token.Token, []diag.Warning, error) {
	return LexRange(src, 0, len(src), token.Pos{Line: 1, Col: 1}, file, mode)
}

// LexRange はソースの一部 [start, end) を字句に分ける。文字列の #{} の中を読むのに使う。
func LexRange(src string, start, end int, pos token.Pos, file string, mode Mode) ([]token.Token, []diag.Warning, error) {
	l := &lexer{src: src[:end], i: start, line: pos.Line, col: pos.Col, file: file, mode: mode}
	if err := l.run(); err != nil {
		return nil, nil, err
	}
	return l.tokens, l.warnings, nil
}

type lexer struct {
	src       string
	i         int
	line, col int
	file      string
	mode      Mode
	depth     int // 括弧の深さ。括弧の中の改行は文の区切りにしない
	tokens    []token.Token
	warnings  []diag.Warning
}

// 長いものから並べる。先に一致したものを使う。
var operators = []string{
	"...", "**", "//", "<<", ">>", "<=", ">=", "..", "==", "!=",
	"+", "-", "*", "/", "%", "&", "|", "^", "~", "<", ">",
	"(", ")", "[", "]", "{", "}", ",", ":",
}

func (l *lexer) run() error {
	for l.i < len(l.src) {
		c := l.src[l.i]
		switch {
		case c == '\n':
			if l.depth == 0 {
				l.emitNewline()
			}
			l.advance()
		case c == ' ' || c == '\t' || c == '\r':
			l.advance()
		case c == '#':
			for l.i < len(l.src) && l.src[l.i] != '\n' {
				l.advance()
			}
		case c == '"':
			if err := l.lexString(); err != nil {
				return err
			}
		case isDigit(c):
			if err := l.lexNumber(); err != nil {
				return err
			}
		case isWordStart(c):
			if err := l.lexWord(); err != nil {
				return err
			}
		default:
			if err := l.lexOperator(); err != nil {
				return err
			}
		}
	}
	l.tokens = append(l.tokens, token.Token{Kind: token.EOF, Pos: l.pos(), Start: l.i, End: l.i})
	return nil
}

func (l *lexer) pos() token.Pos { return token.Pos{Line: l.line, Col: l.col} }

func (l *lexer) advance() {
	r, size := utf8.DecodeRuneInString(l.src[l.i:])
	l.i += size
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
}

func (l *lexer) advanceTo(end int) {
	for l.i < end {
		l.advance()
	}
}

func (l *lexer) errorf(kind diag.Kind, pos token.Pos, target string) error {
	return diag.New(kind, l.file, pos, target)
}

func (l *lexer) emitNewline() {
	n := len(l.tokens)
	if n == 0 || l.tokens[n-1].Kind == token.Newline {
		return
	}
	l.tokens = append(l.tokens, token.Token{Kind: token.Newline, Pos: l.pos(), Start: l.i, End: l.i + 1})
}

func (l *lexer) lexOperator() error {
	pos, start := l.pos(), l.i
	for _, op := range operators {
		if !strings.HasPrefix(l.src[l.i:], op) {
			continue
		}
		if op == "==" || op == "!=" {
			// 等しさは senaji と nii で書く。記号を受け付けると書き方が2つになる
			return l.errorf(diag.AyAmis, pos, op)
		}
		l.advanceTo(l.i + len(op))
		switch op {
		case "(", "[", "{":
			l.depth++
		case ")", "]", "}":
			if l.depth > 0 {
				l.depth--
			}
		}
		l.tokens = append(l.tokens, token.Token{Kind: token.Op, Name: op, Pos: pos, Start: start, End: l.i})
		return nil
	}
	r, _ := utf8.DecodeRuneInString(l.src[l.i:])
	return l.errorf(diag.AyAmis, pos, string(r))
}

func (l *lexer) lexNumber() error {
	pos, start := l.pos(), l.i
	kind := token.Int
	l.skipDigits()
	if l.peek(0) == '.' && isDigit(l.peek(1)) {
		kind = token.Float
		l.advance()
		l.skipDigits()
	}
	if c := l.peek(0); c == 'e' || c == 'E' {
		kind = token.Float
		l.advance()
		if c := l.peek(0); c == '+' || c == '-' {
			l.advance()
		}
		if !isDigit(l.peek(0)) {
			return l.errorf(diag.AyAmis, pos, l.src[start:l.i])
		}
		l.skipDigits()
	}
	if isWordChar(l.peek(0)) {
		return l.errorf(diag.AyAmis, pos, l.src[start:l.i+1])
	}
	text := strings.ReplaceAll(l.src[start:l.i], "_", "")
	l.tokens = append(l.tokens, token.Token{Kind: kind, Name: text, Pos: pos, Start: start, End: l.i})
	return nil
}

func (l *lexer) skipDigits() {
	for isDigit(l.peek(0)) || (l.peek(0) == '_' && isDigit(l.peek(1))) {
		l.advance()
	}
}

func (l *lexer) peek(ahead int) byte {
	if l.i+ahead < len(l.src) {
		return l.src[l.i+ahead]
	}
	return 0
}

func (l *lexer) lexWord() error {
	if l.mode == Fair {
		if e, ok := lexicon.MatchSpelling(l.src[l.i:]); ok {
			end := l.i + len(e.Spelling)
			if end >= len(l.src) || !isWordChar(l.src[end]) {
				pos, start := l.pos(), l.i
				l.advanceTo(end)
				l.emitWord(e, pos, start)
				return nil
			}
		}
	}

	pos, start := l.pos(), l.i
	for l.i < len(l.src) && isWordChar(l.src[l.i]) {
		l.advance()
	}
	text := l.src[start:l.i]
	if strings.Contains(text, "_") {
		// 辞書の語に _ は入らないので、表と比べるまでもなく普通の識別子
		l.tokens = append(l.tokens, token.Token{Kind: token.Ident, Name: text, Pos: pos, Start: start, End: l.i})
		return nil
	}

	e, dialect, ok := lexicon.Lookup(strings.ToLower(text))
	switch {
	case !ok:
		l.tokens = append(l.tokens, token.Token{Kind: token.Ident, Name: text, Pos: pos, Start: start, End: l.i})
	case l.mode == Fair:
		// 清書に骨格や方言が残っているのは、手で書き換えたか壊れたかのどちらか
		return l.errorf(diag.AyAmis, pos, text)
	case dialect:
		l.warnings = append(l.warnings, diag.Warning{Kind: diag.DowGuil, From: text, To: e.Spelling, File: l.file, Pos: pos})
		l.emitWord(e, pos, start)
	case text != e.Name:
		// 大文字を混ぜた IZU や Izu は、キーワードと紛らわしいので受け付けない
		return l.errorf(diag.AyAmis, pos, text)
	default:
		l.emitWord(e, pos, start)
	}
	return nil
}

func (l *lexer) emitWord(e lexicon.Entry, pos token.Pos, start int) {
	kind := token.Ident
	if e.IsKeyword() {
		kind = token.Keyword
	}
	l.tokens = append(l.tokens, token.Token{Kind: kind, Name: e.Name, Pos: pos, Start: start, End: l.i, Lexicon: true})
}

func isDigit(c byte) bool     { return '0' <= c && c <= '9' }
func isWordStart(c byte) bool { return c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') }
func isWordChar(c byte) bool  { return isWordStart(c) || isDigit(c) }
